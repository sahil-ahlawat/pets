package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
)

type HealthService interface {
	ListVaccinationTypes(ctx context.Context, species string) ([]domain.VaccinationType, error)
	CreateVaccination(ctx context.Context, userID, petID uuid.UUID, v *domain.PetVaccination) error
	ListVaccinations(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetVaccination, error)
	CreateMedicalEvent(ctx context.Context, userID, petID uuid.UUID, m *domain.PetMedicalEvent) error
	ListMedicalEvents(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetMedicalEvent, error)
	CreateMedication(ctx context.Context, userID, petID uuid.UUID, m *domain.PetMedication) error
	ListMedications(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetMedication, error)
	GetTimeline(ctx context.Context, userID, petID uuid.UUID, limit int) ([]domain.PetTimelineEvent, error)
	CreateDocument(ctx context.Context, userID, petID uuid.UUID, doc *domain.PetDocument) error
	ListDocuments(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetDocument, error)
}

type healthService struct {
	healthRepo repository.HealthRepository
	petRepo    repository.PetRepository
}

func NewHealthService(healthRepo repository.HealthRepository, petRepo repository.PetRepository) HealthService {
	return &healthService{healthRepo: healthRepo, petRepo: petRepo}
}

func (s *healthService) ListVaccinationTypes(ctx context.Context, species string) ([]domain.VaccinationType, error) {
	return s.healthRepo.ListVaccinationTypes(ctx, species)
}

func (s *healthService) CreateVaccination(ctx context.Context, userID, petID uuid.UUID, v *domain.PetVaccination) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return errors.New("access denied to pet")
	}
	v.PetID = petID
	return s.healthRepo.CreateVaccination(ctx, v)
}

func (s *healthService) ListVaccinations(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetVaccination, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return nil, errors.New("access denied to pet")
	}
	return s.healthRepo.ListVaccinations(ctx, petID)
}

func (s *healthService) CreateMedicalEvent(ctx context.Context, userID, petID uuid.UUID, m *domain.PetMedicalEvent) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return errors.New("access denied to pet")
	}
	m.PetID = petID

	eventType := strings.ToUpper(strings.TrimSpace(m.EventType))
	switch eventType {
	case "CONSULTATION", "DIAGNOSIS", "TREATMENT", "MEDICATION", "SURGERY", "LAB_TEST", "DENTAL", "WEIGHT_CHECK", "OTHER":
		m.EventType = eventType
	default:
		m.EventType = "CONSULTATION"
	}

	return s.healthRepo.CreateMedicalEvent(ctx, m)
}

func (s *healthService) ListMedicalEvents(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetMedicalEvent, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return nil, errors.New("access denied to pet")
	}
	return s.healthRepo.ListMedicalEvents(ctx, petID)
}

func (s *healthService) CreateMedication(ctx context.Context, userID, petID uuid.UUID, m *domain.PetMedication) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return errors.New("access denied to pet")
	}
	m.PetID = petID
	return s.healthRepo.CreateMedication(ctx, m)
}

func (s *healthService) ListMedications(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetMedication, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return nil, errors.New("access denied to pet")
	}
	return s.healthRepo.ListMedications(ctx, petID)
}

func (s *healthService) GetTimeline(ctx context.Context, userID, petID uuid.UUID, limit int) ([]domain.PetTimelineEvent, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return nil, errors.New("access denied to pet")
	}
	return s.healthRepo.GetTimeline(ctx, petID, limit)
}

func (s *healthService) CreateDocument(ctx context.Context, userID, petID uuid.UUID, doc *domain.PetDocument) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return errors.New("access denied to pet")
	}
	doc.PetID = petID
	doc.UploadedBy = userID
	return s.healthRepo.CreateDocument(ctx, doc)
}

func (s *healthService) ListDocuments(ctx context.Context, userID, petID uuid.UUID) ([]domain.PetDocument, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil || !hasAccess {
		return nil, errors.New("access denied to pet")
	}
	return s.healthRepo.ListDocuments(ctx, petID)
}
