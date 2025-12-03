package repository

import (
	"database/sql"

	"quote-vault/errors"
	"quote-vault/models"
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
func (r *QuoteRepository) Create(quote *models.Quote) (*models.Quote, error) {
	query := `INSERT INTO quotes (text, author, category) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, quote.Text, quote.Author, quote.Category)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to create quote")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get last insert id")
	}
	quote.ID = int(id)

	return quote, nil
}

// GetByID retrieves a quote by its ID
func (r *QuoteRepository) GetByID(id int) (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE id = ?`

	quote := &models.Quote{}
	err := r.db.QueryRow(query, id).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.Category,
		&quote.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrQuoteNotFound
		}
		return nil, errors.NewDatabaseError("failed to get quote")
	}

	return quote, nil
}

// GetRandom retrieves a random quote
func (r *QuoteRepository) GetRandom() (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`

	quote := &models.Quote{}
	err := r.db.QueryRow(query).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.Category,
		&quote.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrQuoteNotFound
		}
		return nil, errors.NewDatabaseError("failed to get random quote")
	}

	return quote, nil
}

// GetRandomByCategory retrieves a random quote from a specific category
func (r *QuoteRepository) GetRandomByCategory(category string) (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE category = ? ORDER BY RANDOM() LIMIT 1`

	quote := &models.Quote{}
	err := r.db.QueryRow(query, category).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.Category,
		&quote.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrQuoteNotFound
		}
		return nil, errors.NewDatabaseError("failed to get random quote by category")
	}

	return quote, nil
}

// GetAll retrieves all quotes with pagination
func (r *QuoteRepository) GetAll(limit, offset int) ([]*models.Quote, int, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("failed to get quotes")
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
		)
		if err != nil {
			return nil, 0, errors.NewDatabaseError("failed to scan quote")
		}
		quotes = append(quotes, quote)
	}

	// Get total count
	var total int
	err = r.db.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&total)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("failed to get quote count")
	}

	return quotes, total, nil
}

// GetByCategory retrieves quotes by category with pagination
func (r *QuoteRepository) GetByCategory(category string, limit, offset int) ([]*models.Quote, int, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE category = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("failed to get quotes by category")
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
		)
		if err != nil {
			return nil, 0, errors.NewDatabaseError("failed to scan quote")
		}
		quotes = append(quotes, quote)
	}

	// Get total count for category
	var total int
	err = r.db.QueryRow("SELECT COUNT(*) FROM quotes WHERE category = ?", category).Scan(&total)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("failed to get quote count")
	}

	return quotes, total, nil
}

// GetCategories returns all unique categories
func (r *QuoteRepository) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM quotes WHERE category IS NOT NULL AND category != '' ORDER BY category`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get categories")
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, errors.NewDatabaseError("failed to scan category")
		}
		categories = append(categories, category)
	}

	return categories, nil
}
