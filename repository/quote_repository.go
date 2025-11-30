package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/quote-vault/errors"
	"github.com/quote-vault/models"
)

// QuoteRepository handles database operations for quotes
type QuoteRepository struct {
	db *sql.DB
}

// NewQuoteRepository creates a new quote repository
func NewQuoteRepository(db *sql.DB) *QuoteRepository {
	return &QuoteRepository{
		db: db,
	}
}

// Create adds a new quote to the database
func (r *QuoteRepository) Create(ctx context.Context, quote *models.Quote) error {
	query := `
		INSERT INTO quotes (text, author, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	quote.CreatedAt = now
	quote.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query, quote.Text, quote.Author, quote.Category, quote.CreatedAt, quote.UpdatedAt).Scan(&quote.ID)
	if err != nil {
		return errors.NewDatabaseError("failed to create quote", err)
	}

	return nil
}

// GetByID retrieves a quote by its ID
func (r *QuoteRepository) GetByID(ctx context.Context, id int64) (*models.Quote, error) {
	query := `
		SELECT id, text, author, category, created_at, updated_at
		FROM quotes
		WHERE id = $1`

	quote := &models.Quote{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.Category,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("quote not found")
		}
		return nil, errors.NewDatabaseError("failed to get quote", err)
	}

	return quote, nil
}

// GetRandom retrieves a random quote, optionally filtered by category
func (r *QuoteRepository) GetRandom(ctx context.Context, category string) (*models.Quote, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `
			SELECT id, text, author, category, created_at, updated_at
			FROM quotes
			WHERE category = $1
			ORDER BY RANDOM()
			LIMIT 1`
		args = append(args, category)
	} else {
		query = `
			SELECT id, text, author, category, created_at, updated_at
			FROM quotes
			ORDER BY RANDOM()
			LIMIT 1`
	}

	quote := &models.Quote{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.Category,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("no quotes found")
		}
		return nil, errors.NewDatabaseError("failed to get random quote", err)
	}

	return quote, nil
}

// GetAll retrieves all quotes with pagination and optional category filter
func (r *QuoteRepository) GetAll(ctx context.Context, limit, offset int, category string) ([]*models.Quote, error) {
	var query string
	var args []interface{}
	var whereClause string

	// Build where clause if category is specified
	if category != "" {
		whereClause = "WHERE category = $1"
		args = append(args, category)
	}

	// Build the complete query
	query = fmt.Sprintf(`
		SELECT id, text, author, category, created_at, updated_at
		FROM quotes
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause,
		len(args)+1,
		len(args)+2,
	)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get quotes", err)
	}
	defer rows.Close()

	var quotes []*models.Quote
	for rows.Next() {
		quote := &models.Quote{}
		err := rows.Scan(
			&quote.ID,
			&quote.Text,
			&quote.Author,
			&quote.Category,
			&quote.CreatedAt,
			&quote.UpdatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan quote", err)
		}
		quotes = append(quotes, quote)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewDatabaseError("error iterating rows", err)
	}

	return quotes, nil
}

// GetCount returns the total count of quotes, optionally filtered by category
func (r *QuoteRepository) GetCount(ctx context.Context, category string) (int64, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = "SELECT COUNT(*) FROM quotes WHERE category = $1"
		args = append(args, category)
	} else {
		query = "SELECT COUNT(*) FROM quotes"
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, errors.NewDatabaseError("failed to get quote count", err)
	}

	return count, nil
}

// GetCategories returns all unique categories
func (r *QuoteRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT category
		FROM quotes
		WHERE category IS NOT NULL AND category != ''
		ORDER BY category`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get categories", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, errors.NewDatabaseError("failed to scan category", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewDatabaseError("error iterating category rows", err)
	}

	return categories, nil
}

// Search searches for quotes by text, author, or category
func (r *QuoteRepository) Search(ctx context.Context, searchTerm string, limit, offset int) ([]*models.Quote, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return r.GetAll(ctx, limit, offset, "")
	}

	searchPattern := "%" + strings.ToLower(searchTerm) + "%"
	query := `
		SELECT id, text, author, category, created_at, updated_at
		FROM quotes
		WHERE LOWER(text) LIKE $1
		   OR LOWER(author) LIKE $1
		   OR LOWER(category) LIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, searchPattern, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to search quotes", err)
	}
	defer rows.Close()

	var quotes []*models.Quote
	for rows.Next() {
		quote := &models.Quote{}
		err := rows.Scan(
			&quote.ID,
			&quote.Text,
			&quote.Author,
			&quote.Category,
			&quote.CreatedAt,
			&quote.UpdatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan search result", err)
		}
		quotes = append(quotes, quote)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewDatabaseError("error iterating search results", err)
	}

	return quotes, nil
}