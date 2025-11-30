package database

import (
	"context"

	"github.com/quote-vault/models"
)

// QuoteRepositoryInterface defines the contract for quote repository operations
type QuoteRepositoryInterface interface {
	// Create adds a new quote to the database
	Create(ctx context.Context, quote *models.Quote) error

	// GetByID retrieves a quote by its ID
	GetByID(ctx context.Context, id int64) (*models.Quote, error)

	// GetRandom retrieves a random quote, optionally filtered by category
	GetRandom(ctx context.Context, category string) (*models.Quote, error)

	// GetAll retrieves all quotes with pagination and optional category filter
	GetAll(ctx context.Context, limit, offset int, category string) ([]*models.Quote, error)

	// GetCount returns the total count of quotes, optionally filtered by category
	GetCount(ctx context.Context, category string) (int64, error)

	// GetCategories returns all unique categories
	GetCategories(ctx context.Context) ([]string, error)

	// Search searches for quotes by text, author, or category
	Search(ctx context.Context, searchTerm string, limit, offset int) ([]*models.Quote, error)
}

// DatabaseInterface defines the contract for database operations
type DatabaseInterface interface {
	// Connect establishes a connection to the database
	Connect() error

	// Close closes the database connection
	Close() error

	// Ping checks if the database connection is alive
	Ping() error

	// Migrate runs database migrations
	Migrate() error

	// GetDB returns the underlying database connection
	GetDB() interface{}
}

// TransactionInterface defines the contract for database transactions
type TransactionInterface interface {
	// Begin starts a new database transaction
	Begin(ctx context.Context) (interface{}, error)

	// Commit commits the current transaction
	Commit(tx interface{}) error

	// Rollback rolls back the current transaction
	Rollback(tx interface{}) error
}