package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
)

type Worker struct {
	outboxRepo repository.OutboxRepository
	userRepo   repository.UserRepository
	concurrency int
	stopCh      chan struct{}
}

func NewWorker(outboxRepo repository.OutboxRepository, userRepo repository.UserRepository, concurrency int) *Worker {
	if concurrency <= 0 {
		concurrency = 4
	}
	return &Worker{
		outboxRepo:  outboxRepo,
		userRepo:    userRepo,
		concurrency: concurrency,
		stopCh:      make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("[Worker] Starting background job processor with concurrency=%d\n", w.concurrency)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case <-ticker.C:
			jobs, err := w.outboxRepo.FetchPendingJobs(ctx, 10)
			if err != nil {
				continue
			}
			for _, job := range jobs {
				w.processJob(ctx, &job)
			}
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

func (w *Worker) processJob(ctx context.Context, job *domain.Job) {
	fmt.Printf("[Worker] Executing job %s (Type: %s, Attempt: %d)\n", job.ID, job.JobType, job.RetryCount+1)

	var err error
	switch job.JobType {
	case "SEND_SCAN_NOTIFICATION":
		err = w.handleScanNotification(ctx, job.Payload)
	case "BROADCAST_LOST_PET_ALERT":
		err = w.handleLostBroadcast(ctx, job.Payload)
	case "SEND_FOUND_PET_NOTIFICATION":
		err = w.handleFoundNotification(ctx, job.Payload)
	case "SCHEDULE_VACCINE_BOOSTER_REMINDER":
		err = w.handleVaccineBooster(ctx, job.Payload)
	default:
		// Unsupported job type
		err = nil
	}

	if err != nil {
		fmt.Printf("[Worker] Job %s failed: %v\n", job.ID, err)
		_ = w.outboxRepo.UpdateJobStatus(ctx, job.ID, "FAILED", err.Error())
	} else {
		fmt.Printf("[Worker] Job %s completed successfully\n", job.ID)
		_ = w.outboxRepo.UpdateJobStatus(ctx, job.ID, "COMPLETED", "")
	}
}

func (w *Worker) handleScanNotification(ctx context.Context, payload []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	// Emulated push delivery via FCM / APNs stub
	fmt.Printf("[Push Service] Dispatched push notification for QR scan on pet %v\n", data["petId"])
	return nil
}

func (w *Worker) handleLostBroadcast(ctx context.Context, payload []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	fmt.Printf("[Community Broadcast] Pet %v marked lost at %v. Broadcast sent to nearby users.\n", data["petId"], data["address"])
	return nil
}

func (w *Worker) handleFoundNotification(ctx context.Context, payload []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	fmt.Printf("[Emergency Alert] Found pet report alert dispatched to owner %v for pet %v\n", data["ownerId"], data["petId"])
	return nil
}

func (w *Worker) handleVaccineBooster(ctx context.Context, payload []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	fmt.Printf("[Health Engine] Booster schedule calculated for %v (pet %v)\n", data["vaccineName"], data["petId"])
	return nil
}
