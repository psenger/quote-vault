package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Type       string `json:"type"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewValidationError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Type:       "validation_error",
	}
}

func NewNotFoundError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusNotFound,
		Type:       "not_found_error",
	}
}

func NewDatabaseError(message string) *APIError {
	return &APIError{
		Message:    fmt.Sprintf("Database error: %s", message),
		StatusCode: http.StatusInternalServerError,
		Type:       "database_error",
	}
}

func NewInternalServerError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Type:       "internal_server_error",
	}
}