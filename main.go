package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"quote-vault/database"
	"quote-vault/router"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		os.Exit(0)
	}()

	// Setup router
	r := router.SetupRouter()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}