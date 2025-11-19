package utils

import (
	"encoding/json"
	"net/http"

	appErrors "quote-vault/errors"
	"quote-vault/middleware"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// WriteJSON writes a JSON response to the response writer
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes a successful JSON response
func WriteSuccess(w http.ResponseWriter, data interface{}, message string) error {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	return WriteJSON(w, http.StatusOK, response)
}

// WriteCreated writes a successful creation response
func WriteCreated(w http.ResponseWriter, data interface{}) error {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: "Resource created successfully",
	}
	return WriteJSON(w, http.StatusCreated, response)
}

// WritePaginated writes a paginated response
func WritePaginated(w http.ResponseWriter, data interface{}, pagination Pagination) error {
	response := PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	}
	return WriteJSON(w, http.StatusOK, response)
}

// WriteError writes an error response using the error handler middleware
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	middleware.HandleError(w, r, err)
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, r *http.Request, message string, err error) {
	appErr := appErrors.NewValidationError(message, err)
	WriteError(w, r, appErr)
}

// WriteNotFoundError writes a not found error response
func WriteNotFoundError(w http.ResponseWriter, r *http.Request, resource string) {
	appErr := appErrors.NewNotFoundError(resource)
	WriteError(w, r, appErr)
}

// WriteInternalError writes an internal server error response
func WriteInternalError(w http.ResponseWriter, r *http.Request, message string, err error) {
	appErr := appErrors.NewInternalError(message, err)
	WriteError(w, r, appErr)
}