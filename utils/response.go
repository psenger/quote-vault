package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	Status    int         `json:"status"`
}

type ErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Path      string `json:"path,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    status,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding success response: %v", err)
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := ErrorResponse{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    status,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
		// Fallback to basic error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ValidationErrorResponse(w http.ResponseWriter, errors []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := map[string]interface{}{
		"success":          false,
		"error":            "Validation failed",
		"validation_errors": errors,
		"timestamp":        time.Now().UTC().Format(time.RFC3339),
		"status":           http.StatusBadRequest,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding validation error response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}