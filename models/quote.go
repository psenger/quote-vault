package models

import (
	"time"
)

// Quote represents an inspirational quote with metadata
type Quote struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// QuoteRequest represents the payload for creating a new quote
type QuoteRequest struct {
	Text     string `json:"text" binding:"required"`
	Author   string `json:"author" binding:"required"`
	Category string `json:"category" binding:"required"`
}

// QuoteResponse represents the response structure for quote operations
type QuoteResponse struct {
	Quote   *Quote `json:"quote,omitempty"`
	Message string `json:"message,omitempty"`
}

// QuotesListResponse represents the response for listing quotes with pagination
type QuotesListResponse struct {
	Quotes []Quote `json:"quotes"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}