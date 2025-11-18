package services

import (
	"database/sql"
	"fmt"
	"quote-vault/models"
	"quote-vault/validators"
)

type QuoteService struct {
	db        *sql.DB
	validator *validators.QuoteValidator
}

func NewQuoteService(db *sql.DB) *QuoteService {
	return &QuoteService{
		db:        db,
		validator: validators.NewQuoteValidator(),
	}
}

func (s *QuoteService) CreateQuote(text, author, category string) (*models.Quote, error) {
	if err := s.validator.ValidateQuote(text, author, category); err != nil {
		return nil, err
	}

	query := `INSERT INTO quotes (text, author, category) VALUES ($1, $2, $3) RETURNING id, created_at`
	var quote models.Quote
	quote.Text = text
	quote.Author = author
	quote.Category = category

	err := s.db.QueryRow(query, text, author, category).Scan(&quote.ID, &quote.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return &quote, nil
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT id, text, author, category, created_at FROM quotes WHERE category = $1 ORDER BY RANDOM() LIMIT 1`
		args = []interface{}{category}
	else {
		query = `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`
	}

	var quote models.Quote
	err := s.db.QueryRow(query, args...).Scan(
		&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no quotes found")
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return &quote, nil
}

func (s *QuoteService) GetQuotes(limit, offset int) ([]models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes: %w", err)
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (s *QuoteService) GetQuoteCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get quote count: %w", err)
	}
	return count, nil
}