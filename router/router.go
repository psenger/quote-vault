package router

import (
	"github.com/gorilla/mux"
	"quote-vault/handlers"
)

// SetupRoutes configures all API routes
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	
	// API prefix
	api := r.PathPrefix("/api/v1").Subrouter()
	
	// Quote routes
	api.HandleFunc("/quotes", handlers.GetQuotes).Methods("GET")
	api.HandleFunc("/quotes", handlers.CreateQuote).Methods("POST")
	api.HandleFunc("/quotes/random", handlers.GetRandomQuote).Methods("GET")
	api.HandleFunc("/quotes/random/{category}", handlers.GetRandomQuoteByCategory).Methods("GET")
	
	return r
}