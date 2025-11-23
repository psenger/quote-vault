package services

import (
	"fmt"
	"quote-vault/database"
	"quote-vault/errors"
	"quote-vault/models"
)

type QuoteService struct {
	db database.Database
}

func NewQuoteService(db database.Database) *QuoteService {
	return &QuoteService{db: db}
}

func (s *QuoteService) CreateQuote(quote *models.Quote) error {
	if s.db == nil {
		return errors.NewDatabaseError("database connection is nil")
	}

	err := s.db.CreateQuote(quote)
	if err != nil {
		return errors.NewDatabaseError(fmt.Sprintf("failed to create quote: %v", err))
	}
	return nil
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	if s.db == nil {
		return nil, errors.NewDatabaseError("database connection is nil")
	}

	quote, err := s.db.GetRandomQuote(category)
	if err != nil {
		return nil, errors.NewDatabaseError(fmt.Sprintf("failed to get random quote: %v", err))
	}
	return quote, nil
}

func (s *QuoteService) GetAllQuotes(offset, limit int) ([]*models.Quote, error) {
	if s.db == nil {
		return nil, errors.NewDatabaseError("database connection is nil")
	}

	if offset < 0 {
		return nil, errors.NewValidationError("offset cannot be negative")
	}
	if limit <= 0 || limit > 100 {
		return nil, errors.NewValidationError("limit must be between 1 and 100")
	}

	quotes, err := s.db.GetAllQuotes(offset, limit)
	if err != nil {
		return nil, errors.NewDatabaseError(fmt.Sprintf("failed to get quotes: %v", err))
	}
	return quotes, nil
}