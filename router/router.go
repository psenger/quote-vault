package router

import (
	"net/http"

	"quote-vault/handlers"
	"quote-vault/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(quoteHandler *handlers.QuoteHandler) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Quote routes
	api.HandleFunc("/quotes", middleware.ValidationMiddleware(quoteHandler.CreateQuote)).Methods("POST")
	api.HandleFunc("/quotes", quoteHandler.ListQuotes).Methods("GET")
	api.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	api.HandleFunc("/quotes/search", quoteHandler.SearchQuotes).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return r
}