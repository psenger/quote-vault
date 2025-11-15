package router

import (
	"net/http"

	"quote-vault/handlers"
	"quote-vault/middleware"
)

func SetupRoutes(quoteHandler *handlers.QuoteHandler) http.Handler {
	mux := http.NewServeMux()

	// Quote routes
	mux.HandleFunc("/api/quotes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			quoteHandler.CreateQuote(w, r)
		case http.MethodGet:
			quoteHandler.GetQuotes(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/quotes/random", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		quoteHandler.GetRandomQuote(w, r)
	})

	// Filter routes
	mux.HandleFunc("/api/quotes/category/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		quoteHandler.GetQuotesByCategory(w, r)
	})

	mux.HandleFunc("/api/quotes/author/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		quoteHandler.GetQuotesByAuthor(w, r)
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware
	handler := middleware.CORS(mux)
	handler = middleware.Logging(handler)
	handler = middleware.ValidationMiddleware(handler)

	return handler
}