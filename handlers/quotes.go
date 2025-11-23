package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"quote-vault/errors"
	"quote-vault/services"
	"quote-vault/utils"
	"quote-vault/validators"
)

type QuoteHandler struct {
	service *services.QuoteService
}

func NewQuoteHandler(service *services.QuoteService) *QuoteHandler {
	return &QuoteHandler{service: service}
}

func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var req validators.CreateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validators.ValidateCreateQuote(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	quote, err := h.service.CreateQuote(req.Text, req.Author, req.Category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create quote", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Quote created successfully", quote)
}

func (h *QuoteHandler) GetRandomQuote(c *gin.Context) {
	category := c.Query("category")

	quote, err := h.service.GetRandomQuote(category)
	if err != nil {
		if err == errors.ErrQuoteNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "No quotes found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get random quote", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Random quote retrieved successfully", quote)
}

func (h *QuoteHandler) ListQuotes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	quotes, total, err := h.service.ListQuotes(page, limit, category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list quotes", err)
		return
	}

	response := map[string]interface{}{
		"quotes": quotes,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Quotes retrieved successfully", response)
}

func (h *QuoteHandler) SearchQuotes(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Search query is required", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	quotes, total, err := h.service.SearchQuotes(query, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search quotes", err)
		return
	}

	response := map[string]interface{}{
		"quotes": quotes,
		"query":  query,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Search results retrieved successfully", response)
}