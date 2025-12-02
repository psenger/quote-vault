package handlers

import (
	"net/http"
	"quote-vault/database"
	"quote-vault/utils"
)

type HealthHandler struct {
	db database.Database
}

func NewHealthHandler(db database.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	}

	// Test database connection
	if err := h.db.Ping(); err != nil {
		response.Database = "unhealthy"
		response.Status = "degraded"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		response.Database = "healthy"
	}

	utils.WriteJSONResponse(w, response)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Check if database is accessible
	if err := h.db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		utils.WriteJSONResponse(w, map[string]string{
			"status": "not ready",
			"reason": "database connection failed",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteJSONResponse(w, map[string]string{"status": "ready"})
}