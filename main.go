package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"quote-vault/config"
	"quote-vault/database"
	"quote-vault/handlers"
	"quote-vault/repository"
	"quote-vault/router"
	"quote-vault/services"
)

func main() {
	// Load configuration
	cfg := config.Load()
	httpCfg := config.GetHTTPConfig()

	// Initialize database
	db, err := database.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Setup repository, service, and handlers
	quoteRepo := repository.NewQuoteRepository(db.DB())
	quoteService := services.NewQuoteService(quoteRepo)
	quoteHandler := handlers.NewQuoteHandler(quoteService)
	healthHandler := handlers.NewHealthHandler(db)

	// Setup router (middleware is configured inside router)
	r := router.NewRouter(quoteHandler, healthHandler)

	// Configure HTTP server
	srv := &http.Server{
		Addr:           ":" + httpCfg.Port,
		Handler:        r,
		ReadTimeout:    httpCfg.ReadTimeout,
		WriteTimeout:   httpCfg.WriteTimeout,
		IdleTimeout:    httpCfg.IdleTimeout,
		MaxHeaderBytes: httpCfg.MaxHeaderBytes,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", httpCfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), httpCfg.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
