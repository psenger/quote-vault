package router

import (
	"github.com/gin-gonic/gin"
	"quote-vault/handlers"
	"quote-vault/middleware"
)

func SetupRouter(quoteHandler *handlers.QuoteHandler, healthHandler *handlers.HealthHandler) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())

	// Health check endpoint
	r.GET("/health", healthHandler.HealthCheck)

	// API routes
	v1 := r.Group("/api/v1")
	{
		// Quote endpoints
		v1.POST("/quotes", middleware.ValidationMiddleware(), quoteHandler.CreateQuote)
		v1.GET("/quotes/random", quoteHandler.GetRandomQuote)
		v1.GET("/quotes", quoteHandler.GetAllQuotes)
		v1.GET("/quotes/search", quoteHandler.SearchQuotes)
	}

	return r
}