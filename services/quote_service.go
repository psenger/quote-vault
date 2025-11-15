package services

import (
	"database/sql"
	"fmt"

	"quote-vault/models"
)

type QuoteService struct {
	db *sql.DB
}

func NewQuoteService(db *sql.DB) *QuoteService {
	return &QuoteService{db: db}
}

func (s *QuoteService) CreateQuote(quote *models.Quote) error {
	query := `INSERT INTO quotes (text, author, category) VALUES (?, ?, ?) RETURNING id, created_at`
	err := s.db.QueryRow(query, quote.Text, quote.Author, quote.Category).Scan(&quote.ID, &quote.CreatedAt)
	return err
}

func (s *QuoteService) GetRandomQuote(category, author string) (*models.Quote, error) {
	var query string
	var args []interface{}

	baseQuery := "SELECT id, text, author, category, created_at FROM quotes"
	orderBy := " ORDER BY RANDOM() LIMIT 1"

	if category != "" && author != "" {
		query = baseQuery + " WHERE category = ? AND author = ?" + orderBy
		args = []interface{}{category, author}
	} else if category != "" {
		query = baseQuery + " WHERE category = ?" + orderBy
		args = []interface{}{category}
	} else if author != "" {
		query = baseQuery + " WHERE author = ?" + orderBy
		args = []interface{}{author}
	} else {
		query = baseQuery + orderBy
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
		return nil, err
	}

	return &quote, nil
}

func (s *QuoteService) GetQuotes(page, limit int, category, author string) ([]models.Quote, int, error) {
	offset := (page - 1) * limit

	// Build WHERE clause
	var whereClause string
	var args []interface{}
	var countArgs []interface{}

	if category != "" && author != "" {
		whereClause = " WHERE category = ? AND author = ?"
		args = []interface{}{category, author}
		countArgs = []interface{}{category, author}
	} else if category != "" {
		whereClause = " WHERE category = ?"
		args = []interface{}{category}
		countArgs = []interface{}{category}
	} else if author != "" {
		whereClause = " WHERE author = ?"
		args = []interface{}{author}
		countArgs = []interface{}{author}
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM quotes" + whereClause
	var total int
	err := s.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get quotes with pagination
	query := fmt.Sprintf("SELECT id, text, author, category, created_at FROM quotes%s ORDER BY created_at DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
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

	return quotes, total, nil
}