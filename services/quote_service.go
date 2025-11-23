package services

import (
	"quote-vault/database"
	"quote-vault/errors"
	"quote-vault/models"
	"strings"
	"time"
)

type QuoteService struct {
	db database.Database
}

func NewQuoteService(db database.Database) *QuoteService {
	return &QuoteService{db: db}
}

func (s *QuoteService) CreateQuote(text, author, category string) (*models.Quote, error) {
	quote := &models.Quote{
		Text:      text,
		Author:    author,
		Category:  category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.db.CreateQuote(quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	quote, err := s.db.GetRandomQuote(category)
	if err != nil {
		return nil, err
	}

	if quote == nil {
		return nil, errors.ErrQuoteNotFound
	}

	return quote, nil
}

func (s *QuoteService) ListQuotes(page, limit int, category string) ([]*models.Quote, int, error) {
	offset := (page - 1) * limit

	quotes, err := s.db.ListQuotes(limit, offset, category)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.db.CountQuotes(category)
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (s *QuoteService) SearchQuotes(query string, page, limit int) ([]*models.Quote, int, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*models.Quote{}, 0, nil
	}

	offset := (page - 1) * limit

	quotes, err := s.db.SearchQuotes(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.db.CountSearchResults(query)
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}