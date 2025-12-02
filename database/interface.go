package database

import "quote-vault/models"

type Database interface {
	GetAllQuotes(offset, limit int) ([]models.Quote, error)
	GetQuoteByID(id int) (*models.Quote, error)
	GetRandomQuote() (*models.Quote, error)
	GetRandomQuoteByCategory(category string) (*models.Quote, error)
	CreateQuote(quote *models.Quote) error
	UpdateQuote(quote *models.Quote) error
	DeleteQuote(id int) error
	GetTotalQuotes() (int, error)
	GetQuotesByCategory(category string, offset, limit int) ([]models.Quote, error)
	GetCategories() ([]string, error)
	Ping() error
	Close() error
}