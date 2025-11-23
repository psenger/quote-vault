package services

import (
	"database/sql"
	"fmt"
	"strings"

	"quote-vault/database"
	"quote-vault/models"
)

type QuoteService struct {
	db database.Database
}

func NewQuoteService(db database.Database) *QuoteService {
	return &QuoteService{
		db: db,
	}
}

func (s *QuoteService) CreateQuote(text, author, category string) (*models.Quote, error) {
	quote := &models.Quote{
		Text:     text,
		Author:   author,
		Category: category,
	}

	query := `INSERT INTO quotes (text, author, category) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := s.db.QueryRow(query, quote.Text, quote.Author, quote.Category).Scan(&quote.ID, &quote.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return quote, nil
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT id, text, author, category, created_at FROM quotes WHERE category = $1 ORDER BY RANDOM() LIMIT 1`
		args = append(args, category)
	} else {
		query = `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`
	}

	quote := &models.Quote{}
	err := s.db.QueryRow(query, args...).Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no quotes found")
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return quote, nil
}

func (s *QuoteService) GetAllQuotes(page, limit int) ([]*models.Quote, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM quotes`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get quote count: %w", err)
	}

	// Get quotes with pagination
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*models.Quote
	for rows.Next() {
		quote := &models.Quote{}
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, quote)
	}

	return quotes, total, nil
}

func (s *QuoteService) SearchQuotes(text, author, category string, page, limit int) ([]*models.Quote, int, error) {
	offset := (page - 1) * limit
	
	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if text != "" {
		conditions = append(conditions, fmt.Sprintf("text ILIKE $%d", argIndex))
		args = append(args, "%"+text+"%")
		argIndex++
	}

	if author != "" {
		conditions = append(conditions, fmt.Sprintf("author ILIKE $%d", argIndex))
		args = append(args, "%"+author+"%")
		argIndex++
	}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category ILIKE $%d", argIndex))
		args = append(args, "%"+category+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM quotes WHERE %s", whereClause)
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get search count: %w", err)
	}

	// Get quotes with search and pagination
	args = append(args, limit, offset)
	query := fmt.Sprintf("SELECT id, text, author, category, created_at FROM quotes WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", 
		whereClause, argIndex, argIndex+1)
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*models.Quote
	for rows.Next() {
		quote := &models.Quote{}
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, quote)
	}

	return quotes, total, nil
}