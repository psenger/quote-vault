package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/quote-vault/models"
)

type DB struct {
	conn *sql.DB
}

var database *DB

func Init() error {
	// For now, use SQLite for simplicity
	connStr := "file:quotes.db?cache=shared&mode=rwc"
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	database = &DB{conn: db}
	
	// Create quotes table
	if err := database.createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
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
	
	_, err := db.conn.Exec(query)
	return err
}

func GetDB() *DB {
	return database
}

func (db *DB) AddQuote(quote *models.Quote) error {
	query := `INSERT INTO quotes (text, author, category) VALUES (?, ?, ?)`
	result, err := db.conn.Exec(query, quote.Text, quote.Author, quote.Category)
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
	
	if category != "" {
		query = `SELECT id, text, author, category FROM quotes WHERE category = ? ORDER BY RANDOM() LIMIT 1`
		args = append(args, category)
	} else {
		query = `SELECT id, text, author, category FROM quotes ORDER BY RANDOM() LIMIT 1`
	}
	
	row := db.conn.QueryRow(query, args...)
	
	quote := &models.Quote{}
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category)
	if err != nil {
		return nil, err
	}
	
	return quote, nil
}

func (db *DB) GetAllQuotes(offset, limit int) ([]*models.Quote, error) {
	query := `SELECT id, text, author, category FROM quotes LIMIT ? OFFSET ?`
	rows, err := db.conn.Query(query, limit, offset)
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