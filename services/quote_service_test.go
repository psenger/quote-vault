package services

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"quote-vault/models"
	"quote-vault/repository"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE quotes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL,
			author TEXT NOT NULL,
			category TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestQuoteService_CreateQuote(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewQuoteRepository(db)
	service := NewQuoteService(repo)

	tests := []struct {
		name    string
		quote   *models.Quote
		wantErr bool
	}{
		{
			name: "valid quote",
			quote: &models.Quote{
				Text:     "Test quote text",
				Author:   "Test Author",
				Category: "motivation",
			},
			wantErr: false,
		},
		{
			name: "empty text",
			quote: &models.Quote{
				Text:     "",
				Author:   "Test Author",
				Category: "motivation",
			},
			wantErr: true,
		},
		{
			name: "empty author",
			quote: &models.Quote{
				Text:     "Test quote text",
				Author:   "",
				Category: "motivation",
			},
			wantErr: true,
		},
		{
			name: "empty category defaults to general",
			quote: &models.Quote{
				Text:     "Test quote text",
				Author:   "Test Author",
				Category: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateQuote(tt.quote)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.ID == 0 {
					t.Error("CreateQuote() did not set ID")
				}
				if tt.quote.Category == "" && result.Category != "general" {
					t.Errorf("CreateQuote() category = %v, want 'general' for empty input", result.Category)
				}
			}
		})
	}
}

func TestQuoteService_GetQuoteByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewQuoteRepository(db)
	service := NewQuoteService(repo)

	// Create a test quote
	quote := &models.Quote{
		Text:     "Test quote",
		Author:   "Test Author",
		Category: "test",
	}
	created, err := service.CreateQuote(quote)
	if err != nil {
		t.Fatalf("failed to create test quote: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "valid ID",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "invalid ID (zero)",
			id:      0,
			wantErr: true,
		},
		{
			name:    "invalid ID (negative)",
			id:      -1,
			wantErr: true,
		},
		{
			name:    "non-existing ID",
			id:      9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetQuoteByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuoteByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.ID != tt.id {
				t.Errorf("GetQuoteByID() ID = %v, want %v", result.ID, tt.id)
			}
		})
	}
}

func TestQuoteService_GetRandomQuote(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewQuoteRepository(db)
	service := NewQuoteService(repo)

	// Test empty database
	_, err := service.GetRandomQuote("")
	if err == nil {
		t.Error("GetRandomQuote() should return error for empty database")
	}

	// Insert test quotes
	quotes := []*models.Quote{
		{Text: "Quote 1", Author: "Author 1", Category: "motivation"},
		{Text: "Quote 2", Author: "Author 2", Category: "humor"},
	}
	for _, q := range quotes {
		_, err := service.CreateQuote(q)
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}

	// Test getting random quote without category filter
	result, err := service.GetRandomQuote("")
	if err != nil {
		t.Errorf("GetRandomQuote() error = %v", err)
	}
	if result == nil {
		t.Error("GetRandomQuote() returned nil")
	}

	// Test getting random quote with category filter
	result, err = service.GetRandomQuote("motivation")
	if err != nil {
		t.Errorf("GetRandomQuote() with category error = %v", err)
	}
	if result != nil && result.Category != "motivation" {
		t.Errorf("GetRandomQuote() category = %v, want motivation", result.Category)
	}
}

func TestQuoteService_GetQuotes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewQuoteRepository(db)
	service := NewQuoteService(repo)

	// Insert test quotes
	for i := 0; i < 5; i++ {
		_, err := service.CreateQuote(&models.Quote{
			Text:     "Quote text",
			Author:   "Author",
			Category: "motivation",
		})
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}
	for i := 0; i < 3; i++ {
		_, err := service.CreateQuote(&models.Quote{
			Text:     "Quote text",
			Author:   "Author",
			Category: "humor",
		})
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}

	tests := []struct {
		name      string
		limit     int
		offset    int
		category  string
		wantCount int
		wantTotal int
	}{
		{
			name:      "all quotes",
			limit:     10,
			offset:    0,
			category:  "",
			wantCount: 8,
			wantTotal: 8,
		},
		{
			name:      "motivation category",
			limit:     10,
			offset:    0,
			category:  "motivation",
			wantCount: 5,
			wantTotal: 5,
		},
		{
			name:      "humor category",
			limit:     10,
			offset:    0,
			category:  "humor",
			wantCount: 3,
			wantTotal: 3,
		},
		{
			name:      "pagination",
			limit:     5,
			offset:    0,
			category:  "",
			wantCount: 5,
			wantTotal: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, total, err := service.GetQuotes(tt.limit, tt.offset, tt.category)
			if err != nil {
				t.Errorf("GetQuotes() error = %v", err)
				return
			}
			if len(quotes) != tt.wantCount {
				t.Errorf("GetQuotes() count = %v, want %v", len(quotes), tt.wantCount)
			}
			if total != tt.wantTotal {
				t.Errorf("GetQuotes() total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestQuoteService_GetCategories(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewQuoteRepository(db)
	service := NewQuoteService(repo)

	// Test empty database
	categories, err := service.GetCategories()
	if err != nil {
		t.Errorf("GetCategories() error = %v", err)
	}
	if len(categories) != 0 {
		t.Errorf("GetCategories() should return empty for empty db, got %v", len(categories))
	}

	// Insert quotes with different categories
	testData := []struct {
		text     string
		author   string
		category string
	}{
		{"Quote 1", "Author 1", "motivation"},
		{"Quote 2", "Author 2", "humor"},
		{"Quote 3", "Author 3", "wisdom"},
		{"Quote 4", "Author 4", "motivation"},
	}

	for _, td := range testData {
		_, err := service.CreateQuote(&models.Quote{
			Text:     td.text,
			Author:   td.author,
			Category: td.category,
		})
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}

	categories, err = service.GetCategories()
	if err != nil {
		t.Errorf("GetCategories() error = %v", err)
	}
	if len(categories) != 3 {
		t.Errorf("GetCategories() count = %v, want 3", len(categories))
	}
}
