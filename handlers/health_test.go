package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"quote-vault/database"

	_ "github.com/mattn/go-sqlite3"
)

func setupHealthHandler(t *testing.T) (*HealthHandler, *database.SQLiteDB) {
	db, err := database.NewSQLiteDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	handler := NewHealthHandler(db)
	return handler, db
}

func TestHealthHandler_Health(t *testing.T) {
	handler, db := setupHealthHandler(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Health() status = %v, want %v", rec.Code, http.StatusOK)
	}

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Health() status = %v, want ok", response["status"])
	}

	if response["database"] != "healthy" {
		t.Errorf("Health() database = %v, want healthy", response["database"])
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	handler, db := setupHealthHandler(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rec := httptest.NewRecorder()

	handler.Ready(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Ready() status = %v, want %v", rec.Code, http.StatusOK)
	}

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["status"] != "ready" {
		t.Errorf("Ready() status = %v, want ready", response["status"])
	}
}
