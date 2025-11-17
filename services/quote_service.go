package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"quote-vault/models"
)

type QuoteService struct {
	DB *sql.DB
}

type CreateQuoteRequest struct {
	Text     string `json:"text"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

type QuoteListResponse struct {
	Quotes []models.Quote `json:"quotes"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
	Total  int            `json:"total"`
}

type SearchParams struct {
	Query  string
	Author string
	Page   int
	Limit  int
}

type SearchResponse struct {
	Quotes []models.Quote `json:"quotes"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
	Total  int            `json:"total"`
	Query  string         `json:"query,omitempty"`
	Author string         `json:"author,omitempty"`
}

func NewQuoteService(db *sql.DB) *QuoteService {
	return &QuoteService{DB: db}
}

func (s *QuoteService) CreateQuote(req CreateQuoteRequest) (*models.Quote, error) {
	query := `INSERT INTO quotes (text, author, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	now := time.Now()

	result, err := s.DB.Exec(query, req.Text, req.Author, req.Category, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Quote{
		ID:        int(id),
		Text:      req.Text,
		Author:    req.Author,
		Category:  req.Category,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *QuoteService) GetRandomQuote(category string) (*models.Quote, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = "SELECT id, text, author, category, created_at, updated_at FROM quotes WHERE category = ? ORDER BY RANDOM() LIMIT 1"
		args = append(args, category)
	} else {
		query = "SELECT id, text, author, category, created_at, updated_at FROM quotes ORDER BY RANDOM() LIMIT 1"
	}

	var quote models.Quote
	err := s.DB.QueryRow(query, args...).Scan(
		&quote.ID,
		&quote.Text,
		&quote.Author,
		&quote.Category,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (s *QuoteService) ListQuotes(page, limit int) (*QuoteListResponse, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get quotes
	query := "SELECT id, text, author, category, created_at, updated_at FROM quotes ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
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
			&quote.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return &QuoteListResponse{
		Quotes: quotes,
		Page:   page,
		Limit:  limit,
		Total:  total,
	}, nil
}

func (s *QuoteService) SearchQuotes(params SearchParams) (*SearchResponse, error) {
	var whereClauses []string
	var args []interface{}

	if params.Query != "" {
		whereClauses = append(whereClauses, "text LIKE ?")
		args = append(args, "%"+params.Query+"%")
	}

	if params.Author != "" {
		whereClauses = append(whereClauses, "author LIKE ?")
		args = append(args, "%"+params.Author+"%")
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM quotes %s", whereClause)
	var total int
	err := s.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get quotes with pagination
	offset := (params.Page - 1) * params.Limit
	selectQuery := fmt.Sprintf("SELECT id, text, author, category, created_at, updated_at FROM quotes %s ORDER BY created_at DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, params.Limit, offset)

	rows, err := s.DB.Query(selectQuery, args...)
	if err != nil {
		return nil, err
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
			&quote.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return &SearchResponse{
		Quotes: quotes,
		Page:   params.Page,
		Limit:  params.Limit,
		Total:  total,
		Query:  params.Query,
		Author: params.Author,
	}, nil
}