package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"quote-vault/errors"
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

func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var req struct {
		Text     string `json:"text" binding:"required"`
		Author   string `json:"author" binding:"required"`
		Category string `json:"category" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidRequest.Error())
		return
	}

	quote, err := h.quoteService.CreateQuote(req.Text, req.Author, req.Category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Quote created successfully", quote)
}

func (h *QuoteHandler) GetRandomQuote(c *gin.Context) {
	category := c.Query("category")

	quote, err := h.quoteService.GetRandomQuote(category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Random quote retrieved", quote)
}

func (h *QuoteHandler) GetAllQuotes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	quotes, total, err := h.quoteService.GetAllQuotes(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"quotes": quotes,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Quotes retrieved successfully", response)
}

func (h *QuoteHandler) SearchQuotes(c *gin.Context) {
	text := c.Query("text")
	author := c.Query("author")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	if text == "" && author == "" && category == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "At least one search parameter (text, author, or category) is required")
		return
	}

	quotes, total, err := h.quoteService.SearchQuotes(text, author, category, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"quotes": quotes,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
		"search_params": map[string]string{
			"text":     text,
			"author":   author,
			"category": category,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Search completed successfully", response)
}