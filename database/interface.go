package database

import "quote-vault/models"

type Database interface {
	Connect() error
	Close() error
	CreateQuote(quote *models.Quote) error
	GetRandomQuote(category string) (*models.Quote, error)
	ListQuotes(limit, offset int, category string) ([]*models.Quote, error)
	CountQuotes(category string) (int, error)
	SearchQuotes(query string, limit, offset int) ([]*models.Quote, error)
	CountSearchResults(query string) (int, error)
}