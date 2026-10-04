package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
)

type ReminderRepository interface {
	CreateReminder(ctx context.Context, r *domain.PetReminder) error
	ListUserReminders(ctx context.Context, userID uuid.UUID) ([]domain.PetReminder, error)
	MarkReminderCompleted(ctx context.Context, reminderID, userID uuid.UUID) error
	GetPendingDueReminders(ctx context.Context, dueBefore time.Time) ([]domain.PetReminder, error)
	UpdateReminderNotified(ctx context.Context, reminderID uuid.UUID) error
}

type sqlReminderRepo struct {
	db *DB
}

func NewReminderRepository(db *DB) ReminderRepository {
	return &sqlReminderRepo{db: db}
}

func (r *sqlReminderRepo) CreateReminder(ctx context.Context, rem *domain.PetReminder) error {
	if rem.ID == uuid.Nil {
		rem.ID = uuid.New()
	}
	if rem.Recurrence == "" {
		rem.Recurrence = "ONCE"
	}
	now := time.Now().UTC()

	query := `
		INSERT INTO pet_reminders (id, pet_id, user_id, reminder_type, title, description, due_date, recurrence, is_completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, FALSE, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		rem.ID, rem.PetID, rem.UserID, rem.ReminderType, rem.Title,
		rem.Description, rem.DueDate, rem.Recurrence, now, now,
	)
	return err
}

func (r *sqlReminderRepo) ListUserReminders(ctx context.Context, userID uuid.UUID) ([]domain.PetReminder, error) {
	query := `
		SELECT pr.id, pr.pet_id, p.name, pr.user_id, pr.reminder_type, pr.title, COALESCE(pr.description, ''), pr.due_date, pr.recurrence, pr.is_completed, pr.completed_at
		FROM pet_reminders pr
		JOIN pets p ON p.id = pr.pet_id
		WHERE pr.user_id = $1
		ORDER BY pr.is_completed ASC, pr.due_date ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PetReminder
	for rows.Next() {
		var rem domain.PetReminder
		var compAt sql.NullTime
		if err := rows.Scan(
			&rem.ID, &rem.PetID, &rem.PetName, &rem.UserID, &rem.ReminderType,
			&rem.Title, &rem.Description, &rem.DueDate, &rem.Recurrence,
			&rem.IsCompleted, &compAt,
		); err != nil {
			return nil, err
		}
		if compAt.Valid {
			t := compAt.Time
			rem.CompletedAt = &t
		}
		list = append(list, rem)
	}
	return list, nil
}

func (r *sqlReminderRepo) MarkReminderCompleted(ctx context.Context, reminderID, userID uuid.UUID) error {
	query := `
		UPDATE pet_reminders
		SET is_completed = TRUE, completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, reminderID, userID)
	return err
}

func (r *sqlReminderRepo) GetPendingDueReminders(ctx context.Context, dueBefore time.Time) ([]domain.PetReminder, error) {
	query := `
		SELECT pr.id, pr.pet_id, p.name, pr.user_id, pr.reminder_type, pr.title, COALESCE(pr.description, ''), pr.due_date, pr.recurrence, pr.is_completed, pr.completed_at
		FROM pet_reminders pr
		JOIN pets p ON p.id = pr.pet_id
		WHERE pr.is_completed = FALSE 
		  AND pr.due_date <= $1 
		  AND (pr.last_notified_at IS NULL OR pr.last_notified_at < CURRENT_TIMESTAMP - INTERVAL '1 day')
		LIMIT 100
	`
	rows, err := r.db.QueryContext(ctx, query, dueBefore)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PetReminder
	for rows.Next() {
		var rem domain.PetReminder
		var compAt sql.NullTime
		if err := rows.Scan(
			&rem.ID, &rem.PetID, &rem.PetName, &rem.UserID, &rem.ReminderType,
			&rem.Title, &rem.Description, &rem.DueDate, &rem.Recurrence,
			&rem.IsCompleted, &compAt,
		); err != nil {
			return nil, err
		}
		list = append(list, rem)
	}
	return list, nil
}

func (r *sqlReminderRepo) UpdateReminderNotified(ctx context.Context, reminderID uuid.UUID) error {
	query := `UPDATE pet_reminders SET last_notified_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, reminderID)
	return err
}
