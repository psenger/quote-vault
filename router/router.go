package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"quote-vault/handlers"
	"quote-vault/middleware"
)

// NewRouter creates and configures the main router
func NewRouter(quoteHandler *handlers.QuoteHandler, healthHandler *handlers.HealthHandler) *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.ResponseTime)
	r.Use(middleware.CORS)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ErrorHandler)

	// Health check endpoints
	r.HandleFunc("/health", healthHandler.Health).Methods("GET")
	r.HandleFunc("/health/ready", healthHandler.Ready).Methods("GET")

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Quote routes
	api.HandleFunc("/quotes", quoteHandler.CreateQuote).Methods("POST")
	api.HandleFunc("/quotes", quoteHandler.GetQuotes).Methods("GET")
	api.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	api.HandleFunc("/quotes/random/{category}", quoteHandler.GetRandomQuoteByCategory).Methods("GET")
	api.HandleFunc("/quotes/{id:[0-9]+}", quoteHandler.GetQuote).Methods("GET")
	api.HandleFunc("/quotes/{id:[0-9]+}", quoteHandler.UpdateQuote).Methods("PUT")
	api.HandleFunc("/quotes/{id:[0-9]+}", quoteHandler.DeleteQuote).Methods("DELETE")

	// Category routes
	api.HandleFunc("/categories", quoteHandler.GetCategories).Methods("GET")

	// Set content type for API routes
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	return r
}