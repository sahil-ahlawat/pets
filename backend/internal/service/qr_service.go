package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"pet-platform/backend/internal/config"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
)

type QRService interface {
	GetQRCodes(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetQRCode, error)
	GenerateNewQRCode(ctx context.Context, userID, petID uuid.UUID, tagLabel string) (*domain.PetQRCode, error)
	RevokeQRCode(ctx context.Context, userID, qrID uuid.UUID, reason string) error
	GetPublicProfile(ctx context.Context, token string, clientIP, userAgent string, lat, lng *float64) (*domain.PublicPetProfile, error)
	GenerateQRCodePNG(token string, size int) ([]byte, error)
	ContactOwner(ctx context.Context, token string, req *domain.ContactOwnerRequest) error
	ReportFound(ctx context.Context, token string, req *domain.ReportFoundRequest) error
}

type qrService struct {
	cfg      *config.Config
	qrRepo   repository.QRRepository
	petRepo  repository.PetRepository
	userRepo repository.UserRepository
}

func NewQRService(cfg *config.Config, qrRepo repository.QRRepository, petRepo repository.PetRepository, userRepo repository.UserRepository) QRService {
	return &qrService{
		cfg:      cfg,
		qrRepo:   qrRepo,
		petRepo:  petRepo,
		userRepo: userRepo,
	}
}

func (s *qrService) GetQRCodes(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetQRCode, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return nil, errors.New("access denied")
	}
	return s.qrRepo.GetQRCodesByPetID(ctx, petID)
}

func (s *qrService) GenerateNewQRCode(ctx context.Context, userID, petID uuid.UUID, tagLabel string) (*domain.PetQRCode, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return nil, errors.New("access denied")
	}
	return s.qrRepo.GenerateNewQRCode(ctx, petID, tagLabel)
}

func (s *qrService) RevokeQRCode(ctx context.Context, userID, qrID uuid.UUID, reason string) error {
	// Revoke tag
	return s.qrRepo.RevokeQRCode(ctx, qrID, reason)
}

func (s *qrService) GetPublicProfile(ctx context.Context, token string, clientIP, userAgent string, lat, lng *float64) (*domain.PublicPetProfile, error) {
	pub, petID, ownerID, qrID, err := s.qrRepo.GetPublicPetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if pub == nil {
		return nil, errors.New("invalid or revoked QR token")
	}

	// Record scan in background/transaction asynchronously
	if petID != nil && qrID != nil {
		scan := &domain.QRScan{
			ID:              uuid.New(),
			QRCodeID:        *qrID,
			PetID:           *petID,
			ApproxLatitude:  lat,
			ApproxLongitude: lng,
			IPHash:          clientIP,
			UserAgent:       userAgent,
		}
		_ = s.qrRepo.RecordScan(ctx, scan)

		// Notify owner of scan
		if ownerID != nil {
			title := fmt.Sprintf("Tag Scanned: %s", pub.Name)
			body := "Someone just scanned your pet's digital QR collar tag!"
			_ = s.userRepo.CreateNotification(ctx, &domain.Notification{
				UserID:           *ownerID,
				PetID:            petID,
				NotificationType: "QR_SCANNED",
				Title:            title,
				Body:             body,
			})
		}
	}

	return pub, nil
}

func (s *qrService) GenerateQRCodePNG(token string, size int) ([]byte, error) {
	if size <= 0 || size > 2048 {
		size = 512
	}
	// Deep link URL that opens native app if installed, or web profile
	targetURL := fmt.Sprintf("%s/p/%s", s.cfg.PublicBaseURL, token)
	return qrcode.Encode(targetURL, qrcode.Medium, size)
}

func (s *qrService) ContactOwner(ctx context.Context, token string, req *domain.ContactOwnerRequest) error {
	pub, petID, ownerID, _, err := s.qrRepo.GetPublicPetByToken(ctx, token)
	if err != nil || pub == nil || ownerID == nil {
		return errors.New("unable to reach pet owner")
	}

	title := fmt.Sprintf("💬 Message Regarding %s", pub.Name)
	contactFrom := req.FinderName
	if contactFrom == "" {
		contactFrom = "A community member"
	}
	body := fmt.Sprintf("%s sent a message: \"%s\". Finder phone: %s, Email: %s",
		contactFrom, req.Message, req.FinderPhone, req.FinderEmail)

	return s.userRepo.CreateNotification(ctx, &domain.Notification{
		UserID:           *ownerID,
		PetID:            petID,
		NotificationType: "COMMUNITY_MESSAGE",
		Title:            title,
		Body:             body,
	})
}

func (s *qrService) ReportFound(ctx context.Context, token string, req *domain.ReportFoundRequest) error {
	if req.Message == "" {
		return errors.New("report message is required")
	}
	return s.qrRepo.RecordFoundReport(ctx, token, req)
}
