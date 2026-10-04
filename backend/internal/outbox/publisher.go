package outbox

import (
	"context"
	"fmt"
	"time"

	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
)

type EventPublisher struct {
	repo       repository.OutboxRepository
	pollPeriod time.Duration
	stopCh     chan struct{}
}

func NewEventPublisher(repo repository.OutboxRepository, pollPeriod time.Duration) *EventPublisher {
	if pollPeriod <= 0 {
		pollPeriod = 2 * time.Second
	}
	return &EventPublisher{
		repo:       repo,
		pollPeriod: pollPeriod,
		stopCh:     make(chan struct{}),
	}
}

func (p *EventPublisher) Start(ctx context.Context) {
	ticker := time.NewTicker(p.pollPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *EventPublisher) Stop() {
	close(p.stopCh)
}

func (p *EventPublisher) processBatch(ctx context.Context) {
	events, err := p.repo.FetchUnpublishedEvents(ctx, 25)
	if err != nil {
		fmt.Printf("Error fetching outbox events: %v\n", err)
		return
	}

	for _, ev := range events {
		p.publishEvent(ctx, &ev)
	}
}

func (p *EventPublisher) publishEvent(ctx context.Context, ev *domain.OutboxEvent) {
	// 1. Enqueue corresponding background jobs based on event type
	switch ev.EventType {
	case "QRScanned":
		_ = p.repo.EnqueueJob(ctx, "notifications", "SEND_SCAN_NOTIFICATION", ev.Payload, time.Now().UTC())
	case "PetMarkedLost":
		_ = p.repo.EnqueueJob(ctx, "broadcast", "BROADCAST_LOST_PET_ALERT", ev.Payload, time.Now().UTC())
	case "PetFound":
		_ = p.repo.EnqueueJob(ctx, "notifications", "SEND_FOUND_PET_NOTIFICATION", ev.Payload, time.Now().UTC())
	case "VaccinationAdded":
		_ = p.repo.EnqueueJob(ctx, "reminders", "SCHEDULE_VACCINE_BOOSTER_REMINDER", ev.Payload, time.Now().UTC())
	}

	// 2. Mark event as successfully published
	if err := p.repo.MarkEventPublished(ctx, ev.ID); err != nil {
		fmt.Printf("Error marking outbox event %s published: %v\n", ev.ID, err)
	} else {
		fmt.Printf("[Outbox] Published event %s (Type: %s, Aggregate: %s)\n", ev.ID, ev.EventType, ev.AggregateID)
	}
}
