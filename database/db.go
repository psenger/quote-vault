package database

import (
	"database/sql"
	"quote-vault/models"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(connectionString string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Ping() error {
	return p.db.Ping()
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) GetAllQuotes(offset, limit int) ([]models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := p.db.Query(query, limit, offset)
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

func (p *PostgresDB) GetQuoteByID(id int) (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE id = $1`
	row := p.db.QueryRow(query, id)

	var quote models.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (p *PostgresDB) GetRandomQuote() (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`
	row := p.db.QueryRow(query)

	var quote models.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (p *PostgresDB) GetRandomQuoteByCategory(category string) (*models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE category = $1 ORDER BY RANDOM() LIMIT 1`
	row := p.db.QueryRow(query, category)

	var quote models.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func (p *PostgresDB) CreateQuote(quote *models.Quote) error {
	query := `INSERT INTO quotes (text, author, category) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := p.db.QueryRow(query, quote.Text, quote.Author, quote.Category).Scan(&quote.ID, &quote.CreatedAt)
	return err
}

func (p *PostgresDB) UpdateQuote(quote *models.Quote) error {
	query := `UPDATE quotes SET text = $1, author = $2, category = $3 WHERE id = $4`
	_, err := p.db.Exec(query, quote.Text, quote.Author, quote.Category, quote.ID)
	return err
}

func (p *PostgresDB) DeleteQuote(id int) error {
	query := `DELETE FROM quotes WHERE id = $1`
	_, err := p.db.Exec(query, id)
	return err
}

func (p *PostgresDB) GetTotalQuotes() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM quotes`
	err := p.db.QueryRow(query).Scan(&count)
	return count, err
}

func (p *PostgresDB) GetQuotesByCategory(category string, offset, limit int) ([]models.Quote, error) {
	query := `SELECT id, text, author, category, created_at FROM quotes WHERE category = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := p.db.Query(query, category, limit, offset)
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

func (p *PostgresDB) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM quotes ORDER BY category`
	rows, err := p.db.Query(query)
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