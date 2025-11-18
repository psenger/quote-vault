package services

import (
	"errors"
	"quote-vault/database"
	"quote-vault/models"
)

var (
	ErrQuoteNotFound = errors.New("quote not found")
	ErrInvalidLimit  = errors.New("limit must be between 1 and 100")
	ErrInvalidOffset = errors.New("offset must be non-negative")
)

// QuoteService handles business logic for quotes
type QuoteService struct {
	repo database.QuoteRepository
}

// NewQuoteService creates a new quote service
func NewQuoteService(repo database.QuoteRepository) *QuoteService {
	return &QuoteService{
		repo: repo,
	}
}

// CreateQuote creates a new quote
func (s *QuoteService) CreateQuote(text, author, category string) (*models.Quote, error) {
	quote := &models.Quote{
		Text:     text,
		Author:   author,
		Category: category,
	}

	err := s.repo.CreateQuote(quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

// GetQuoteByID retrieves a quote by ID
func (s *QuoteService) GetQuoteByID(id int) (*models.Quote, error) {
	quote, err := s.repo.GetQuoteByID(id)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, ErrQuoteNotFound
	}
	return quote, nil
}

// GetRandomQuote retrieves a random quote
func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	quote, err := s.repo.GetRandomQuote(category)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, ErrQuoteNotFound
	}
	return quote, nil
}

// GetQuotes retrieves quotes with pagination
func (s *QuoteService) GetQuotes(limit, offset int, category string) (*models.PaginatedQuotes, error) {
	if err := s.validatePagination(limit, offset); err != nil {
		return nil, err
	}

	var quotes []*models.Quote
	var total int
	var err error

	if category != "" {
		quotes, err = s.repo.GetQuotesByCategory(category, limit, offset)
		if err != nil {
			return nil, err
		}
		total, err = s.repo.GetTotalQuotesByCategoryCount(category)
	} else {
		quotes, err = s.repo.GetQuotes(limit, offset)
		if err != nil {
			return nil, err
		}
		total, err = s.repo.GetTotalQuotesCount()
	}

	if err != nil {
		return nil, err
	}

	return &models.PaginatedQuotes{
		Quotes: quotes,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetCategories retrieves all unique categories
func (s *QuoteService) GetCategories() ([]string, error) {
	return s.repo.GetCategories()
}

func (s *QuoteService) validatePagination(limit, offset int) error {
	if limit < 1 || limit > 100 {
		return ErrInvalidLimit
	}
	if offset < 0 {
		return ErrInvalidOffset
	}
	return nil
}