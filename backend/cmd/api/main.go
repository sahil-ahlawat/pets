package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"pet-platform/backend/internal/config"
	"pet-platform/backend/internal/handler"
	"pet-platform/backend/internal/middleware"
	"pet-platform/backend/internal/outbox"
	"pet-platform/backend/internal/repository"
	"pet-platform/backend/internal/service"
	"pet-platform/backend/internal/storage"
)

func main() {
	cfg := config.Load()

	// 1. Database
	db, err := repository.NewDB(cfg)
	if err != nil {
		fmt.Printf("Fatal database init error: %v\n", err)
	}

	// 2. Repositories
	userRepo := repository.NewUserRepository(db)
	petRepo := repository.NewPetRepository(db)
	healthRepo := repository.NewHealthRepository(db)
	qrRepo := repository.NewQRRepository(db)
	vetRepo := repository.NewVetRepository(db)
	reminderRepo := repository.NewReminderRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)

	// 3. Storage
	storageSvc := storage.NewStorageService(cfg)

	// 4. Services
	authSvc := service.NewAuthService(cfg, userRepo)
	petSvc := service.NewPetService(petRepo)
	healthSvc := service.NewHealthService(healthRepo, petRepo)
	qrSvc := service.NewQRService(cfg, qrRepo, petRepo, userRepo)
	vetSvc := service.NewVetService(vetRepo)
	reminderSvc := service.NewReminderService(reminderRepo, petRepo)

	// 5. Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	petHandler := handler.NewPetHandler(petSvc)
	healthHandler := handler.NewHealthHandler(healthSvc)
	qrHandler := handler.NewQRHandler(qrSvc)
	publicHandler := handler.NewPublicHandler(qrSvc)
	vetHandler := handler.NewVetHandler(vetSvc)
	reminderHandler := handler.NewReminderHandler(reminderSvc)
	storageHandler := handler.NewStorageHandler(storageSvc)
	healthzHandler := handler.NewHealthzHandler(db)

	// 6. Router Setup
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Total-Count", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health & Ops endpoints
	r.Get("/healthz", healthzHandler.Healthz)
	r.Get("/readyz", healthzHandler.Readyz)
	r.Get("/metrics", healthzHandler.Metrics)

	// Deep-link redirect: /p/{token} -> frontend public pet page
	r.Get("/p/{token}", func(w http.ResponseWriter, req *http.Request) {
		token := chi.URLParam(req, "token")
		http.Redirect(w, req, fmt.Sprintf("/public-pet.html?token=%s", token), http.StatusFound)
	})

	// API v1
	r.Route("/api/v1", func(v1 chi.Router) {
		// Public Auth
		v1.Post("/auth/register", authHandler.Register)
		v1.Post("/auth/login", authHandler.Login)
		v1.Post("/auth/refresh", authHandler.Refresh)
		v1.Post("/auth/otp/request", authHandler.RequestOTP)
		v1.Post("/auth/otp/verify", authHandler.VerifyOTP)

		// Public QR Scan Experience
		v1.Get("/public/pets/{token}", publicHandler.GetPublicPet)
		v1.Post("/public/pets/{token}/contact", publicHandler.ContactOwner)
		v1.Post("/public/pets/{token}/found", publicHandler.ReportFound)
		v1.Get("/public/qr/{token}.png", qrHandler.RenderQRCodePNG)

		// Public Catalog
		v1.Get("/vaccination-types", healthHandler.ListVaccinationTypes)
		v1.Get("/vets", vetHandler.SearchVets)
		v1.Get("/vets/nearby", vetHandler.FindNearbyVets)

		// Protected Routes
		v1.Group(func(auth chi.Router) {
			auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			// User
			auth.Get("/auth/me", authHandler.Me)

			// Pets
			auth.Get("/pets", petHandler.ListPets)
			auth.Post("/pets", petHandler.CreatePet)
			auth.Get("/pets/{id}", petHandler.GetPet)
			auth.Patch("/pets/{id}", petHandler.UpdatePet)
			auth.Delete("/pets/{id}", petHandler.DeletePet)
			auth.Post("/pets/{id}/lost", petHandler.MarkLost)
			auth.Post("/pets/{id}/reunited", petHandler.MarkReunited)

			// QR Tag Management
			auth.Get("/pets/{id}/qr", qrHandler.ListQRCodes)
			auth.Post("/pets/{id}/qr/generate", qrHandler.GenerateQRCode)
			auth.Post("/qr/{id}/revoke", qrHandler.RevokeQRCode)

			// Health & Clinical Records
			auth.Get("/pets/{id}/vaccinations", healthHandler.ListVaccinations)
			auth.Post("/pets/{id}/vaccinations", healthHandler.CreateVaccination)
			auth.Get("/pets/{id}/medical-events", healthHandler.ListMedicalEvents)
			auth.Post("/pets/{id}/medical-events", healthHandler.CreateMedicalEvent)
			auth.Get("/pets/{id}/medications", healthHandler.ListMedications)
			auth.Post("/pets/{id}/medications", healthHandler.CreateMedication)
			auth.Get("/pets/{id}/timeline", healthHandler.GetTimeline)

			// Documents
			auth.Get("/pets/{id}/documents", healthHandler.ListDocuments)
			auth.Post("/pets/{id}/documents", healthHandler.SaveDocument)
			auth.Post("/storage/presigned-url", storageHandler.GeneratePresignedURL)

			// Reminders
			auth.Get("/reminders", reminderHandler.ListReminders)
			auth.Post("/reminders", reminderHandler.CreateReminder)
			auth.Post("/reminders/{id}/complete", reminderHandler.CompleteReminder)
		})
	})

	// 7. Background Outbox Publisher
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	publisher := outbox.NewEventPublisher(outboxRepo, 2*time.Second)
	go publisher.Start(ctx)

	// 8. Server Lifecycle
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("🐾 Pet Digital Identity Platform API running on %s (ENV=%s)\n", addr, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down API server gracefully...")
	publisher.Stop()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	fmt.Println("Server exited cleanly.")
}
