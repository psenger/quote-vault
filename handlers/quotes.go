package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quote-vault/errors"
	"quote-vault/models"
	"quote-vault/services"
	"quote-vault/utils"
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
	var req models.QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	quote := &models.Quote{
		Text:     req.Text,
		Author:   req.Author,
		Category: req.Category,
	}

	result, err := h.quoteService.CreateQuote(quote)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create quote")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, result)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	quote, err := h.quoteService.GetRandomQuote(category)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get random quote")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, quote)
}

func (h *QuoteHandler) GetRandomQuoteByCategory(w http.ResponseWriter, r *http.Request) {
	// This handler is for the /quotes/random/{category} route
	// Category is extracted from the URL path by the router
	h.GetRandomQuote(w, r)
}

func (h *QuoteHandler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	category := r.URL.Query().Get("category")
	offset := (page - 1) * limit

	quotes, total, err := h.quoteService.GetQuotes(limit, offset, category)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to list quotes")
		return
	}

	response := map[string]interface{}{
		"quotes": quotes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *QuoteHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL - this depends on router setup
	// For gorilla/mux, we'd use mux.Vars(r)["id"]
	utils.ErrorResponse(w, http.StatusNotImplemented, "Not implemented")
}

func (h *QuoteHandler) UpdateQuote(w http.ResponseWriter, r *http.Request) {
	utils.ErrorResponse(w, http.StatusNotImplemented, "Not implemented")
}

func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	utils.ErrorResponse(w, http.StatusNotImplemented, "Not implemented")
}

func (h *QuoteHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.quoteService.GetCategories()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get categories")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, categories)
}
