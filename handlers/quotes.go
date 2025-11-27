package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quote-vault/errors"
	"quote-vault/models"
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
	var quote models.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if err := validators.ValidateQuote(&quote); err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdQuote, err := h.quoteService.CreateQuote(&quote)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, createdQuote, http.StatusCreated)
}

func (h *QuoteHandler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.NewPaginationParams(r)
	category := r.URL.Query().Get("category")

	quotes, total, err := h.quoteService.GetQuotes(paginationParams.Limit, paginationParams.Offset, category)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	paginationMeta := paginationParams.CalculateMeta(total)
	response := utils.NewPaginatedResponse(quotes, paginationMeta)

	utils.SuccessResponse(w, response, http.StatusOK)
}

func (h *QuoteHandler) GetQuoteByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/quotes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID, http.StatusBadRequest)
		return
	}

	quote, err := h.quoteService.GetQuoteByID(id)
	if err != nil {
		if err == errors.ErrQuoteNotFound {
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		}
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, quote, http.StatusOK)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	quote, err := h.quoteService.GetRandomQuote(category)
	if err != nil {
		if err == errors.ErrQuoteNotFound {
			utils.ErrorResponse(w, err, http.StatusNotFound)
			return
		}
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, quote, http.StatusOK)
}