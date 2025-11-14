package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"quote-vault/services"
)

type QuoteHandler struct {
	quoteService *services.QuoteService
}

type QuoteResponse struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

type CreateQuoteRequest struct {
	Text     string `json:"text" validate:"required,min=10"`
	Author   string `json:"author" validate:"required"`
	Category string `json:"category" validate:"required"`
}

type PaginatedQuotesResponse struct {
	Quotes []QuoteResponse `json:"quotes"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
	Total  int             `json:"total"`
}

func NewQuoteHandler(quoteService *services.QuoteService) *QuoteHandler {
	return &QuoteHandler{
		quoteService: quoteService,
	}
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	quote, err := h.quoteService.CreateQuote(req.Text, req.Author, req.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := QuoteResponse{
		ID:       quote.ID,
		Text:     quote.Text,
		Author:   quote.Author,
		Category: quote.Category,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	// Trim whitespace and handle empty string
	category = strings.TrimSpace(category)
	if category == "" {
		category = ""
	}

	quote, err := h.quoteService.GetRandomQuote(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if quote == nil {
		http.Error(w, "No quotes found", http.StatusNotFound)
		return
	}

	response := QuoteResponse{
		ID:       quote.ID,
		Text:     quote.Text,
		Author:   quote.Author,
		Category: quote.Category,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *QuoteHandler) ListQuotes(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	quotes, total, err := h.quoteService.ListQuotes(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseQuotes := make([]QuoteResponse, len(quotes))
	for i, quote := range quotes {
		responseQuotes[i] = QuoteResponse{
			ID:       quote.ID,
			Text:     quote.Text,
			Author:   quote.Author,
			Category: quote.Category,
		}
	}

	response := PaginatedQuotesResponse{
		Quotes: responseQuotes,
		Page:   page,
		Limit:  limit,
		Total:  total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}