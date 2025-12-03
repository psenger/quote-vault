package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with additional context
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Error types
const (
	TypeValidation   = "validation_error"
	TypeNotFound     = "not_found"
	TypeDatabase     = "database_error"
	TypeInternal     = "internal_error"
	TypeBadRequest   = "bad_request"
	TypeConflict     = "conflict"
)

// Common errors
var (
	ErrQuoteNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Quote not found",
		Type:    TypeNotFound,
	}

	ErrEmptyQuoteText = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Quote text cannot be empty",
		Type:    TypeValidation,
	}

	ErrEmptyAuthor = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Author cannot be empty",
		Type:    TypeValidation,
	}

	ErrInvalidID = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid quote ID",
		Type:    TypeValidation,
	}

	ErrInvalidCategory = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid category specified",
		Type:    TypeBadRequest,
	}

	ErrInvalidPagination = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid pagination parameters",
		Type:    TypeBadRequest,
	}

	ErrDatabaseConnection = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Database connection failed",
		Type:    TypeDatabase,
	}

	ErrQuoteExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Quote already exists",
		Type:    TypeConflict,
	}
)

// NewValidationError creates a new validation error with specific details
func NewValidationError(message, detail string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Type:    TypeValidation,
		Detail:  detail,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Type:    TypeDatabase,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Type:    TypeInternal,
	}
}