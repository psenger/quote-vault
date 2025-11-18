package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/quote-vault/database"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Message  string `json:"message"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status:  "ok",
		Message: "Quote Vault API is running",
	}

	// Check database health
	if err := database.HealthCheck(); err != nil {
		response.Status = "error"
		response.Database = "disconnected"
		response.Message = "Database connection failed"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		response.Database = "connected"
	}

	json.NewEncoder(w).Encode(response)
}