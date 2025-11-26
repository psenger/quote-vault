package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quote-vault/config"
	"quote-vault/database"
	"quote-vault/middleware"
	"quote-vault/router"
)

func main() {
	// Load configuration
	cfg := config.Load()
	httpCfg := config.GetHTTPConfig()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Setup router with middleware
	r := router.SetupRoutes(db)
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())

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