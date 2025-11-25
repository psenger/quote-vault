package repository

import (
	"fmt"
	"quote-vault/models"

	"gorm.io/gorm"
)

// QuoteRepository interface defines the contract for quote data operations
type QuoteRepository interface {
	Create(quote *models.Quote) error
	GetAll(limit, offset int) ([]*models.Quote, error)
	GetByCategory(category string, limit, offset int) ([]*models.Quote, error)
	GetRandom() (*models.Quote, error)
	GetRandomByCategory(category string) (*models.Quote, error)
	Count() (int64, error)
	CountByCategory(category string) (int64, error)
}

// quoteRepository implements QuoteRepository interface
type quoteRepository struct {
	db *gorm.DB
}

// NewQuoteRepository creates a new quote repository instance
func NewQuoteRepository(db *gorm.DB) QuoteRepository {
	return &quoteRepository{
		db: db,
	}
}

// Create adds a new quote to the database
func (r *quoteRepository) Create(quote *models.Quote) error {
	result := r.db.Create(quote)
	if result.Error != nil {
		return fmt.Errorf("failed to create quote: %w", result.Error)
	}
	return nil
}

// GetAll retrieves all quotes with pagination
func (r *quoteRepository) GetAll(limit, offset int) ([]*models.Quote, error) {
	var quotes []*models.Quote
	result := r.db.Limit(limit).Offset(offset).Find(&quotes)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch quotes: %w", result.Error)
	}
	return quotes, nil
}

// GetByCategory retrieves quotes by category with pagination
func (r *quoteRepository) GetByCategory(category string, limit, offset int) ([]*models.Quote, error) {
	var quotes []*models.Quote
	result := r.db.Where("category = ?", category).Limit(limit).Offset(offset).Find(&quotes)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch quotes by category: %w", result.Error)
	}
	return quotes, nil
}

// GetRandom retrieves a random quote from all quotes
func (r *quoteRepository) GetRandom() (*models.Quote, error) {
	var quote models.Quote
	result := r.db.Order("RANDOM()").First(&quote)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch random quote: %w", result.Error)
	}
	return &quote, nil
}

// GetRandomByCategory retrieves a random quote from a specific category
func (r *quoteRepository) GetRandomByCategory(category string) (*models.Quote, error) {
	var quote models.Quote
	result := r.db.Where("category = ?", category).Order("RANDOM()").First(&quote)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch random quote by category: %w", result.Error)
	}
	return &quote, nil
}

// Count returns the total number of quotes
func (r *quoteRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&models.Quote{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count quotes: %w", result.Error)
	}
	return count, nil
}

// CountByCategory returns the total number of quotes in a specific category
func (r *quoteRepository) CountByCategory(category string) (int64, error) {
	var count int64
	result := r.db.Model(&models.Quote{}).Where("category = ?", category).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count quotes by category: %w", result.Error)
	}
	return count, nil
}