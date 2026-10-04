package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
)

type OutboxRepository interface {
	FetchUnpublishedEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkEventPublished(ctx context.Context, id uuid.UUID) error
	RecordConsumerProcessed(ctx context.Context, eventID uuid.UUID, consumerName string) error
	IsConsumerProcessed(ctx context.Context, eventID uuid.UUID, consumerName string) (bool, error)
	EnqueueJob(ctx context.Context, queue, jobType string, payload interface{}, runAt time.Time) error
	FetchPendingJobs(ctx context.Context, limit int) ([]domain.Job, error)
	UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status, lastError string) error
}

type sqlOutboxRepo struct {
	db *DB
}

func NewOutboxRepository(db *DB) OutboxRepository {
	return &sqlOutboxRepo{db: db}
}

func (r *sqlOutboxRepo) FetchUnpublishedEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, event_type, aggregate_type, aggregate_id, payload, created_at, is_published
		FROM outbox_events
		WHERE is_published = FALSE
		ORDER BY created_at ASC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.OutboxEvent
	for rows.Next() {
		var ev domain.OutboxEvent
		var payloadBytes []byte
		if err := rows.Scan(&ev.ID, &ev.EventType, &ev.AggregateType, &ev.AggregateID, &payloadBytes, &ev.CreatedAt, &ev.IsPublished); err != nil {
			return nil, err
		}
		ev.Payload = payloadBytes
		events = append(events, ev)
	}
	return events, nil
}

func (r *sqlOutboxRepo) MarkEventPublished(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE outbox_events SET is_published = TRUE, published_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *sqlOutboxRepo) RecordConsumerProcessed(ctx context.Context, eventID uuid.UUID, consumerName string) error {
	query := `
		INSERT INTO outbox_consumer_tracking (id, event_id, consumer_name, processed_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (event_id, consumer_name) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, uuid.New(), eventID, consumerName)
	return err
}

func (r *sqlOutboxRepo) IsConsumerProcessed(ctx context.Context, eventID uuid.UUID, consumerName string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM outbox_consumer_tracking WHERE event_id = $1 AND consumer_name = $2`
	err := r.db.QueryRowContext(ctx, query, eventID, consumerName).Scan(&count)
	return count > 0, err
}

func (r *sqlOutboxRepo) EnqueueJob(ctx context.Context, queue, jobType string, payload interface{}, runAt time.Time) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if runAt.IsZero() {
		runAt = time.Now().UTC()
	}

	query := `
		INSERT INTO jobs (id, queue, job_type, payload, status, run_at, max_retries, retry_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'QUEUED', $5, 5, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err = r.db.ExecContext(ctx, query, uuid.New(), queue, jobType, payloadBytes, runAt)
	return err
}

func (r *sqlOutboxRepo) FetchPendingJobs(ctx context.Context, limit int) ([]domain.Job, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		UPDATE jobs
		SET status = 'PROCESSING', started_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id IN (
			SELECT id FROM jobs
			WHERE status = 'QUEUED' AND run_at <= CURRENT_TIMESTAMP
			ORDER BY run_at ASC
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, queue, job_type, payload, status, run_at, retry_count
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []domain.Job
	for rows.Next() {
		var j domain.Job
		var payloadBytes []byte
		if err := rows.Scan(&j.ID, &j.Queue, &j.JobType, &payloadBytes, &j.Status, &j.RunAt, &j.RetryCount); err != nil {
			return nil, err
		}
		j.Payload = payloadBytes
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *sqlOutboxRepo) UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status, lastError string) error {
	var query string
	if status == "COMPLETED" {
		query = `UPDATE jobs SET status = 'COMPLETED', finished_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
		_, err := r.db.ExecContext(ctx, query, jobID)
		return err
	} else if status == "FAILED" {
		query = `
			UPDATE jobs SET
				retry_count = retry_count + 1,
				last_error = $2,
				status = CASE WHEN retry_count + 1 >= max_retries THEN 'DEAD_LETTER'::job_status ELSE 'QUEUED'::job_status END,
				run_at = CURRENT_TIMESTAMP + (INTERVAL '1 second' * POWER(2, retry_count + 1)),
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
		_, err := r.db.ExecContext(ctx, query, jobID, lastError)
		return err
	}
	return nil
}
