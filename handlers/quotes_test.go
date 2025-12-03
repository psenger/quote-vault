package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"quote-vault/repository"
	"quote-vault/services"
)

func setupTestHandler(t *testing.T) (*QuoteHandler, *sql.DB) {
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

	repo := repository.NewQuoteRepository(db)
	service := services.NewQuoteService(repo)
	handler := NewQuoteHandler(service)

	return handler, db
}

func TestQuoteHandler_CreateQuote(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	tests := []struct {
		name           string
		body           map[string]string
		wantStatusCode int
	}{
		{
			name: "valid quote",
			body: map[string]string{
				"text":     "Test quote text",
				"author":   "Test Author",
				"category": "motivation",
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "missing text",
			body: map[string]string{
				"author":   "Test Author",
				"category": "motivation",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing author",
			body: map[string]string{
				"text":     "Test quote text",
				"category": "motivation",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			body:           map[string]string{},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/quotes", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateQuote(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("CreateQuote() status = %v, want %v", rec.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestQuoteHandler_CreateQuote_InvalidJSON(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/quotes", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateQuote(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("CreateQuote() with invalid JSON status = %v, want %v", rec.Code, http.StatusBadRequest)
	}
}

func TestQuoteHandler_GetRandomQuote(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Test with empty database
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quotes/random", nil)
	rec := httptest.NewRecorder()
	handler.GetRandomQuote(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("GetRandomQuote() empty db status = %v, want %v", rec.Code, http.StatusNotFound)
	}

	// Insert a test quote
	body, _ := json.Marshal(map[string]string{
		"text":     "Test quote",
		"author":   "Test Author",
		"category": "motivation",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.CreateQuote(rec, req)

	// Test GetRandomQuote with data
	req = httptest.NewRequest(http.MethodGet, "/api/v1/quotes/random", nil)
	rec = httptest.NewRecorder()
	handler.GetRandomQuote(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GetRandomQuote() with data status = %v, want %v", rec.Code, http.StatusOK)
	}
}

func TestQuoteHandler_GetRandomQuote_WithCategory(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Insert test quotes
	quotes := []map[string]string{
		{"text": "Motivation quote", "author": "Author 1", "category": "motivation"},
		{"text": "Humor quote", "author": "Author 2", "category": "humor"},
	}

	for _, q := range quotes {
		body, _ := json.Marshal(q)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.CreateQuote(rec, req)
	}

	// Test with category filter
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quotes/random?category=motivation", nil)
	rec := httptest.NewRecorder()
	handler.GetRandomQuote(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GetRandomQuote() with category status = %v, want %v", rec.Code, http.StatusOK)
	}

	// Verify the response contains motivation category
	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	if data, ok := response["data"].(map[string]interface{}); ok {
		if data["category"] != "motivation" {
			t.Errorf("GetRandomQuote() category = %v, want motivation", data["category"])
		}
	}
}

func TestQuoteHandler_GetQuotes(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Insert test quotes
	for i := 0; i < 15; i++ {
		body, _ := json.Marshal(map[string]string{
			"text":     "Test quote",
			"author":   "Test Author",
			"category": "test",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.CreateQuote(rec, req)
	}

	tests := []struct {
		name           string
		queryParams    string
		wantStatusCode int
		wantCount      int
	}{
		{
			name:           "default pagination",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantCount:      10,
		},
		{
			name:           "custom limit",
			queryParams:    "?limit=5",
			wantStatusCode: http.StatusOK,
			wantCount:      5,
		},
		{
			name:           "second page",
			queryParams:    "?page=2&limit=10",
			wantStatusCode: http.StatusOK,
			wantCount:      5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/quotes"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			handler.GetQuotes(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("GetQuotes() status = %v, want %v", rec.Code, tt.wantStatusCode)
			}

			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if quotes, ok := data["quotes"].([]interface{}); ok {
					if len(quotes) != tt.wantCount {
						t.Errorf("GetQuotes() count = %v, want %v", len(quotes), tt.wantCount)
					}
				}
			}
		})
	}
}

func TestQuoteHandler_GetCategories(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Test empty categories
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	rec := httptest.NewRecorder()
	handler.GetCategories(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GetCategories() empty status = %v, want %v", rec.Code, http.StatusOK)
	}

	// Insert quotes with different categories
	categories := []string{"motivation", "humor", "wisdom"}
	for _, cat := range categories {
		body, _ := json.Marshal(map[string]string{
			"text":     "Test quote",
			"author":   "Test Author",
			"category": cat,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.CreateQuote(rec, req)
	}

	// Test GetCategories with data
	req = httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	rec = httptest.NewRecorder()
	handler.GetCategories(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GetCategories() with data status = %v, want %v", rec.Code, http.StatusOK)
	}

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	if data, ok := response["data"].([]interface{}); ok {
		if len(data) != 3 {
			t.Errorf("GetCategories() count = %v, want 3", len(data))
		}
	}
}
