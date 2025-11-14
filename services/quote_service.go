package services

import (
	"database/sql"
	"errors"
	"math/rand"
	"quote-vault/database"
	"quote-vault/models"
	"time"
)

type QuoteService struct {
	db *database.DB
}

func NewQuoteService(db *database.DB) *QuoteService {
	rand.Seed(time.Now().UnixNano())
	return &QuoteService{db: db}
}

func (s *QuoteService) CreateQuote(text, author, category string) (*models.Quote, error) {
	if text == "" || author == "" || category == "" {
		return nil, errors.New("text, author, and category are required")
	}

	query := `INSERT INTO quotes (text, author, category, created_at) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, text, author, category, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Quote{
		ID:        int(id),
		Text:      text,
		Author:    author,
		Category:  category,
		CreatedAt: time.Now(),
	}, nil
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT id, text, author, category, created_at FROM quotes WHERE category = ? ORDER BY RANDOM() LIMIT 1`
		args = []interface{}{category}
	} else {
		query = `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`
		args = []interface{}{}
	}

	var quote models.Quote
	err := s.db.QueryRow(query, args...).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.Category,
		&quote.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &quote, nil
}

func (s *QuoteService) ListQuotes(page, limit int) ([]models.Quote, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get quotes with pagination
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		err := rows.Scan(
			&quote.ID,
			&quote.Text,
			&quote.Author,
			&quote.Category,
			&quote.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		quotes = append(quotes, quote)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}