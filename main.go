package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Quote represents an inspirational quote
type Quote struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

// QuoteRequest represents the request body for creating a new quote
type QuoteRequest struct {
	Text     string `json:"text"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

// QuoteResponse represents paginated quote response
type QuoteResponse struct {
	Quotes     []Quote `json:"quotes"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
}

// In-memory storage for quotes
var quotes []Quote
var nextID = 1

func main() {
	// Initialize with some sample quotes
	initializeSampleQuotes()

	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/quotes", addQuoteHandler).Methods("POST")
	r.HandleFunc("/quotes", getQuotesHandler).Methods("GET")
	r.HandleFunc("/quotes/random", getRandomQuoteHandler).Methods("GET")

	fmt.Println("Quote Vault API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Initialize with sample quotes for testing
func initializeSampleQuotes() {
	sampleQuotes := []QuoteRequest{
		{"The only way to do great work is to love what you do.", "Steve Jobs", "motivation"},
		{"Innovation distinguishes between a leader and a follower.", "Steve Jobs", "innovation"},
		{"Life is what happens to you while you're busy making other plans.", "John Lennon", "life"},
		{"The future belongs to those who believe in the beauty of their dreams.", "Eleanor Roosevelt", "dreams"},
	}

	for _, sq := range sampleQuotes {
		quotes = append(quotes, Quote{
			ID:       nextID,
			Text:     sq.Text,
			Author:   sq.Author,
			Category: sq.Category,
		})
		nextID++
	}
}

// Handler to add a new quote
func addQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Text == "" || req.Author == "" {
		http.Error(w, "Text and author are required", http.StatusBadRequest)
		return
	}

	// Set default category if not provided
	if req.Category == "" {
		req.Category = "general"
	}

	// Create new quote
	newQuote := Quote{
		ID:       nextID,
		Text:     req.Text,
		Author:   req.Author,
		Category: req.Category,
	}

	quotes = append(quotes, newQuote)
	nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newQuote)
}

// Handler to get all quotes with pagination
func getQuotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	total := len(quotes)
	totalPages := (total + limit - 1) / limit

	// Calculate start and end indices
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		start = total
		end = total
	} else if end > total {
		end = total
	}

	// Get quotes for current page
	pageQuotes := quotes[start:end]

	response := QuoteResponse{
		Quotes:     pageQuotes,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	json.NewEncoder(w).Encode(response)
}

// Handler to get a random quote
func getRandomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(quotes) == 0 {
		http.Error(w, "No quotes available", http.StatusNotFound)
		return
	}

	category := r.URL.Query().Get("category")

	// Filter quotes by category if specified
	var availableQuotes []Quote
	if category != "" {
		for _, quote := range quotes {
			if quote.Category == category {
				availableQuotes = append(availableQuotes, quote)
			}
		}
		if len(availableQuotes) == 0 {
			http.Error(w, fmt.Sprintf("No quotes found for category '%s'", category), http.StatusNotFound)
			return
		}
	} else {
		availableQuotes = quotes
	}

	// Get random quote
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(availableQuotes))
	randomQuote := availableQuotes[randomIndex]

	json.NewEncoder(w).Encode(randomQuote)
}