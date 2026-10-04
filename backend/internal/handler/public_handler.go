package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/service"
)

type PublicHandler struct {
	qrService service.QRService
}

func NewPublicHandler(qrService service.QRService) *PublicHandler {
	return &PublicHandler{qrService: qrService}
}

func (h *PublicHandler) GetPublicPet(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Token required")
		return
	}

	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}
	userAgent := r.UserAgent()

	var latPtr, lngPtr *float64
	if latStr := r.URL.Query().Get("lat"); latStr != "" {
		if val, err := strconv.ParseFloat(latStr, 64); err == nil {
			latPtr = &val
		}
	}
	if lngStr := r.URL.Query().Get("lng"); lngStr != "" {
		if val, err := strconv.ParseFloat(lngStr, 64); err == nil {
			lngPtr = &val
		}
	}

	profile, err := h.qrService.GetPublicProfile(r.Context(), token, clientIP, userAgent, latPtr, lngPtr)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

func (h *PublicHandler) ContactOwner(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Token required")
		return
	}

	var req domain.ContactOwnerRequest
	if err := readJSON(r, &req); err != nil || req.Message == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Message is required")
		return
	}

	if err := h.qrService.ContactOwner(r.Context(), token, &req); err != nil {
		writeError(w, http.StatusBadRequest, "CONTACT_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, "Message sent to pet owner securely")
}

func (h *PublicHandler) ReportFound(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Token required")
		return
	}

	var req domain.ReportFoundRequest
	if err := readJSON(r, &req); err != nil || req.Message == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Message is required")
		return
	}

	if err := h.qrService.ReportFound(r.Context(), token, &req); err != nil {
		writeError(w, http.StatusBadRequest, "REPORT_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, "Found pet report submitted. Owner has been alerted immediately.")
}
