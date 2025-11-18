package database

import (
	"database/sql"
	"fmt"
	"quote-vault/models"

	_ "github.com/lib/pq"
)

// PostgreSQLRepository implements the QuoteRepository interface
type PostgreSQLRepository struct {
	db *sql.DB
}

// NewPostgreSQLRepository creates a new PostgreSQL repository
func NewPostgreSQLRepository(dataSourceName string) (*PostgreSQLRepository, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create quotes table if it doesn't exist
	if err := createQuotesTable(db); err != nil {
		return nil, fmt.Errorf("failed to create quotes table: %w", err)
	}

	return &PostgreSQLRepository{db: db}, nil
}

// CreateQuote adds a new quote to the database
func (r *PostgreSQLRepository) CreateQuote(quote *models.Quote) error {
	query := `INSERT INTO quotes (text, author, category) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(query, quote.Text, quote.Author, quote.Category).Scan(&quote.ID, &quote.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create quote: %w", err)
	}
	return nil
}

// GetQuoteByID retrieves a quote by its ID
func (r *PostgreSQLRepository) GetQuoteByID(id int) (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE id = $1`
	quote := &models.Quote{}
	err := r.db.QueryRow(query, id).Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get quote by ID: %w", err)
	}
	return quote, nil
}

// GetRandomQuote retrieves a random quote, optionally filtered by category
func (r *PostgreSQLRepository) GetRandomQuote(category string) (*models.Quote, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT id, text, author, category, created_at FROM quotes WHERE category = $1 ORDER BY RANDOM() LIMIT 1`
		args = append(args, category)
	} else {
		query = `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`
	}

	quote := &models.Quote{}
	err := r.db.QueryRow(query, args...).Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}
	return quote, nil
}

// GetQuotes retrieves quotes with pagination
func (r *PostgreSQLRepository) GetQuotes(limit, offset int) ([]*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	return r.executeQuoteQuery(query, limit, offset)
}

// GetQuotesByCategory retrieves quotes by category with pagination
func (r *PostgreSQLRepository) GetQuotesByCategory(category string, limit, offset int) ([]*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE category = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	return r.executeQuoteQueryWithCategory(query, category, limit, offset)
}

// GetTotalQuotesCount returns the total number of quotes
func (r *PostgreSQLRepository) GetTotalQuotesCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM quotes`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total quotes count: %w", err)
	}
	return count, nil
}

// GetTotalQuotesByCategoryCount returns the total number of quotes in a category
func (r *PostgreSQLRepository) GetTotalQuotesByCategoryCount(category string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM quotes WHERE category = $1`
	err := r.db.QueryRow(query, category).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total quotes count by category: %w", err)
	}
	return count, nil
}

// GetCategories returns all unique categories
func (r *PostgreSQLRepository) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM quotes ORDER BY category`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

// Close closes the database connection
func (r *PostgreSQLRepository) Close() error {
	return r.db.Close()
}

// Helper functions

func createQuotesTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS quotes (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL,
			author VARCHAR(255) NOT NULL,
			category VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func (r *PostgreSQLRepository) executeQuoteQuery(query string, args ...interface{}) ([]*models.Quote, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute quote query: %w", err)
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *PostgreSQLRepository) executeQuoteQueryWithCategory(query string, category string, limit, offset int) ([]*models.Quote, error) {
	rows, err := r.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute quote query with category: %w", err)
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *PostgreSQLRepository) scanQuotes(rows *sql.Rows) ([]*models.Quote, error) {
	var quotes []*models.Quote
	for rows.Next() {
		quote := &models.Quote{}
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, quote)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating quotes: %w", err)
	}

	return quotes, nil
}