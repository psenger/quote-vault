package router

import (
	"net/http"
	"quote-vault/handlers"
	"quote-vault/middleware"
)

// SetupRoutes configures and returns the HTTP router
func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Apply CORS middleware to all routes
	mux.HandleFunc("/quotes", middleware.CORS(middleware.Logging(middleware.ValidateQuoteInput(handlers.HandleQuotes))))
	mux.HandleFunc("/quotes/random", middleware.CORS(middleware.Logging(handlers.GetRandomQuote)))

	return mux
}