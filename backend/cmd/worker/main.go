package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"pet-platform/backend/internal/config"
	"pet-platform/backend/internal/repository"
	"pet-platform/backend/internal/worker"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewDB(cfg)
	if err != nil {
		fmt.Printf("Fatal database connection error: %v\n", err)
	}

	userRepo := repository.NewUserRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)

	w := worker.NewWorker(outboxRepo, userRepo, 4)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go w.Start(ctx)

	fmt.Println("🚀 Asynchronous Job Worker is online and waiting for jobs...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Terminating worker process gracefully...")
	w.Stop()
	fmt.Println("Worker terminated.")
}
