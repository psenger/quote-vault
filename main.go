package main

import (
	"log"
	"net/http"

	"quote-vault/config"
	"quote-vault/database"
	"quote-vault/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Initialize(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	r := router.SetupRouter(db, cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}