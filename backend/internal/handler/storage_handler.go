package handler

import (
	"net/http"

	"pet-platform/backend/internal/storage"
)

type StorageHandler struct {
	storageService storage.StorageService
}

func NewStorageHandler(storageService storage.StorageService) *StorageHandler {
	return &StorageHandler{storageService: storageService}
}

func (h *StorageHandler) GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Filename string `json:"filename"`
		MimeType string `json:"mimeType"`
		Category string `json:"category"`
	}
	if err := readJSON(r, &body); err != nil || body.Filename == "" || body.MimeType == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "filename and mimeType are required")
		return
	}

	if body.Category == "" {
		body.Category = "PET_PHOTO"
	}

	res, err := h.storageService.GeneratePresignedUploadURL(body.Filename, body.MimeType, body.Category)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "STORAGE_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)
}
