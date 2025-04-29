package main

import (
	"log"
	"net/http"

	"github.com/lmousom/passless-auth/internal/api/routes"
	"github.com/lmousom/passless-auth/internal/config"
	"github.com/lmousom/passless-auth/internal/middleware"
)

func main() {
	cfg := &config.Config{} // Initialize config
	router, err := routes.SetupRouter(cfg)
	if err != nil {
		log.Fatalf("Failed to setup router: %v", err)
	}

	// Wrap the router with error handling middleware
	handler := middleware.ErrorHandler(router)

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
