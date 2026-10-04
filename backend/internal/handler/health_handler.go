package handler

import (
	"net/http"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/service"
)

type HealthHandler struct {
	healthService service.HealthService
}

func NewHealthHandler(healthService service.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) ListVaccinationTypes(w http.ResponseWriter, r *http.Request) {
	species := r.URL.Query().Get("species")
	types, err := h.healthService.ListVaccinationTypes(r.Context(), species)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, types)
}

func (h *HealthHandler) ListVaccinations(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	vacs, err := h.healthService.ListVaccinations(r.Context(), userID, petID)
	if err != nil {
		writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, vacs)
}

func (h *HealthHandler) CreateVaccination(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	var req struct {
		VaccineName    string  `json:"vaccineName"`
		AdministeredAt string  `json:"administeredAt"`
		NextDueDate    *string `json:"nextDueDate"`
		VetName        string  `json:"vetName"`
		ClinicName     string  `json:"clinicName"`
		BatchLotNumber string  `json:"batchLotNumber"`
		Manufacturer   string  `json:"manufacturer"`
		Notes          string  `json:"notes"`
	}
	if err := readJSON(r, &req); err != nil || req.VaccineName == "" || req.AdministeredAt == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "vaccineName and administeredAt are required")
		return
	}

	vac := &domain.PetVaccination{
		VaccineName:    req.VaccineName,
		AdministeredAt: req.AdministeredAt,
		NextDueDate:    req.NextDueDate,
		VetName:        req.VetName,
		ClinicName:     req.ClinicName,
		BatchLotNumber: req.BatchLotNumber,
		Manufacturer:   req.Manufacturer,
		Notes:          req.Notes,
	}

	if err := h.healthService.CreateVaccination(r.Context(), userID, petID, vac); err != nil {
		writeError(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, vac)
}

func (h *HealthHandler) ListMedicalEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	events, err := h.healthService.ListMedicalEvents(r.Context(), userID, petID)
	if err != nil {
		writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, events)
}

func (h *HealthHandler) CreateMedicalEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	var req domain.PetMedicalEvent
	if err := readJSON(r, &req); err != nil || req.Title == "" || req.EventDate == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "title and eventDate are required")
		return
	}

	if err := h.healthService.CreateMedicalEvent(r.Context(), userID, petID, &req); err != nil {
		writeError(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, req)
}

func (h *HealthHandler) ListMedications(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	meds, err := h.healthService.ListMedications(r.Context(), userID, petID)
	if err != nil {
		writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, meds)
}

func (h *HealthHandler) CreateMedication(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	var req domain.PetMedication
	if err := readJSON(r, &req); err != nil || req.MedicationName == "" || req.Dosage == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "medicationName and dosage are required")
		return
	}

	if err := h.healthService.CreateMedication(r.Context(), userID, petID, &req); err != nil {
		writeError(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, req)
}

func (h *HealthHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	events, err := h.healthService.GetTimeline(r.Context(), userID, petID, 50)
	if err != nil {
		writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, events)
}

func (h *HealthHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	docs, err := h.healthService.ListDocuments(r.Context(), userID, petID)
	if err != nil {
		writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, docs)
}

func (h *HealthHandler) SaveDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	petID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid pet ID")
		return
	}

	var req struct {
		Title         string `json:"title"`
		Category      string `json:"category"`
		StorageKey    string `json:"storageKey"`
		FileURL       string `json:"fileUrl"`
		MimeType      string `json:"mimeType"`
		FileSizeBytes int64  `json:"fileSizeBytes"`
	}
	if err := readJSON(r, &req); err != nil || req.Title == "" || req.StorageKey == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "title and storageKey are required")
		return
	}

	doc := &domain.PetDocument{
		ID:            uuid.New(),
		PetID:         petID,
		UploadedBy:    userID,
		Title:         req.Title,
		Category:      req.Category,
		StorageKey:    req.StorageKey,
		FileURL:       req.FileURL,
		MimeType:      req.MimeType,
		FileSizeBytes: req.FileSizeBytes,
		IsVerified:    true,
	}

	if err := h.healthService.CreateDocument(r.Context(), userID, petID, doc); err != nil {
		writeError(w, http.StatusBadRequest, "SAVE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, doc)
}
