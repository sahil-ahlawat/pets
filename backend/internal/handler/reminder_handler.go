package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/service"
)

type ReminderHandler struct {
	reminderService service.ReminderService
}

func NewReminderHandler(reminderService service.ReminderService) *ReminderHandler {
	return &ReminderHandler{reminderService: reminderService}
}

func (h *ReminderHandler) ListReminders(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	reminders, err := h.reminderService.ListUserReminders(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, reminders)
}

func (h *ReminderHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	var req struct {
		PetID        string `json:"petId"`
		ReminderType string `json:"reminderType"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		DueDate      string `json:"dueDate"`
		Recurrence   string `json:"recurrence"`
	}
	if err := readJSON(r, &req); err != nil || req.PetID == "" || req.Title == "" || req.DueDate == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "petId, title, and dueDate are required")
		return
	}

	petUUID, err := uuid.Parse(req.PetID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid petId format")
		return
	}

	dueTime, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		// Try date format YYYY-MM-DD
		dueTime, err = time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_DATE", "dueDate must be RFC3339 or YYYY-MM-DD")
			return
		}
	}

	reminder := &domain.PetReminder{
		PetID:        petUUID,
		ReminderType: req.ReminderType,
		Title:        req.Title,
		Description:  req.Description,
		DueDate:      dueTime,
		Recurrence:   req.Recurrence,
	}

	if err := h.reminderService.CreateReminder(r.Context(), userID, reminder); err != nil {
		writeError(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, reminder)
}

func (h *ReminderHandler) CompleteReminder(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	remID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid reminder ID")
		return
	}

	if err := h.reminderService.MarkReminderCompleted(r.Context(), remID, userID); err != nil {
		writeError(w, http.StatusBadRequest, "COMPLETE_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, "Reminder marked as completed")
}
