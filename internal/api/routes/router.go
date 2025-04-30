package routes

import (
	"github.com/gorilla/mux"
	handlers "github.com/lmousom/passless-auth/internal/api/handlers"
	"github.com/lmousom/passless-auth/internal/config"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/lmousom/passless-auth/internal/services/sms"
)

func SetupRouter(cfg *config.Config) (*mux.Router, error) {
	r := mux.NewRouter()

	// Middleware
	rateLimiter, err := middleware.RateLimiter()
	if err != nil {
		return nil, err
	}

	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RequestLogger)
	r.Use(rateLimiter)

	// Initialize services
	smsService, err := sms.NewTwilioService(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	sendOtpHandler := handlers.NewSendOtpHandler(smsService)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/sendOtp", sendOtpHandler.Handle).Methods("POST")
	api.HandleFunc("/verifyOtp", handlers.VerifyOtpHandler).Methods("POST")
	api.HandleFunc("/login", handlers.VerificationHandler).Methods("GET")
	api.HandleFunc("/refreshToken", handlers.RefreshTokenHandler).Methods("POST")
	api.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	api.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	return r, nil
}
