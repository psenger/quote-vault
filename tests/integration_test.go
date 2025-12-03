package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"quote-vault/database"
	"quote-vault/handlers"
	"quote-vault/repository"
	"quote-vault/router"
	"quote-vault/services"
)

func setupTestServer(t *testing.T) (*httptest.Server, *database.SQLiteDB) {
	db, err := database.NewSQLiteDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	repo := repository.NewQuoteRepository(db.DB())
	service := services.NewQuoteService(repo)
	quoteHandler := handlers.NewQuoteHandler(service)
	healthHandler := handlers.NewHealthHandler(db)

	r := router.NewRouter(quoteHandler, healthHandler)
	server := httptest.NewServer(r)

	return server, db
}

func TestIntegration_HealthEndpoint(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /health status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if response["status"] != "ok" {
		t.Errorf("GET /health status = %v, want ok", response["status"])
	}
}

func TestIntegration_ReadyEndpoint(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	resp, err := http.Get(server.URL + "/health/ready")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /health/ready status = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestIntegration_CreateAndGetQuote(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Create a quote
	quoteData := map[string]string{
		"text":     "The only way to do great work is to love what you do.",
		"author":   "Steve Jobs",
		"category": "motivation",
	}
	body, _ := json.Marshal(quoteData)

	resp, err := http.Post(
		server.URL+"/api/v1/quotes",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to create quote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("POST /api/v1/quotes status = %v, want %v", resp.StatusCode, http.StatusCreated)
	}

	var createResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createResponse)

	if !createResponse["success"].(bool) {
		t.Error("POST /api/v1/quotes success = false, want true")
	}

	// Verify the quote was created by getting random quote
	resp, err = http.Get(server.URL + "/api/v1/quotes/random")
	if err != nil {
		t.Fatalf("failed to get random quote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /api/v1/quotes/random status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var randomResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&randomResponse)

	data := randomResponse["data"].(map[string]interface{})
	if data["text"] != quoteData["text"] {
		t.Errorf("GET /api/v1/quotes/random text = %v, want %v", data["text"], quoteData["text"])
	}
}

func TestIntegration_ListQuotesWithPagination(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Create multiple quotes
	for i := 0; i < 15; i++ {
		quoteData := map[string]string{
			"text":     "Test quote text",
			"author":   "Test Author",
			"category": "test",
		}
		body, _ := json.Marshal(quoteData)

		resp, err := http.Post(
			server.URL+"/api/v1/quotes",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatalf("failed to create quote: %v", err)
		}
		resp.Body.Close()
	}

	// Test default pagination
	resp, err := http.Get(server.URL + "/api/v1/quotes")
	if err != nil {
		t.Fatalf("failed to list quotes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /api/v1/quotes status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var listResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&listResponse)

	data := listResponse["data"].(map[string]interface{})
	quotes := data["quotes"].([]interface{})

	if len(quotes) != 10 {
		t.Errorf("GET /api/v1/quotes default count = %v, want 10", len(quotes))
	}

	total := int(data["total"].(float64))
	if total != 15 {
		t.Errorf("GET /api/v1/quotes total = %v, want 15", total)
	}

	// Test custom pagination
	resp, err = http.Get(server.URL + "/api/v1/quotes?page=2&limit=10")
	if err != nil {
		t.Fatalf("failed to list quotes page 2: %v", err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&listResponse)
	data = listResponse["data"].(map[string]interface{})
	quotes = data["quotes"].([]interface{})

	if len(quotes) != 5 {
		t.Errorf("GET /api/v1/quotes page 2 count = %v, want 5", len(quotes))
	}
}

func TestIntegration_FilterByCategory(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Create quotes with different categories
	categories := []string{"motivation", "motivation", "humor", "wisdom"}
	for _, cat := range categories {
		quoteData := map[string]string{
			"text":     "Test quote",
			"author":   "Test Author",
			"category": cat,
		}
		body, _ := json.Marshal(quoteData)

		resp, err := http.Post(
			server.URL+"/api/v1/quotes",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatalf("failed to create quote: %v", err)
		}
		resp.Body.Close()
	}

	// Test filtering by category
	resp, err := http.Get(server.URL + "/api/v1/quotes?category=motivation")
	if err != nil {
		t.Fatalf("failed to filter quotes: %v", err)
	}
	defer resp.Body.Close()

	var listResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&listResponse)

	data := listResponse["data"].(map[string]interface{})
	quotes := data["quotes"].([]interface{})

	if len(quotes) != 2 {
		t.Errorf("GET /api/v1/quotes?category=motivation count = %v, want 2", len(quotes))
	}

	// Verify all returned quotes are in the motivation category
	for _, q := range quotes {
		quote := q.(map[string]interface{})
		if quote["category"] != "motivation" {
			t.Errorf("returned quote category = %v, want motivation", quote["category"])
		}
	}
}

func TestIntegration_GetCategories(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Create quotes with different categories
	categories := []string{"motivation", "humor", "wisdom", "motivation"}
	for _, cat := range categories {
		quoteData := map[string]string{
			"text":     "Test quote",
			"author":   "Test Author",
			"category": cat,
		}
		body, _ := json.Marshal(quoteData)

		resp, err := http.Post(
			server.URL+"/api/v1/quotes",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatalf("failed to create quote: %v", err)
		}
		resp.Body.Close()
	}

	// Get categories
	resp, err := http.Get(server.URL + "/api/v1/categories")
	if err != nil {
		t.Fatalf("failed to get categories: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /api/v1/categories status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var catResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&catResponse)

	data := catResponse["data"].([]interface{})
	if len(data) != 3 {
		t.Errorf("GET /api/v1/categories count = %v, want 3 unique categories", len(data))
	}
}

func TestIntegration_GetRandomQuoteByCategory(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// Create quotes with different categories
	quotesData := []map[string]string{
		{"text": "Motivation quote 1", "author": "Author 1", "category": "motivation"},
		{"text": "Motivation quote 2", "author": "Author 2", "category": "motivation"},
		{"text": "Humor quote", "author": "Author 3", "category": "humor"},
	}

	for _, q := range quotesData {
		body, _ := json.Marshal(q)
		resp, err := http.Post(
			server.URL+"/api/v1/quotes",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatalf("failed to create quote: %v", err)
		}
		resp.Body.Close()
	}

	// Get random quote by category
	resp, err := http.Get(server.URL + "/api/v1/quotes/random?category=motivation")
	if err != nil {
		t.Fatalf("failed to get random quote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /api/v1/quotes/random?category=motivation status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	if data["category"] != "motivation" {
		t.Errorf("random quote category = %v, want motivation", data["category"])
	}

	// Test non-existing category
	resp, err = http.Get(server.URL + "/api/v1/quotes/random?category=nonexistent")
	if err != nil {
		t.Fatalf("failed to get random quote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /api/v1/quotes/random?category=nonexistent status = %v, want %v", resp.StatusCode, http.StatusNotFound)
	}
}

func TestIntegration_ResponseHeaders(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check for security headers
	if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Error("missing X-Content-Type-Options header")
	}

	if resp.Header.Get("X-Frame-Options") != "DENY" {
		t.Error("missing X-Frame-Options header")
	}

	// Check for request ID
	if resp.Header.Get("X-Request-ID") == "" {
		t.Error("missing X-Request-ID header")
	}
}
