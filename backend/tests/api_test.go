package tests

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"pet-platform/backend/internal/middleware"
)

func TestPasswordHashing(t *testing.T) {
	password := "SecretPetPassword123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		t.Errorf("Password validation failed: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword")); err == nil {
		t.Errorf("Expected failure for wrong password")
	}
}

func TestJWTTokenLifecycle(t *testing.T) {
	secret := "test-secret-key-12345"
	userID := uuid.New().String()
	email := "petowner@example.com"
	role := "PET_OWNER"

	claims := &middleware.JWTClaims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pet-identity-platform",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Signing failed: %v", err)
	}

	// Parse & verify
	parsedClaims := &middleware.JWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenStr, parsedClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsedToken.Valid {
		t.Fatalf("Parsing failed: %v", err)
	}

	if parsedClaims.UserID != userID || parsedClaims.Email != email {
		t.Errorf("Claims mismatch: got %v, expected %v", parsedClaims.UserID, userID)
	}
}
