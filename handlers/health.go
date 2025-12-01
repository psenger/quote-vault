package handlers

import (
	"net/http"
	"time"

	"quote-vault/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "quote-vault",
		"version":   "1.0.0",
	}

	utils.SuccessResponse(w, http.StatusOK, health)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Here you could add database connectivity checks
	// For now, just return ready status
	ready := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks": map[string]string{
			"database": "ok",
		},
	}

	utils.SuccessResponse(w, http.StatusOK, ready)
}