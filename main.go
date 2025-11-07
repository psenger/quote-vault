package main

import (
	"log"
	"net/http"
	"quote-vault/database"
	"quote-vault/router"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Setup routes
	r := router.SetupRoutes()

	// Start server
	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}