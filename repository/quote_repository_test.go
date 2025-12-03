package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"quote-vault/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create quotes table
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

func TestQuoteRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuoteRepository(db)

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
				Category: "test",
			},
			wantErr: false,
		},
		{
			name: "another valid quote",
			quote: &models.Quote{
				Text:     "Another test quote",
				Author:   "Another Author",
				Category: "motivation",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(tt.quote)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.ID == 0 {
					t.Error("Create() did not set ID")
				}
				if result.Text != tt.quote.Text {
					t.Errorf("Create() text = %v, want %v", result.Text, tt.quote.Text)
				}
			}
		})
	}
}

func TestQuoteRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuoteRepository(db)

	// Insert a test quote
	quote := &models.Quote{
		Text:     "Test quote for GetByID",
		Author:   "Test Author",
		Category: "test",
	}
	created, err := repo.Create(quote)
	if err != nil {
		t.Fatalf("failed to create test quote: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing quote",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existing quote",
			id:      9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.ID != tt.id {
				t.Errorf("GetByID() ID = %v, want %v", result.ID, tt.id)
			}
		})
	}
}

func TestQuoteRepository_GetRandom(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuoteRepository(db)

	// Test with empty database
	_, err := repo.GetRandom()
	if err == nil {
		t.Error("GetRandom() should return error for empty database")
	}

	// Insert test quotes
	quotes := []*models.Quote{
		{Text: "Quote 1", Author: "Author 1", Category: "cat1"},
		{Text: "Quote 2", Author: "Author 2", Category: "cat2"},
		{Text: "Quote 3", Author: "Author 3", Category: "cat1"},
	}
	for _, q := range quotes {
		_, err := repo.Create(q)
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}

	// Test GetRandom returns a quote
	result, err := repo.GetRandom()
	if err != nil {
		t.Errorf("GetRandom() error = %v", err)
	}
	if result == nil {
		t.Error("GetRandom() returned nil")
	}
}

func TestQuoteRepository_GetRandomByCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuoteRepository(db)

	// Insert test quotes
	quotes := []*models.Quote{
		{Text: "Motivation 1", Author: "Author 1", Category: "motivation"},
		{Text: "Motivation 2", Author: "Author 2", Category: "motivation"},
		{Text: "Humor 1", Author: "Author 3", Category: "humor"},
	}
	for _, q := range quotes {
		_, err := repo.Create(q)
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}

	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{
			name:     "existing category",
			category: "motivation",
			wantErr:  false,
		},
		{
			name:     "another existing category",
			category: "humor",
			wantErr:  false,
		},
		{
			name:     "non-existing category",
			category: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetRandomByCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandomByCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Category != tt.category {
				t.Errorf("GetRandomByCategory() category = %v, want %v", result.Category, tt.category)
			}
		})
	}
}

func TestQuoteRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuoteRepository(db)

	// Insert test quotes
	for i := 0; i < 15; i++ {
		_, err := repo.Create(&models.Quote{
			Text:     "Quote text",
			Author:   "Author",
			Category: "test",
		})
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
		wantTotal int
	}{
		{
			name:      "first page",
			limit:     10,
			offset:    0,
			wantCount: 10,
			wantTotal: 15,
		},
		{
			name:      "second page",
			limit:     10,
			offset:    10,
			wantCount: 5,
			wantTotal: 15,
		},
		{
			name:      "small limit",
			limit:     5,
			offset:    0,
			wantCount: 5,
			wantTotal: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, total, err := repo.GetAll(tt.limit, tt.offset)
			if err != nil {
				t.Errorf("GetAll() error = %v", err)
				return
			}
			if len(quotes) != tt.wantCount {
				t.Errorf("GetAll() count = %v, want %v", len(quotes), tt.wantCount)
			}
			if total != tt.wantTotal {
				t.Errorf("GetAll() total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestQuoteRepository_GetByCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuoteRepository(db)

	// Insert test quotes with different categories
	categories := []string{"motivation", "motivation", "motivation", "humor", "humor"}
	for i, cat := range categories {
		_, err := repo.Create(&models.Quote{
			Text:     "Quote text",
			Author:   "Author",
			Category: cat,
		})
		if err != nil {
			t.Fatalf("failed to create test quote %d: %v", i, err)
		}
	}

	tests := []struct {
		name      string
		category  string
		limit     int
		offset    int
		wantCount int
		wantTotal int
	}{
		{
			name:      "motivation category",
			category:  "motivation",
			limit:     10,
			offset:    0,
			wantCount: 3,
			wantTotal: 3,
		},
		{
			name:      "humor category",
			category:  "humor",
			limit:     10,
			offset:    0,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:      "non-existing category",
			category:  "nonexistent",
			limit:     10,
			offset:    0,
			wantCount: 0,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, total, err := repo.GetByCategory(tt.category, tt.limit, tt.offset)
			if err != nil {
				t.Errorf("GetByCategory() error = %v", err)
				return
			}
			if len(quotes) != tt.wantCount {
				t.Errorf("GetByCategory() count = %v, want %v", len(quotes), tt.wantCount)
			}
			if total != tt.wantTotal {
				t.Errorf("GetByCategory() total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestQuoteRepository_GetCategories(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewQuoteRepository(db)

	// Test empty database
	categories, err := repo.GetCategories()
	if err != nil {
		t.Errorf("GetCategories() error = %v", err)
	}
	if len(categories) != 0 {
		t.Errorf("GetCategories() should return empty slice for empty db, got %v", len(categories))
	}

	// Insert test quotes with different categories
	testCategories := []string{"motivation", "humor", "wisdom", "motivation"}
	for _, cat := range testCategories {
		_, err := repo.Create(&models.Quote{
			Text:     "Quote text",
			Author:   "Author",
			Category: cat,
		})
		if err != nil {
			t.Fatalf("failed to create test quote: %v", err)
		}
	}

	// Should return unique categories
	categories, err = repo.GetCategories()
	if err != nil {
		t.Errorf("GetCategories() error = %v", err)
	}
	if len(categories) != 3 {
		t.Errorf("GetCategories() count = %v, want 3", len(categories))
	}
}
