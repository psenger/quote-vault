package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/quote-vault/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database connection established")
	return database, nil
}

func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		author TEXT NOT NULL,
		category TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	return err
}

func (db *DB) CreateQuote(quote *models.Quote) error {
	query := `INSERT INTO quotes (text, author, category) VALUES (?, ?, ?)`
	result, err := db.Exec(query, quote.Text, quote.Author, quote.Category)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	quote.ID = int(id)
	return nil
}

func (db *DB) GetRandomQuote(category string) (*models.Quote, error) {
	var query string
	var args []interface{}

	if category == "" {
		query = `SELECT id, text, author, category FROM quotes ORDER BY RANDOM() LIMIT 1`
	} else {
		query = `SELECT id, text, author, category FROM quotes WHERE category = ? ORDER BY RANDOM() LIMIT 1`
		args = append(args, category)
	}

	row := db.QueryRow(query, args...)

	quote := &models.Quote{}
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (db *DB) GetQuotes(limit, offset int) ([]*models.Quote, error) {
	query := `SELECT id, text, author, category FROM quotes ORDER BY id LIMIT ? OFFSET ?`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []*models.Quote
	for rows.Next() {
		quote := &models.Quote{}
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}