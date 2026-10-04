package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pet-platform/backend/internal/config"
	"pet-platform/backend/internal/repository"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewDB(cfg)
	if err != nil {
		fmt.Printf("Fatal database connection error: %v\n", err)
	}

	reminderRepo := repository.NewReminderRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("⏰ Cron Scheduler is active (monitoring due reminders and boosters)...")

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				checkDueReminders(ctx, reminderRepo, outboxRepo)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Scheduler shutting down...")
}

func checkDueReminders(ctx context.Context, remRepo repository.ReminderRepository, outboxRepo repository.OutboxRepository) {
	now := time.Now().UTC()
	reminders, err := remRepo.GetPendingDueReminders(ctx, now)
	if err != nil {
		return
	}

	for _, rem := range reminders {
		fmt.Printf("[Scheduler] Enqueuing due reminder notification for pet %s: %s\n", rem.PetName, rem.Title)
		_ = outboxRepo.EnqueueJob(ctx, "notifications", "SEND_REMINDER_NOTIFICATION", map[string]interface{}{
			"reminderId": rem.ID.String(),
			"userId":     rem.UserID.String(),
			"petName":    rem.PetName,
			"title":      rem.Title,
		}, now)

		_ = remRepo.UpdateReminderNotified(ctx, rem.ID)
	}
}
