package main

import (
	"log"
	"net/http"

	"github.com/lmousom/passless-auth/internal/api/routes"
	"github.com/lmousom/passless-auth/internal/config"
)

func main() {
	cfg := &config.Config{} // Initialize config
	router, err := routes.SetupRouter(cfg)
	if err != nil {
		log.Fatalf("Failed to setup router: %v", err)
	}

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
