package handler

import (
	"context"
	"net/http"
	"time"

	"pet-platform/backend/internal/repository"
)

type HealthzHandler struct {
	db *repository.DB
}

func NewHealthzHandler(db *repository.DB) *HealthzHandler {
	return &HealthzHandler{db: db}
}

func (h *HealthzHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"UP"}`))
}

func (h *HealthzHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		writeError(w, http.StatusServiceUnavailable, "DEPENDENCY_DOWN", "Database ping failed")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"READY","database":"UP"}`))
}

func (h *HealthzHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("# HELP pet_platform_api_up Service up indicator\n# TYPE pet_platform_api_up gauge\npet_platform_api_up 1\n"))
}
