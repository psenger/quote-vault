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
	return &QuoteHandler{quoteService: quoteService}
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.quoteService.CreateQuote(&quote); err != nil {
		http.Error(w, "Failed to create quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quote)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	author := r.URL.Query().Get("author")

	quote, err := h.quoteService.GetRandomQuote(category, author)
	if err != nil {
		http.Error(w, "Failed to get random quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *QuoteHandler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10
	category := r.URL.Query().Get("category")
	author := r.URL.Query().Get("author")

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	quotes, total, err := h.quoteService.GetQuotes(page, limit, category, author)
	if err != nil {
		http.Error(w, "Failed to get quotes", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"quotes": quotes,
		"page":   page,
		"limit":  limit,
		"total":  total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *QuoteHandler) GetQuotesByCategory(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Category not specified", http.StatusBadRequest)
		return
	}
	category := pathParts[3]

	page := 1
	limit := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	quotes, total, err := h.quoteService.GetQuotes(page, limit, category, "")
	if err != nil {
		http.Error(w, "Failed to get quotes by category", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"quotes":   quotes,
		"category": category,
		"page":     page,
		"limit":    limit,
		"total":    total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *QuoteHandler) GetQuotesByAuthor(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Author not specified", http.StatusBadRequest)
		return
	}
	author := pathParts[3]

	page := 1
	limit := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	quotes, total, err := h.quoteService.GetQuotes(page, limit, "", author)
	if err != nil {
		http.Error(w, "Failed to get quotes by author", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"quotes": quotes,
		"author": author,
		"page":   page,
		"limit":  limit,
		"total":  total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}