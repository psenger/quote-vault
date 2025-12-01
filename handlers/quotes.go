package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quote-vault/errors"
	"quote-vault/services"
	"quote-vault/utils"
	"quote-vault/validators"
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
	var req validators.CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := validators.ValidateCreateQuote(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	quote, err := h.quoteService.CreateQuote(req.Text, req.Author, req.Category)
	if err != nil {
		if customErr, ok := err.(*errors.ValidationError); ok {
			utils.ErrorResponse(w, http.StatusBadRequest, customErr.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create quote")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, quote)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	quote, err := h.quoteService.GetRandomQuote(category)
	if err != nil {
		if customErr, ok := err.(*errors.NotFoundError); ok {
			utils.ErrorResponse(w, http.StatusNotFound, customErr.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get random quote")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, quote)
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

	category := r.URL.Query().Get("category")

	quotes, total, err := h.quoteService.ListQuotes(page, limit, category)
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