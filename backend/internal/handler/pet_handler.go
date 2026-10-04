package handler

import (
	"net/http"
	"strconv"

	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/service"
)

type PetHandler struct {
	petService service.PetService
}

func NewPetHandler(petService service.PetService) *PetHandler {
	return &PetHandler{petService: petService}
}

func (h *PetHandler) ListPets(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	cursor := r.URL.Query().Get("cursor")

	pets, nextCursor, err := h.petService.ListPets(r.Context(), userID, limit, cursor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pets":       pets,
		"nextCursor": nextCursor,
	})
}

func (h *PetHandler) CreatePet(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	var req domain.CreatePetRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	pet, err := h.petService.CreatePet(r.Context(), userID, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, pet)
}

func (h *PetHandler) GetPet(w http.ResponseWriter, r *http.Request) {
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

	pet, err := h.petService.GetPet(r.Context(), userID, petID)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, pet)
}

func (h *PetHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
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

	var req domain.UpdatePetRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	pet, err := h.petService.UpdatePet(r.Context(), userID, petID, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, pet)
}

func (h *PetHandler) DeletePet(w http.ResponseWriter, r *http.Request) {
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

	if err := h.petService.DeletePet(r.Context(), userID, petID); err != nil {
		writeError(w, http.StatusBadRequest, "DELETE_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, "Pet archived successfully")
}

func (h *PetHandler) MarkLost(w http.ResponseWriter, r *http.Request) {
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

	var req domain.MarkLostRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	if err := h.petService.MarkPetLost(r.Context(), userID, petID, &req); err != nil {
		writeError(w, http.StatusBadRequest, "MARK_LOST_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, "Pet marked as LOST. Broadcast activated.")
}

func (h *PetHandler) MarkReunited(w http.ResponseWriter, r *http.Request) {
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

	if err := h.petService.MarkPetReunited(r.Context(), userID, petID); err != nil {
		writeError(w, http.StatusBadRequest, "MARK_REUNITED_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, "Pet successfully marked as reunited! Emergency mode deactivated.")
}
