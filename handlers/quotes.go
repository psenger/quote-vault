package handlers

import (
	"net/http"
	"quote-vault/errors"
	"quote-vault/models"
	"quote-vault/services"
	"quote-vault/utils"
	"quote-vault/validators"

	"github.com/gin-gonic/gin"
)

type QuoteHandler struct {
	service   *services.QuoteService
	validator *validators.QuoteValidator
}

func NewQuoteHandler(service *services.QuoteService, validator *validators.QuoteValidator) *QuoteHandler {
	return &QuoteHandler{
		service:   service,
		validator: validator,
	}
}

// CreateQuote handles POST /quotes
func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var quote models.Quote

	if err := c.ShouldBindJSON(&quote); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.ValidateQuote(&quote); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CreateQuote(&quote); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(c, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create quote")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, quote, "Quote created successfully")
}

// GetQuotes handles GET /quotes
func (h *QuoteHandler) GetQuotes(c *gin.Context) {
	result, err := h.service.GetQuotes(c)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(c, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch quotes")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "Quotes retrieved successfully")
}

// GetRandomQuote handles GET /quotes/random
func (h *QuoteHandler) GetRandomQuote(c *gin.Context) {
	category := c.Query("category")

	quote, err := h.service.GetRandomQuote(category)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(c, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch random quote")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, quote, "Random quote retrieved successfully")
}