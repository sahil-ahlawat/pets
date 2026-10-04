package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, userID uuid.UUID, r *domain.PetReminder) error
	ListUserReminders(ctx context.Context, userID uuid.UUID) ([]domain.PetReminder, error)
	MarkReminderCompleted(ctx context.Context, reminderID, userID uuid.UUID) error
}

type reminderService struct {
	reminderRepo repository.ReminderRepository
	petRepo      repository.PetRepository
}

func NewReminderService(reminderRepo repository.ReminderRepository, petRepo repository.PetRepository) ReminderService {
	return &reminderService{reminderRepo: reminderRepo, petRepo: petRepo}
}

func (s *reminderService) CreateReminder(ctx context.Context, userID uuid.UUID, r *domain.PetReminder) error {
	hasAccess, err := s.petRepo.UserHasPetAccess(ctx, userID, r.PetID)
	if err != nil || !hasAccess {
		return errors.New("access denied to pet")
	}
	r.UserID = userID
	return s.reminderRepo.CreateReminder(ctx, r)
}

func (s *reminderService) ListUserReminders(ctx context.Context, userID uuid.UUID) ([]domain.PetReminder, error) {
	return s.reminderRepo.ListUserReminders(ctx, userID)
}

func (s *reminderService) MarkReminderCompleted(ctx context.Context, reminderID, userID uuid.UUID) error {
	return s.reminderRepo.MarkReminderCompleted(ctx, reminderID, userID)
}
