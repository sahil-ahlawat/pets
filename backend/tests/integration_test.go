package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/config"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/repository"
	"pet-platform/backend/internal/service"
)

func TestEndToEndPlatformFlow(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping live DB test")
	}

	cfg := &config.Config{
		DatabaseURL:        dbURL,
		JWTSecret:          "test-secret-key-998811",
		JWTAccessExpiryHr:  24,
		JWTRefreshExpiryHr: 720,
		PublicBaseURL:      "http://localhost:3000",
	}

	db, err := repository.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	ctx := context.Background()
	userRepo := repository.NewUserRepository(db)
	petRepo := repository.NewPetRepository(db)
	healthRepo := repository.NewHealthRepository(db)
	qrRepo := repository.NewQRRepository(db)
	reminderRepo := repository.NewReminderRepository(db)

	authSvc := service.NewAuthService(cfg, userRepo)
	petSvc := service.NewPetService(petRepo)
	healthSvc := service.NewHealthService(healthRepo, petRepo)
	qrSvc := service.NewQRService(cfg, qrRepo, petRepo, userRepo)
	reminderSvc := service.NewReminderService(reminderRepo, petRepo)

	// 1. User Registration
	testEmail := "owner_" + uuid.New().String()[:8] + "@example.com"
	regRes, err := authSvc.Register(ctx, &domain.RegisterRequest{
		Email:       testEmail,
		Password:    "Password123!",
		FirstName:   "Sarah",
		LastName:    "Jenkins",
		PhoneNumber: "+15551234567",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	userID := regRes.User.ID

	// 2. User Login
	loginRes, err := authSvc.Login(ctx, &domain.LoginRequest{
		Email:    testEmail,
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginRes.Tokens.AccessToken == "" {
		t.Fatalf("Expected access token")
	}

	// 3. Create Companion
	dob := "2022-04-10"
	weight := 28.5
	pet, err := petSvc.CreatePet(ctx, userID, &domain.CreatePetRequest{
		Name:                  "Bruno",
		Species:               domain.SpeciesDog,
		Breed:                 "Golden Retriever",
		Gender:                "MALE",
		IsNeutered:            true,
		DateOfBirth:           &dob,
		WeightKg:              &weight,
		Bio:                   "Loves tennis balls and swimming.",
		Allergies:             "No known allergies",
		EmergencyContactName:  "Dr. Elena",
		EmergencyContactPhone: "+15559876543",
		MicrochipNumber:       "985141002938491",
	})
	if err != nil {
		t.Fatalf("CreatePet failed: %v", err)
	}
	if pet.Name != "Bruno" || pet.ActiveQRToken == "" {
		t.Fatalf("Invalid pet created: %+v", pet)
	}

	// 4. List Pets
	pets, _, err := petSvc.ListPets(ctx, userID, 10, "")
	if err != nil || len(pets) == 0 {
		t.Fatalf("ListPets failed: %v", err)
	}

	// 5. Public QR Profile Scan (Simulate phone camera scan)
	pub, err := qrSvc.GetPublicProfile(ctx, pet.ActiveQRToken, "192.168.1.1", "Mobile Safari", nil, nil)
	if err != nil {
		t.Fatalf("GetPublicProfile failed: %v", err)
	}
	if pub.Name != "Bruno" || pub.Status != domain.PetStatusActive {
		t.Fatalf("Public pet profile mismatch: %+v", pub)
	}

	// 6. Record Vaccination
	nextDue := "2027-04-10"
	err = healthSvc.CreateVaccination(ctx, userID, pet.ID, &domain.PetVaccination{
		VaccineName:    "Rabies 3-Year",
		AdministeredAt: "2024-04-10",
		NextDueDate:    &nextDue,
		ClinicName:     "MetroVet Hospital",
		VetName:        "Dr. Elena",
		BatchLotNumber: "LOT-99238",
	})
	if err != nil {
		t.Fatalf("CreateVaccination failed: %v", err)
	}

	vacs, err := healthSvc.ListVaccinations(ctx, userID, pet.ID)
	if err != nil || len(vacs) == 0 {
		t.Fatalf("ListVaccinations failed: %v", err)
	}

	// 7. Schedule Reminder
	err = reminderSvc.CreateReminder(ctx, userID, &domain.PetReminder{
		PetID:        pet.ID,
		ReminderType: "FLEA_TICK",
		Title:        "Monthly NexGard",
		DueDate:      time.Now().Add(48 * time.Hour),
		Recurrence:   "MONTHLY",
	})
	if err != nil {
		t.Fatalf("CreateReminder failed: %v", err)
	}

	// 8. Mark Pet as LOST (Emergency Mode)
	err = petSvc.MarkPetLost(ctx, userID, pet.ID, &domain.MarkLostRequest{
		LastSeenAddress: "Central Park West, New York",
		RewardAmount:    "$500",
		Circumstances:   "Slipped harness near fountain",
	})
	if err != nil {
		t.Fatalf("MarkPetLost failed: %v", err)
	}

	// Verify Public Profile reflects LOST state
	lostPub, err := qrSvc.GetPublicProfile(ctx, pet.ActiveQRToken, "192.168.1.2", "Chrome Android", nil, nil)
	if err != nil || lostPub.Status != domain.PetStatusLost {
		t.Fatalf("Expected public profile to show LOST: %+v", lostPub)
	}

	// 9. Finder reports pet found
	lat := 40.785091
	lng := -73.968285
	err = qrSvc.ReportFound(ctx, pet.ActiveQRToken, &domain.ReportFoundRequest{
		FinderName:   "Alex Rivera",
		FinderPhone:  "+15553334444",
		FoundAddress: "Near 72nd St entrance",
		Latitude:     &lat,
		Longitude:    &lng,
		Message:      "I have Bruno on a leash here!",
	})
	if err != nil {
		t.Fatalf("ReportFound failed: %v", err)
	}

	// 10. Mark Pet as REUNITED
	err = petSvc.MarkPetReunited(ctx, userID, pet.ID)
	if err != nil {
		t.Fatalf("MarkPetReunited failed: %v", err)
	}

	reunitedPub, err := qrSvc.GetPublicProfile(ctx, pet.ActiveQRToken, "192.168.1.1", "Safari", nil, nil)
	if err != nil || reunitedPub.Status != domain.PetStatusActive {
		t.Fatalf("Expected active status after reunite: %+v", reunitedPub)
	}

	// 11. Verify QR Scan Count incremented
	qrs, err := qrSvc.GetQRCodes(ctx, userID, pet.ID)
	if err != nil || len(qrs) == 0 {
		t.Fatalf("Failed to fetch pet QR codes: %v", err)
	}
	if qrs[0].ScanCount < 1 {
		t.Fatalf("Expected scan count >= 1, got %d", qrs[0].ScanCount)
	}

	// 12. Update Pet Profile
	newBio := "Loves agility training and swimming."
	updatedPet, err := petSvc.UpdatePet(ctx, userID, pet.ID, &domain.UpdatePetRequest{
		Bio: &newBio,
	})
	if err != nil {
		t.Fatalf("UpdatePet failed: %v", err)
	}
	if updatedPet.Bio != newBio {
		t.Fatalf("Expected bio %s, got %s", newBio, updatedPet.Bio)
	}

	// 13. Create Medical Event
	err = healthSvc.CreateMedicalEvent(ctx, userID, pet.ID, &domain.PetMedicalEvent{
		EventType: "CONSULTATION",
		EventDate: "2024-05-15",
		Title:     "Annual Wellness Exam",
		Diagnosis: "Healthy weight, clear eyes and ears",
		VetName:   "Dr. Elena",
	})
	if err != nil {
		t.Fatalf("CreateMedicalEvent failed: %v", err)
	}

	events, err := healthSvc.ListMedicalEvents(ctx, userID, pet.ID)
	if err != nil || len(events) == 0 {
		t.Fatalf("ListMedicalEvents failed: %v", err)
	}

	// 14. Vet Discovery Search
	vetRepo := repository.NewVetRepository(db)
	vetSvc := service.NewVetService(vetRepo)
	vets, err := vetSvc.SearchVets(ctx, "", "", false)
	if err != nil {
		t.Fatalf("SearchVets failed: %v", err)
	}
	t.Logf("Found %d seeded vets in database", len(vets))
}
