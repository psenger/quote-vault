package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"quote-vault/models"
	"quote-vault/services"
)

type QuoteHandler struct {
	quoteService *services.QuoteService
}

func NewQuoteHandler(quoteService *services.QuoteService) *QuoteHandler {
	return &QuoteHandler{
		quoteService: quoteService,
	}
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if quote.Text == "" || quote.Author == "" {
		http.Error(w, "Text and author are required", http.StatusBadRequest)
		return
	}

	createdQuote, err := h.quoteService.CreateQuote(quote)
	if err != nil {
		http.Error(w, "Failed to create quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdQuote)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	
	// Handle empty category parameter
	if category != "" {
		category = strings.TrimSpace(category)
		if category == "" {
			category = ""
		}
	}

	quote, err := h.quoteService.GetRandomQuote(category)
	if err != nil {
		http.Error(w, "Failed to get random quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *QuoteHandler) ListQuotes(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	category := r.URL.Query().Get("category")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Handle empty category parameter
	if category != "" {
		category = strings.TrimSpace(category)
		if category == "" {
			category = ""
		}
	}

	quotes, total, err := h.quoteService.ListQuotes(page, limit, category)
	if err != nil {
		http.Error(w, "Failed to list quotes", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"quotes":      quotes,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + limit - 1) / limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}