package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"quote-vault/handlers"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize handlers
	quoteHandler := handlers.NewQuoteHandler()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "quote-vault",
		})
	})

	// Quote routes
	api := r.Group("/api/v1")
	{
		api.POST("/quotes", quoteHandler.AddQuote)
		api.GET("/quotes", quoteHandler.ListQuotes)
		api.GET("/quotes/random", quoteHandler.GetRandomQuote)
		api.GET("/quotes/random/:category", quoteHandler.GetRandomQuoteByCategory)
	}

	// Start server
	log.Println("Starting Quote Vault API server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}