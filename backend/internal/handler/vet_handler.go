package handler

import (
	"net/http"
	"strconv"

	"pet-platform/backend/internal/service"
)

type VetHandler struct {
	vetService service.VetService
}

func NewVetHandler(vetService service.VetService) *VetHandler {
	return &VetHandler{vetService: vetService}
}

func (h *VetHandler) SearchVets(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	city := r.URL.Query().Get("city")
	emerg := r.URL.Query().Get("emergencyOnly") == "true"

	vets, err := h.vetService.SearchVets(r.Context(), q, city, emerg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, vets)
}

func (h *VetHandler) FindNearbyVets(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("latitude")
	lngStr := r.URL.Query().Get("longitude")
	if latStr == "" || lngStr == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "latitude and longitude query parameters are required")
		return
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)
	if err1 != nil || err2 != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid latitude or longitude format")
		return
	}

	radiusKm := 25.0
	if rStr := r.URL.Query().Get("radiusKm"); rStr != "" {
		if val, err := strconv.ParseFloat(rStr, 64); err == nil {
			radiusKm = val
		}
	}

	vets, err := h.vetService.FindNearbyVets(r.Context(), lat, lng, radiusKm)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, vets)
}
