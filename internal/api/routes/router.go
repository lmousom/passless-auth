package routes

import (
	"github.com/gorilla/mux"
	handlers "github.com/lmousom/passless-auth/internal/api/handlers"
	"github.com/lmousom/passless-auth/internal/auth"
	"github.com/lmousom/passless-auth/internal/config"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/lmousom/passless-auth/internal/services/sms"
	"github.com/lmousom/passless-auth/internal/storage"
)

func SetupRouter(cfg *config.Config) (*mux.Router, error) {
	r := mux.NewRouter()

	// Middleware
	rateLimiter, err := middleware.RateLimiter()
	if err != nil {
		return nil, err
	}

	// Apply security middleware
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RequestLogger)
	r.Use(rateLimiter)

	// Apply metrics middleware if enabled
	if cfg.Metrics.Enabled {
		r.Use(middleware.MetricsMiddleware)
	}

	// Apply tracing middleware if enabled
	if cfg.Tracing.Enabled {
		r.Use(middleware.TracingMiddleware)
	}

	// Initialize services
	smsService, err := sms.NewTwilioService(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize Redis client
	redisClient, err := storage.NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	sendOtpHandler := handlers.NewSendOtpHandler(smsService)
	twoFAManager := auth.NewTwoFAManager(cfg)
	verifyOtpHandler := handlers.NewVerifyOtpHandler(redisClient, twoFAManager)
	twoFAHandler := handlers.NewTwoFAHandler(twoFAManager, redisClient)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/sendOtp", sendOtpHandler.Handle).Methods("POST")
	api.HandleFunc("/verifyOtp", verifyOtpHandler.Handle).Methods("POST")
	api.HandleFunc("/login", handlers.VerificationHandler).Methods("GET")
	api.HandleFunc("/refreshToken", handlers.RefreshTokenHandler).Methods("POST")
	api.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	api.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// 2FA routes
	api.HandleFunc("/2fa/enable", twoFAHandler.Enable2FA).Methods("POST")
	api.HandleFunc("/2fa/verify", twoFAHandler.Verify2FA).Methods("POST")
	api.HandleFunc("/2fa/disable", twoFAHandler.Disable2FA).Methods("POST")

	return r, nil
}
