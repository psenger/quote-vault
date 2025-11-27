package services

import (
	"quote-vault/errors"
	"quote-vault/models"
	"quote-vault/repository"
)

type QuoteService struct {
	quoteRepo *repository.QuoteRepository
}

func NewQuoteService(quoteRepo *repository.QuoteRepository) *QuoteService {
	return &QuoteService{
		quoteRepo: quoteRepo,
	}
}

func (s *QuoteService) CreateQuote(quote *models.Quote) (*models.Quote, error) {
	if quote.Text == "" {
		return nil, errors.ErrEmptyQuoteText
	}

	if quote.Author == "" {
		return nil, errors.ErrEmptyAuthor
	}

	if quote.Category == "" {
		quote.Category = "general"
	}

	return s.quoteRepo.Create(quote)
}

func (s *QuoteService) GetQuotes(limit, offset int, category string) ([]*models.Quote, int, error) {
	if category == "" {
		return s.quoteRepo.GetAll(limit, offset)
	}
	return s.quoteRepo.GetByCategory(category, limit, offset)
}

func (s *QuoteService) GetQuoteByID(id int) (*models.Quote, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidID
	}

	return s.quoteRepo.GetByID(id)
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	if category == "" {
		return s.quoteRepo.GetRandom()
	}
	return s.quoteRepo.GetRandomByCategory(category)
}

func (s *QuoteService) GetCategories() ([]string, error) {
	return s.quoteRepo.GetCategories()
}