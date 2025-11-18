package database

import "quote-vault/models"

// QuoteRepository defines the interface for quote data operations
type QuoteRepository interface {
	// CreateQuote adds a new quote to the database
	CreateQuote(quote *models.Quote) error
	
	// GetQuoteByID retrieves a quote by its ID
	GetQuoteByID(id int) (*models.Quote, error)
	
	// GetRandomQuote retrieves a random quote, optionally filtered by category
	GetRandomQuote(category string) (*models.Quote, error)
	
	// GetQuotes retrieves quotes with pagination
	GetQuotes(limit, offset int) ([]*models.Quote, error)
	
	// GetQuotesByCategory retrieves quotes by category with pagination
	GetQuotesByCategory(category string, limit, offset int) ([]*models.Quote, error)
	
	// GetTotalQuotesCount returns the total number of quotes
	GetTotalQuotesCount() (int, error)
	
	// GetTotalQuotesByCategoryCount returns the total number of quotes in a category
	GetTotalQuotesByCategoryCount(category string) (int, error)
	
	// GetCategories returns all unique categories
	GetCategories() ([]string, error)
	
	// Close closes the database connection
	Close() error
}