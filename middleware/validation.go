package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResponse represents the validation error response
type ValidationResponse struct {
	Error  string            `json:"error"`
	Errors []ValidationError `json:"errors"`
}

// ValidateQuoteInput validates quote input data
func ValidateQuoteInput(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		var quote struct {
			Text     string `json:"text"`
			Author   string `json:"author"`
			Category string `json:"category"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&quote); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ValidationResponse{
				Error: "Invalid JSON format",
			})
			return
		}

		var errors []ValidationError

		// Validate quote text
		if strings.TrimSpace(quote.Text) == "" {
			errors = append(errors, ValidationError{
				Field:   "text",
				Message: "Quote text is required",
			})
		} else if len(quote.Text) > 1000 {
			errors = append(errors, ValidationError{
				Field:   "text",
				Message: "Quote text must be less than 1000 characters",
			})
		}

		// Validate author
		if strings.TrimSpace(quote.Author) == "" {
			errors = append(errors, ValidationError{
				Field:   "author",
				Message: "Author is required",
			})
		} else if len(quote.Author) > 100 {
			errors = append(errors, ValidationError{
				Field:   "author",
				Message: "Author name must be less than 100 characters",
			})
		}

		// Validate category
		if strings.TrimSpace(quote.Category) == "" {
			errors = append(errors, ValidationError{
				Field:   "category",
				Message: "Category is required",
			})
		} else if len(quote.Category) > 50 {
			errors = append(errors, ValidationError{
				Field:   "category",
				Message: "Category must be less than 50 characters",
			})
		}

		if len(errors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ValidationResponse{
				Error:  "Validation failed",
				Errors: errors,
			})
			return
		}

		next.ServeHTTP(w, r)
	}
}