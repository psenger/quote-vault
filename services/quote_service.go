package services

import (
	"quote-vault/errors"
	"quote-vault/models"
	"quote-vault/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuoteService struct {
	repo repository.QuoteRepository
}

func NewQuoteService(repo repository.QuoteRepository) *QuoteService {
	return &QuoteService{
		repo: repo,
	}
}

func (s *QuoteService) CreateQuote(quote *models.Quote) error {
	return s.repo.Create(quote)
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	if category == "" {
		return s.repo.GetRandom()
	}
	return s.repo.GetRandomByCategory(category)
}

type PaginatedQuotes struct {
	Quotes     []*models.Quote `json:"quotes"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	HasNext    bool            `json:"has_next"`
	HasPrev    bool            `json:"has_prev"`
}

func (s *QuoteService) GetQuotes(c *gin.Context) (*PaginatedQuotes, error) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var quotes []*models.Quote
	var totalCount int64
	var err error

	if category != "" {
		quotes, err = s.repo.GetByCategory(category, limit, offset)
		if err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch quotes by category")
		}
		totalCount, err = s.repo.CountByCategory(category)
	} else {
		quotes, err = s.repo.GetAll(limit, offset)
		if err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch quotes")
		}
		totalCount, err = s.repo.Count()
	}

	if err != nil {
		return nil, errors.NewInternalServerError("Failed to count quotes")
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	return &PaginatedQuotes{
		Quotes:     quotes,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}, nil
}