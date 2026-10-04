package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
	"pet-platform/backend/internal/middleware"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

func writeSuccess(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, domain.StandardResponse{
		Success: true,
		Message: message,
	})
}

func readJSON(r *http.Request, dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func parseUUIDParam(r *http.Request, paramName string) (uuid.UUID, error) {
	str := chi.URLParam(r, paramName)
	return uuid.Parse(str)
}

func getUserID(r *http.Request) (uuid.UUID, bool) {
	return middleware.GetUserID(r.Context())
}
