package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	appErrors "quote-vault/errors"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Type    string `json:"type,omitempty"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// ErrorHandler middleware handles application errors and converts them to HTTP responses
func ErrorHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a custom response writer to capture errors
			cw := &customWriter{
				ResponseWriter: w,
				request:        r,
			}

			next.ServeHTTP(cw, r)
		})
	}
}

// customWriter wraps the response writer to handle panics
type customWriter struct {
	http.ResponseWriter
	request *http.Request
}

// HandleError processes application errors and sends appropriate HTTP responses
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	// Check if it's our custom AppError
	if appErr, ok := err.(*appErrors.AppError); ok {
		response := ErrorResponse{
			Error:   appErr.Message,
			Type:    string(appErr.Type),
			Code:    appErr.Code,
			Message: appErr.Message,
		}

		// Log internal errors
		if appErr.Type == appErrors.ErrorTypeInternal || appErr.Type == appErrors.ErrorTypeDatabase {
			log.Printf("Internal error on %s %s: %v", r.Method, r.URL.Path, appErr)
		}

		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle generic errors
	log.Printf("Unhandled error on %s %s: %v", r.Method, r.URL.Path, err)
	response := ErrorResponse{
		Error:   "Internal server error",
		Type:    string(appErrors.ErrorTypeInternal),
		Code:    http.StatusInternalServerError,
		Message: "An unexpected error occurred",
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}

// RecoverPanic recovers from panics and converts them to internal server errors
func RecoverPanic() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Printf("Panic recovered on %s %s: %v", r.Method, r.URL.Path, rec)
					
					err := appErrors.NewInternalError("Server panic occurred", nil)
					HandleError(w, r, err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}