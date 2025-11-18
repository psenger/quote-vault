package router

import (
	"github.com/gorilla/mux"
	"github.com/quote-vault/handlers"
	"github.com/quote-vault/middleware"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CorsMiddleware)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// Quote routes
	api.HandleFunc("/quotes", handlers.GetQuotes).Methods("GET")
	api.HandleFunc("/quotes", middleware.ValidateQuoteMiddleware(handlers.CreateQuote)).Methods("POST")
	api.HandleFunc("/quotes/random", handlers.GetRandomQuote).Methods("GET")
	api.HandleFunc("/quotes/random/{category}", handlers.GetRandomQuoteByCategory).Methods("GET")

	return r
}