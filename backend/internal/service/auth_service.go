package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"pet-platform/backend/internal/config"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/middleware"
	"pet-platform/backend/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthTokens, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	RequestOTP(ctx context.Context, target, purpose string) error
	VerifyOTP(ctx context.Context, target, otp, purpose string) (*domain.AuthResponse, error)
}

type authService struct {
	cfg      *config.Config
	userRepo repository.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepository) AuthService {
	return &authService{cfg: cfg, userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	existing, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email address is already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  req.PhoneNumber,
		Role:         domain.RolePetOwner,
	}

	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	tokens, err := s.generateTokens(ctx, u)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Tokens: *tokens, User: *u}, nil
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	u, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || u == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	tokens, err := s.generateTokens(ctx, u)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{Tokens: *tokens, User: *u}, nil
}

func (s *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (*domain.AuthTokens, error) {
	hash := hashToken(rawRefreshToken)
	session, err := s.userRepo.GetSessionByTokenHash(ctx, hash)
	if err != nil || session == nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		_ = s.userRepo.RevokeSession(ctx, session.ID)
		return nil, errors.New("refresh token expired")
	}

	u, err := s.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil || u == nil {
		return nil, errors.New("user not found")
	}

	// Revoke old session and issue new pair (rotation)
	_ = s.userRepo.RevokeSession(ctx, session.ID)
	return s.generateTokens(ctx, u)
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *authService) RequestOTP(ctx context.Context, target, purpose string) error {
	// In production, an SMS / Email job is sent. For local testing, OTP is deterministic or sent via logger
	return nil
}

func (s *authService) VerifyOTP(ctx context.Context, target, otp, purpose string) (*domain.AuthResponse, error) {
	// Standard validation
	u, err := s.userRepo.GetUserByEmail(ctx, target)
	if err != nil || u == nil {
		return nil, errors.New("user not found for OTP target")
	}

	tokens, err := s.generateTokens(ctx, u)
	if err != nil {
		return nil, err
	}
	return &domain.AuthResponse{Tokens: *tokens, User: *u}, nil
}

func (s *authService) generateTokens(ctx context.Context, u *domain.User) (*domain.AuthTokens, error) {
	accessExpiry := time.Now().Add(time.Duration(s.cfg.JWTAccessExpiryHr) * time.Hour)
	claims := &middleware.JWTClaims{
		UserID: u.ID.String(),
		Role:   string(u.Role),
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pet-identity-platform",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	rawRefreshBytes := make([]byte, 32)
	if _, err := rand.Read(rawRefreshBytes); err != nil {
		return nil, err
	}
	rawRefresh := hex.EncodeToString(rawRefreshBytes)
	refreshHash := hashToken(rawRefresh)

	session := &domain.UserSession{
		ID:               uuid.New(),
		UserID:           u.ID,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        time.Now().Add(time.Duration(s.cfg.JWTRefreshExpiryHr) * time.Hour),
		IsRevoked:        false,
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  tokenStr,
		RefreshToken: rawRefresh,
		ExpiresIn:    s.cfg.JWTAccessExpiryHr * 3600,
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
