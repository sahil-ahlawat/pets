package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"pet-platform/backend/internal/service"
)

type QRHandler struct {
	qrService service.QRService
}

func NewQRHandler(qrService service.QRService) *QRHandler {
	return &QRHandler{qrService: qrService}
}

func (h *QRHandler) ListQRCodes(w http.ResponseWriter, r *http.Request) {
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

	codes, err := h.qrService.GetQRCodes(r.Context(), userID, petID)
	if err != nil {
		writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, codes)
}

func (h *QRHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
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

	var body struct {
		TagLabel string `json:"tagLabel"`
	}
	_ = readJSON(r, &body)

	code, err := h.qrService.GenerateNewQRCode(r.Context(), userID, petID, body.TagLabel)
	if err != nil {
		writeError(w, http.StatusBadRequest, "GENERATE_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, code)
}

func (h *QRHandler) RevokeQRCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	qrID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid QR code ID")
		return
	}

	var body struct {
		Reason string `json:"reason"`
	}
	_ = readJSON(r, &body)

	if err := h.qrService.RevokeQRCode(r.Context(), userID, qrID, body.Reason); err != nil {
		writeError(w, http.StatusBadRequest, "REVOKE_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, "QR Code revoked successfully")
}

func (h *QRHandler) RenderQRCodePNG(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	size := 512
	if s := r.URL.Query().Get("size"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil {
			size = parsed
		}
	}

	pngBytes, err := h.qrService.GenerateQRCodePNG(token, size)
	if err != nil {
		http.Error(w, "Failed to render QR", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngBytes)
}
