package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
)

type PetService interface {
	CreatePet(ctx context.Context, userID uuid.UUID, req *domain.CreatePetRequest) (*domain.PetDetail, error)
	ListPets(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.Pet, string, error)
	GetPet(ctx context.Context, userID, petID uuid.UUID) (*domain.PetDetail, error)
	UpdatePet(ctx context.Context, userID, petID uuid.UUID, req *domain.UpdatePetRequest) (*domain.PetDetail, error)
	DeletePet(ctx context.Context, userID, petID uuid.UUID) error
	MarkPetLost(ctx context.Context, userID, petID uuid.UUID, req *domain.MarkLostRequest) error
	MarkPetReunited(ctx context.Context, userID, petID uuid.UUID) error
}

type petService struct {
	petRepo repository.PetRepository
}

func NewPetService(petRepo repository.PetRepository) PetService {
	return &petService{petRepo: petRepo}
}

func (s *petService) CreatePet(ctx context.Context, userID uuid.UUID, req *domain.CreatePetRequest) (*domain.PetDetail, error) {
	if req.Name == "" {
		return nil, errors.New("pet name is required")
	}
	if req.Species == "" {
		req.Species = domain.SpeciesDog
	}

	pet := &domain.Pet{
		ID:         uuid.New(),
		Name:       req.Name,
		Species:    req.Species,
		Breed:      req.Breed,
		Gender:     req.Gender,
		IsNeutered: req.IsNeutered,
		WeightKg:   req.WeightKg,
		AvatarURL:  req.AvatarURL,
	}

	return s.petRepo.CreatePet(ctx, userID, pet, req)
}

func (s *petService) ListPets(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.Pet, string, error) {
	return s.petRepo.ListPetsByUserID(ctx, userID, limit, cursor)
}

func (s *petService) GetPet(ctx context.Context, userID, petID uuid.UUID) (*domain.PetDetail, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("access denied: pet does not belong to user")
	}

	detail, err := s.petRepo.GetPetByID(ctx, petID)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errors.New("pet not found")
	}
	return detail, nil
}

func (s *petService) UpdatePet(ctx context.Context, userID, petID uuid.UUID, req *domain.UpdatePetRequest) (*domain.PetDetail, error) {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("access denied")
	}
	return s.petRepo.UpdatePet(ctx, petID, req)
}

func (s *petService) DeletePet(ctx context.Context, userID, petID uuid.UUID) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("access denied")
	}
	return s.petRepo.DeletePet(ctx, petID)
}

func (s *petService) MarkPetLost(ctx context.Context, userID, petID uuid.UUID, req *domain.MarkLostRequest) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("access denied")
	}
	if req.LastSeenAddress == "" {
		return errors.New("last seen location is required when reporting a lost pet")
	}
	return s.petRepo.MarkPetLost(ctx, petID, userID, req)
}

func (s *petService) MarkPetReunited(ctx context.Context, userID, petID uuid.UUID) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, petID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("access denied")
	}
	return s.petRepo.MarkPetReunited(ctx, petID, userID)
}
