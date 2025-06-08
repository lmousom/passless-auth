package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lmousom/passless-auth/internal/api/routes"
	"github.com/lmousom/passless-auth/internal/config"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize configuration manager
	cfgManager, err := config.NewConfigManager()
	if err != nil {
		log.Fatalf("Failed to initialize configuration manager: %v", err)
	}
	defer cfgManager.Close()

	// Get initial configuration
	cfg := cfgManager.GetConfig()

	// Setup router
	router, err := routes.SetupRouter(cfg)
	if err != nil {
		log.Fatalf("Failed to setup router: %v", err)
	}

	// Wrap the router with error handling middleware
	handler := middleware.ErrorHandler(router)

	// Create server with configuration
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start metrics server if enabled
	if cfg.Metrics.Enabled {
		metricsMux := http.NewServeMux()
		metricsMux.Handle(cfg.Metrics.Path, promhttp.Handler())
		metricsServer := &http.Server{
			Addr:    ":" + cfg.Metrics.Port,
			Handler: metricsMux,
		}
		go func() {
			log.Printf("Metrics server starting on :%s%s", cfg.Metrics.Port, cfg.Metrics.Path)
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Metrics server failed to start: %v", err)
			}
		}()
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on :%s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Subscribe to configuration changes
	configChan := cfgManager.Subscribe()
	defer cfgManager.Unsubscribe(configChan)

	// Handle configuration changes
	go func() {
		for newCfg := range configChan {
			log.Printf("Configuration updated, restarting server...")

			// Create new server with updated configuration
			newServer := &http.Server{
				Addr:         ":" + newCfg.Server.Port,
				Handler:      handler,
				ReadTimeout:  newCfg.Server.ReadTimeout,
				WriteTimeout: newCfg.Server.WriteTimeout,
				IdleTimeout:  newCfg.Server.IdleTimeout,
			}

			// Gracefully shutdown old server
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down server: %v", err)
			}
			cancel()

			// Start new server
			server = newServer
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("Server failed to start: %v", err)
				}
			}()
			log.Printf("Server restarted on :%s", newCfg.Server.Port)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
