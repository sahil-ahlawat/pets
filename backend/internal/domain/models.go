package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Roles & Enums
type UserRole string

const (
	RolePetOwner     UserRole = "PET_OWNER"
	RoleVet          UserRole = "VET"
	RoleShelterAdmin UserRole = "SHELTER_ADMIN"
	RoleSystemAdmin  UserRole = "SYSTEM_ADMIN"
)

type PetStatus string

const (
	PetStatusActive   PetStatus = "ACTIVE"
	PetStatusLost     PetStatus = "LOST"
	PetStatusFound    PetStatus = "FOUND"
	PetStatusReunited PetStatus = "REUNITED"
	PetStatusArchived PetStatus = "ARCHIVED"
)

type PetSpecies string

const (
	SpeciesDog     PetSpecies = "DOG"
	SpeciesCat     PetSpecies = "CAT"
	SpeciesBird    PetSpecies = "BIRD"
	SpeciesRabbit  PetSpecies = "RABBIT"
	SpeciesHorse   PetSpecies = "HORSE"
	SpeciesReptile PetSpecies = "REPTILE"
	SpeciesOther   PetSpecies = "OTHER"
)

// User models
type User struct {
	ID               uuid.UUID  `json:"id"`
	Email            string     `json:"email"`
	PasswordHash     string     `json:"-"`
	FirstName        string     `json:"firstName"`
	LastName         string     `json:"lastName"`
	PhoneNumber      string     `json:"phoneNumber,omitempty"`
	PhoneVerified    bool       `json:"phoneVerified"`
	EmailVerified    bool       `json:"emailVerified"`
	Role             UserRole   `json:"role"`
	AvatarURL        string     `json:"avatarUrl,omitempty"`
	Timezone         string     `json:"timezone"`
	AllowDirectCalls bool       `json:"allowDirectCalls"`
	IsActive         bool       `json:"isActive"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"-"`
}

type UserSession struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"userId"`
	RefreshTokenHash string    `json:"-"`
	UserAgent        string    `json:"userAgent"`
	ClientIP         string    `json:"clientIp"`
	ExpiresAt        time.Time `json:"expiresAt"`
	IsRevoked        bool      `json:"isRevoked"`
	CreatedAt        time.Time `json:"createdAt"`
}

type UserDevice struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	DeviceToken string    `json:"deviceToken"`
	Platform    string    `json:"platform"`
	DeviceName  string    `json:"deviceName"`
	AppVersion  string    `json:"appVersion"`
	IsActive    bool      `json:"isActive"`
	LastSeenAt  time.Time `json:"lastSeenAt"`
}

// Pet models
type Pet struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Species            PetSpecies `json:"species"`
	Breed              string     `json:"breed,omitempty"`
	Gender             string     `json:"gender,omitempty"`
	IsNeutered         bool       `json:"isNeutered"`
	DateOfBirth        *string    `json:"dateOfBirth,omitempty"`
	ApproxAgeYears     *int       `json:"approximateAgeYears,omitempty"`
	ColorMarkings      string     `json:"colorMarkings,omitempty"`
	WeightKg           *float64   `json:"weightKg,omitempty"`
	AvatarURL          string     `json:"avatarUrl,omitempty"`
	Status             PetStatus  `json:"status"`
	IsArchived         bool       `json:"isArchived"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"-"`
}

type PetProfile struct {
	ID                     uuid.UUID `json:"id"`
	PetID                  uuid.UUID `json:"petId"`
	Bio                    string    `json:"bio,omitempty"`
	DietaryPreferences     string    `json:"dietaryPreferences,omitempty"`
	FavoriteToysActivities string    `json:"favoriteToysActivities,omitempty"`
	BehavioralNotes        string    `json:"behavioralNotes,omitempty"`
	Allergies              string    `json:"allergies,omitempty"`
	ChronicConditions      string    `json:"chronicConditions,omitempty"`
	ShowOwnerName          bool      `json:"showOwnerName"`
	ShowOwnerPhone         bool      `json:"showOwnerPhone"`
	ShowMicrochip          bool      `json:"showMicrochip"`
	ShowAllergiesEmergency bool      `json:"showAllergiesEmergency"`
	ShowVetInfo            bool      `json:"showVetInfo"`
	CustomMessage          string    `json:"customMessage,omitempty"`
}

type PetEmergencyProfile struct {
	ID                      uuid.UUID `json:"id"`
	PetID                   uuid.UUID `json:"petId"`
	EmergencyContactName    string    `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone   string    `json:"emergencyContactPhone,omitempty"`
	EmergencyContactRelation string   `json:"emergencyContactRelation,omitempty"`
	CriticalMedicalAlert    string    `json:"criticalMedicalAlert,omitempty"`
	SpecialCareInstructions string    `json:"specialCareInstructions,omitempty"`
	InsuranceProvider       string    `json:"insuranceProvider,omitempty"`
	InsurancePolicyNumber   string    `json:"insurancePolicyNumber,omitempty"`
}

type PetIdentification struct {
	ID             uuid.UUID `json:"id"`
	PetID          uuid.UUID `json:"petId"`
	IDType         string    `json:"idType"` // MICROCHIP, TATTOO, MUNICIPAL_LICENSE
	IDValue        string    `json:"idValue"`
	IssuerRegistry string    `json:"issuerRegistry,omitempty"`
	ImplantDate    *string   `json:"implantDate,omitempty"`
	Verified       bool      `json:"verified"`
}

type PetQRCode struct {
	ID               uuid.UUID  `json:"id"`
	PetID            uuid.UUID  `json:"petId"`
	Token            string     `json:"token"`
	Status           string     `json:"status"` // ACTIVE, INACTIVE, REVOKED
	TagLabel         string     `json:"tagLabel"`
	PhysicalTagType  string     `json:"physicalTagType"`
	ScanCount        int        `json:"scanCount"`
	LastScannedAt    *time.Time `json:"lastScannedAt,omitempty"`
	PublicURL        string     `json:"publicUrl"`
	CreatedAt        time.Time  `json:"createdAt"`
}

type QRScan struct {
	ID             uuid.UUID  `json:"id"`
	QRCodeID       uuid.UUID  `json:"qrCodeId"`
	PetID          uuid.UUID  `json:"petId"`
	ScannedAt      time.Time  `json:"scannedAt"`
	ApproxLatitude *float64   `json:"approxLatitude,omitempty"`
	ApproxLongitude *float64  `json:"approxLongitude,omitempty"`
	AccuracyMeters *float64   `json:"accuracyMeters,omitempty"`
	ApproxAddress  string     `json:"approxAddress,omitempty"`
	City           string     `json:"city,omitempty"`
	CountryCode    string     `json:"countryCode,omitempty"`
	IPHash         string     `json:"ipHash,omitempty"`
	UserAgent      string     `json:"userAgent,omitempty"`
}

// Health & Clinical records
type VaccinationType struct {
	ID                       uuid.UUID  `json:"id"`
	Species                  PetSpecies `json:"species"`
	Name                     string     `json:"name"`
	Code                     string     `json:"code"`
	IsCore                   bool       `json:"isCore"`
	RecommendedIntervalMonths int       `json:"recommendedIntervalMonths"`
	Description              string     `json:"description,omitempty"`
}

type PetVaccination struct {
	ID                uuid.UUID  `json:"id"`
	PetID             uuid.UUID  `json:"petId"`
	VaccinationTypeID *uuid.UUID `json:"vaccinationTypeId,omitempty"`
	VaccineName       string     `json:"vaccineName"`
	AdministeredAt    string     `json:"administeredAt"`
	NextDueDate       *string    `json:"nextDueDate,omitempty"`
	VetName           string     `json:"vetName,omitempty"`
	ClinicName        string     `json:"clinicName,omitempty"`
	BatchLotNumber    string     `json:"batchLotNumber,omitempty"`
	Manufacturer      string     `json:"manufacturer,omitempty"`
	DocumentID        *uuid.UUID `json:"documentId,omitempty"`
	Notes             string     `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
}

type PetMedicalEvent struct {
	ID               uuid.UUID  `json:"id"`
	PetID            uuid.UUID  `json:"petId"`
	EventType        string     `json:"eventType"`
	EventDate        string     `json:"eventDate"`
	Title            string     `json:"title"`
	Diagnosis        string     `json:"diagnosis,omitempty"`
	TreatmentSummary string     `json:"treatmentSummary,omitempty"`
	VetName          string     `json:"vetName,omitempty"`
	ClinicName       string     `json:"clinicName,omitempty"`
	CostAmount       *float64   `json:"costAmount,omitempty"`
	CostCurrency     string     `json:"costCurrency,omitempty"`
	DocumentID       *uuid.UUID `json:"documentId,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
}

type PetMedication struct {
	ID             uuid.UUID `json:"id"`
	PetID          uuid.UUID `json:"petId"`
	MedicationName string    `json:"medicationName"`
	Dosage         string    `json:"dosage"`
	Frequency      string    `json:"frequency"`
	StartDate      string    `json:"startDate"`
	EndDate        *string   `json:"endDate,omitempty"`
	PrescribedBy   string    `json:"prescribedBy,omitempty"`
	Reason         string    `json:"reason,omitempty"`
	IsActive       bool      `json:"isActive"`
}

type PetReminder struct {
	ID                uuid.UUID  `json:"id"`
	PetID             uuid.UUID  `json:"petId"`
	PetName           string     `json:"petName,omitempty"`
	UserID            uuid.UUID  `json:"userId"`
	ReminderType      string     `json:"reminderType"`
	Title             string     `json:"title"`
	Description       string     `json:"description,omitempty"`
	DueDate           time.Time  `json:"dueDate"`
	Recurrence        string     `json:"recurrence"`
	IsCompleted       bool       `json:"isCompleted"`
	CompletedAt       *time.Time `json:"completedAt,omitempty"`
}

type PetTimelineEvent struct {
	ID            uuid.UUID `json:"id"`
	PetID         uuid.UUID `json:"petId"`
	EventCategory string    `json:"eventCategory"`
	Timestamp     time.Time `json:"timestamp"`
	Title         string    `json:"title"`
	Description   string    `json:"description,omitempty"`
	IconName      string    `json:"iconName,omitempty"`
}

type PetDocument struct {
	ID            uuid.UUID `json:"id"`
	PetID         uuid.UUID `json:"petId"`
	UploadedBy    uuid.UUID `json:"uploadedBy"`
	Title         string    `json:"title"`
	Category      string    `json:"category"`
	StorageKey    string    `json:"storageKey"`
	FileURL       string    `json:"fileUrl"`
	MimeType      string    `json:"mimeType"`
	FileSizeBytes int64     `json:"fileSizeBytes"`
	IsVerified    bool      `json:"isVerified"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Vet struct {
	ID             uuid.UUID `json:"id"`
	ClinicName     string    `json:"clinicName"`
	DoctorName     string    `json:"doctorName,omitempty"`
	LicenseNumber  string    `json:"licenseNumber,omitempty"`
	PhoneNumber    string    `json:"phoneNumber"`
	EmergencyPhone string    `json:"emergencyPhone,omitempty"`
	Email          string    `json:"email,omitempty"`
	WebsiteURL     string    `json:"websiteUrl,omitempty"`
	AddressLine1   string    `json:"addressLine1"`
	City           string    `json:"city"`
	State          string    `json:"state,omitempty"`
	PostalCode     string    `json:"postalCode,omitempty"`
	CountryCode    string    `json:"countryCode"`
	Latitude       *float64  `json:"latitude,omitempty"`
	Longitude      *float64  `json:"longitude,omitempty"`
	Is24hEmergency bool      `json:"is24hEmergency"`
	Specializations []string `json:"specializations,omitempty"`
	Rating         float64   `json:"rating"`
	ReviewCount    int       `json:"reviewCount"`
	DistanceKm     *float64  `json:"distanceKm,omitempty"`
}

type Notification struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"userId"`
	PetID            *uuid.UUID `json:"petId,omitempty"`
	NotificationType string    `json:"notificationType"`
	Title            string    `json:"title"`
	Body             string    `json:"body"`
	IsRead           bool      `json:"isRead"`
	CreatedAt        time.Time `json:"createdAt"`
}

// Outbox and Jobs
type OutboxEvent struct {
	ID            uuid.UUID       `json:"id"`
	EventType     string          `json:"eventType"`
	AggregateType string          `json:"aggregateType"`
	AggregateID   uuid.UUID       `json:"aggregateId"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"createdAt"`
	IsPublished   bool            `json:"isPublished"`
}

type Job struct {
	ID         uuid.UUID       `json:"id"`
	Queue      string          `json:"queue"`
	JobType    string          `json:"jobType"`
	Payload    json.RawMessage `json:"payload"`
	Status     string          `json:"status"`
	RunAt      time.Time       `json:"runAt"`
	RetryCount int             `json:"retryCount"`
}

// Request and Response DTOs
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type AuthResponse struct {
	Tokens AuthTokens `json:"tokens"`
	User   User       `json:"user"`
}

type CreatePetRequest struct {
	Name                  string     `json:"name"`
	Species               PetSpecies `json:"species"`
	Breed                 string     `json:"breed"`
	Gender                string     `json:"gender"`
	IsNeutered            bool       `json:"isNeutered"`
	DateOfBirth           *string    `json:"dateOfBirth"`
	WeightKg              *float64   `json:"weightKg"`
	AvatarURL             string     `json:"avatarUrl"`
	Bio                   string     `json:"bio"`
	Allergies             string     `json:"allergies"`
	EmergencyContactName  string     `json:"emergencyContactName"`
	EmergencyContactPhone string     `json:"emergencyContactPhone"`
	MicrochipNumber       string     `json:"microchipNumber"`
}

type UpdatePetRequest struct {
	Name                   *string  `json:"name"`
	Breed                  *string  `json:"breed"`
	Gender                 *string  `json:"gender"`
	IsNeutered             *bool    `json:"isNeutered"`
	WeightKg               *float64 `json:"weightKg"`
	AvatarURL              *string  `json:"avatarUrl"`
	Bio                    *string  `json:"bio"`
	Allergies              *string  `json:"allergies"`
	ChronicConditions      *string  `json:"chronicConditions"`
	EmergencyContactName   *string  `json:"emergencyContactName"`
	EmergencyContactPhone  *string  `json:"emergencyContactPhone"`
	CriticalMedicalAlert   *string  `json:"criticalMedicalAlert"`
	ShowOwnerName          *bool    `json:"showOwnerName"`
	ShowOwnerPhone         *bool    `json:"showOwnerPhone"`
	ShowMicrochip          *bool    `json:"showMicrochip"`
	ShowAllergiesEmergency *bool    `json:"showAllergiesEmergency"`
	ShowVetInfo            *bool    `json:"showVetInfo"`
}

type PetDetail struct {
	Pet
	Bio                    string  `json:"bio"`
	DietaryPreferences     string  `json:"dietaryPreferences"`
	Allergies              string  `json:"allergies"`
	ChronicConditions      string  `json:"chronicConditions"`
	EmergencyContactName   string  `json:"emergencyContactName"`
	EmergencyContactPhone  string  `json:"emergencyContactPhone"`
	CriticalMedicalAlert   string  `json:"criticalMedicalAlert"`
	MicrochipNumber        string  `json:"microchipNumber"`
	ActiveQRToken          string  `json:"activeQrToken"`
	QRURL                  string  `json:"qrUrl"`
	ShowOwnerName          bool    `json:"showOwnerName"`
	ShowOwnerPhone         bool    `json:"showOwnerPhone"`
	ShowMicrochip          bool    `json:"showMicrochip"`
	ShowAllergiesEmergency bool    `json:"showAllergiesEmergency"`
	ShowVetInfo            bool    `json:"showVetInfo"`
}

type PublicPetProfile struct {
	Token                string      `json:"token"`
	Name                 string      `json:"name"`
	Species              string      `json:"species"`
	Breed                string      `json:"breed,omitempty"`
	Gender               string      `json:"gender,omitempty"`
	AgeDescription       string      `json:"ageDescription,omitempty"`
	AvatarURL            string      `json:"avatarUrl,omitempty"`
	Status               PetStatus   `json:"status"`
	LostDetails          *LostDetails `json:"lostDetails,omitempty"`
	OwnerName            string      `json:"ownerName,omitempty"`
	OwnerPhone           string      `json:"ownerPhone,omitempty"`
	AllowDirectCalls     bool        `json:"allowDirectCalls"`
	MicrochipNumber      string      `json:"microchipNumber,omitempty"`
	Allergies            string      `json:"allergies,omitempty"`
	CriticalMedicalAlert string      `json:"criticalMedicalAlert,omitempty"`
	SpecialInstructions  string      `json:"specialInstructions,omitempty"`
	EmergencyContactName string      `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone string     `json:"emergencyContactPhone,omitempty"`
}

type LostDetails struct {
	LastSeenAddress string    `json:"lastSeenAddress"`
	LastSeenCity    string    `json:"lastSeenCity,omitempty"`
	Circumstances   string    `json:"circumstances,omitempty"`
	RewardAmount    string    `json:"rewardAmount,omitempty"`
	LostDate        time.Time `json:"lostDate"`
}

type ContactOwnerRequest struct {
	FinderName      string   `json:"finderName"`
	FinderPhone     string   `json:"finderPhone"`
	FinderEmail     string   `json:"finderEmail"`
	Message         string   `json:"message"`
	ApproxLatitude  *float64 `json:"approxLatitude"`
	ApproxLongitude *float64 `json:"approxLongitude"`
}

type ReportFoundRequest struct {
	FinderName   string   `json:"finderName"`
	FinderPhone  string   `json:"finderPhone"`
	Message      string   `json:"message"`
	FoundAddress string   `json:"foundAddress"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	PhotoURL     string   `json:"photoUrl"`
}

type MarkLostRequest struct {
	LastSeenAddress string `json:"lastSeenAddress"`
	LastSeenCity    string `json:"lastSeenCity"`
	Circumstances   string `json:"circumstances"`
	RewardAmount    string `json:"rewardAmount"`
}

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
