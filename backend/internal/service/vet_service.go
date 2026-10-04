package service

import (
	"context"

	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
)

type VetService interface {
	SearchVets(ctx context.Context, query, city string, emergencyOnly bool) ([]domain.Vet, error)
	FindNearbyVets(ctx context.Context, lat, lng, radiusKm float64) ([]domain.Vet, error)
}

type vetService struct {
	vetRepo repository.VetRepository
}

func NewVetService(vetRepo repository.VetRepository) VetService {
	return &vetService{vetRepo: vetRepo}
}

func (s *vetService) SearchVets(ctx context.Context, query, city string, emergencyOnly bool) ([]domain.Vet, error) {
	return s.vetRepo.SearchVets(ctx, query, city, emergencyOnly)
}

func (s *vetService) FindNearbyVets(ctx context.Context, lat, lng, radiusKm float64) ([]domain.Vet, error) {
	return s.vetRepo.FindNearbyVets(ctx, lat, lng, radiusKm)
}
