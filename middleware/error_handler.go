package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"quote-vault/errors"
)

// ErrorHandler middleware handles application errors and converts them to JSON responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to capture panics and errors
		ew := &errorResponseWriter{
			ResponseWriter: w,
			request:        r,
		}

		// Recover from panics
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				handleError(w, r, errors.NewInternalError("An unexpected error occurred"))
			}
		}()

		next.ServeHTTP(ew, r)
	})
}

// errorResponseWriter wraps http.ResponseWriter to handle errors
type errorResponseWriter struct {
	http.ResponseWriter
	request *http.Request
	written bool
}

// WriteError writes an error response
func (ew *errorResponseWriter) WriteError(err error) {
	if ew.written {
		return
	}
	handleError(ew.ResponseWriter, ew.request, err)
	ew.written = true
}

// handleError processes and responds to errors
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	// Log the error with request context
	log.Printf("Error handling request %s %s: %v", r.Method, r.URL.Path, err)

	// Handle different error types
	var appErr *errors.AppError
	switch e := err.(type) {
	case *errors.AppError:
		appErr = e
	case error:
		// Convert generic errors to internal server errors
		appErr = errors.NewInternalError("An internal error occurred")
	default:
		appErr = errors.NewInternalError("Unknown error")
	}

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)

	// Create error response
	errorResponse := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    appErr.Code,
			"message": appErr.Message,
			"type":    appErr.Type,
			"detail":  appErr.Detail,
		},
	}

	// Encode and send response
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
		// Fallback to plain text
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}