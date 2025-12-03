package database

import (
	"database/sql"
	"quote-vault/models"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		author TEXT NOT NULL,
		category TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func (s *SQLiteDB) Ping() error {
	return s.db.Ping()
}

func (s *SQLiteDB) DB() *sql.DB {
	return s.db
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

func (s *SQLiteDB) GetAllQuotes(offset, limit int) ([]models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (s *SQLiteDB) GetQuoteByID(id int) (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var quote models.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (s *SQLiteDB) GetRandomQuote() (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`
	row := s.db.QueryRow(query)

	var quote models.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (s *SQLiteDB) GetRandomQuoteByCategory(category string) (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE category = ? ORDER BY RANDOM() LIMIT 1`
	row := s.db.QueryRow(query, category)

	var quote models.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (s *SQLiteDB) CreateQuote(quote *models.Quote) error {
	query := `INSERT INTO quotes (text, author, category) VALUES (?, ?, ?)`
	result, err := s.db.Exec(query, quote.Text, quote.Author, quote.Category)
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

func (s *SQLiteDB) UpdateQuote(quote *models.Quote) error {
	query := `UPDATE quotes SET text = ?, author = ?, category = ? WHERE id = ?`
	_, err := s.db.Exec(query, quote.Text, quote.Author, quote.Category, quote.ID)
	return err
}

func (s *SQLiteDB) DeleteQuote(id int) error {
	query := `DELETE FROM quotes WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SQLiteDB) GetTotalQuotes() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM quotes`
	err := s.db.QueryRow(query).Scan(&count)
	return count, err
}

func (s *SQLiteDB) GetQuotesByCategory(category string, offset, limit int) ([]models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE category = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (s *SQLiteDB) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM quotes ORDER BY category`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
