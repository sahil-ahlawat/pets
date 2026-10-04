package service

import (
	"context"
	"errors"
	"strings"

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

	remType := strings.ToUpper(strings.TrimSpace(r.ReminderType))
	switch remType {
	case "VACCINATION", "MEDICATION", "DEWORMING", "FLEA_TICK", "VET_APPOINTMENT", "GROOMING", "CUSTOM":
		r.ReminderType = remType
	default:
		r.ReminderType = "CUSTOM"
	}

	recurrence := strings.ToUpper(strings.TrimSpace(r.Recurrence))
	switch recurrence {
	case "ONCE", "DAILY", "WEEKLY", "MONTHLY", "QUARTERLY", "BIANNUAL", "YEARLY":
		r.Recurrence = recurrence
	default:
		r.Recurrence = "ONCE"
	}

	return s.reminderRepo.CreateReminder(ctx, r)
}

func (s *reminderService) ListUserReminders(ctx context.Context, userID uuid.UUID) ([]domain.PetReminder, error) {
	return s.reminderRepo.ListUserReminders(ctx, userID)
}

func (s *reminderService) MarkReminderCompleted(ctx context.Context, reminderID, userID uuid.UUID) error {
	return s.reminderRepo.MarkReminderCompleted(ctx, reminderID, userID)
}
