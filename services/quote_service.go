package services

import (
	"errors"
	"math/rand"
	"quote-vault/models"
	"strings"
	"time"
)

type QuoteService struct {
	quotes []models.Quote
	nextID int
}

func NewQuoteService() *QuoteService {
	rand.Seed(time.Now().UnixNano())
	return &QuoteService{
		quotes: make([]models.Quote, 0),
		nextID: 1,
	}
}

func (s *QuoteService) AddQuote(text, author, category string) (*models.Quote, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("quote text cannot be empty")
	}
	if strings.TrimSpace(author) == "" {
		return nil, errors.New("author cannot be empty")
	}
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("category cannot be empty")
	}

	quote := models.Quote{
		ID:       s.nextID,
		Text:     strings.TrimSpace(text),
		Author:   strings.TrimSpace(author),
		Category: strings.TrimSpace(category),
	}

	s.quotes = append(s.quotes, quote)
	s.nextID++

	return &quote, nil
}

func (s *QuoteService) GetRandomQuote() (*models.Quote, error) {
	if len(s.quotes) == 0 {
		return nil, errors.New("no quotes available")
	}

	index := rand.Intn(len(s.quotes))
	return &s.quotes[index], nil
}

func (s *QuoteService) GetRandomQuoteByCategory(category string) (*models.Quote, error) {
	categoryQuotes := make([]models.Quote, 0)

	for _, quote := range s.quotes {
		if strings.EqualFold(quote.Category, category) {
			categoryQuotes = append(categoryQuotes, quote)
		}
	}

	if len(categoryQuotes) == 0 {
		return nil, errors.New("no quotes found for this category")
	}

	index := rand.Intn(len(categoryQuotes))
	return &categoryQuotes[index], nil
}

func (s *QuoteService) GetAllQuotes(offset, limit int) ([]models.Quote, error) {
	if offset < 0 || limit < 0 {
		return nil, errors.New("offset and limit must be non-negative")
	}

	totalQuotes := len(s.quotes)
	if offset >= totalQuotes {
		return []models.Quote{}, nil
	}

	end := offset + limit
	if end > totalQuotes {
		end = totalQuotes
	}

	return s.quotes[offset:end], nil
}

func (s *QuoteService) GetQuoteCount() int {
	return len(s.quotes)
}