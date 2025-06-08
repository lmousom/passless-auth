package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lmousom/passless-auth/internal/auth"
	"github.com/lmousom/passless-auth/internal/errors"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/lmousom/passless-auth/internal/storage"
	"github.com/lmousom/passless-auth/models/twofa"
)

type TwoFAHandler struct {
	twoFAManager *auth.TwoFAManager
	redisClient  *storage.RedisClient
}

func NewTwoFAHandler(twoFAManager *auth.TwoFAManager, redisClient *storage.RedisClient) *TwoFAHandler {
	return &TwoFAHandler{
		twoFAManager: twoFAManager,
		redisClient:  redisClient,
	}
}

func (h *TwoFAHandler) Enable2FA(w http.ResponseWriter, r *http.Request) {
	var req twofa.Enable2FARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid request body", err))
		return
	}

	// Check if 2FA is already enabled
	ctx := r.Context()
	enabled, err := h.redisClient.GetTwoFAEnabled(ctx, req.Phone)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to check 2FA status", err))
		return
	}
	if enabled {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("2FA is already enabled", nil))
		return
	}

	// Generate secret key
	secretKey, err := h.twoFAManager.GenerateSecretKey(req.Phone)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to generate 2FA secret key", err))
		return
	}

	// Store secret key
	if err := h.redisClient.SetTwoFASecret(ctx, req.Phone, secretKey); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to store 2FA secret key", err))
		return
	}

	// Set 2FA enabled status
	if err := h.redisClient.SetTwoFAEnabled(ctx, req.Phone, true); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to enable 2FA", err))
		return
	}

	// Generate QR code
	qrCode, err := h.twoFAManager.GenerateQRCode(req.Phone, secretKey)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to generate QR code", err))
		return
	}

	response := &twofa.Enable2FAResponse{
		Status:    "success",
		Message:   "2FA setup initiated",
		SecretKey: secretKey,
		QRCode:    qrCode,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}

func (h *TwoFAHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	var req twofa.Verify2FARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid request body", err))
		return
	}

	ctx := r.Context()

	// Check if 2FA is enabled
	enabled, err := h.redisClient.GetTwoFAEnabled(ctx, req.Phone)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to check 2FA status", err))
		return
	}
	if !enabled {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("2FA is not enabled", nil))
		return
	}

	// Get secret key
	secretKey, err := h.redisClient.GetTwoFASecret(ctx, req.Phone)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to get 2FA secret key", err))
		return
	}

	// Check attempts
	attempts, err := h.redisClient.IncrementTwoFAAttempts(ctx, req.Phone)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to track 2FA attempts", err))
		return
	}
	if attempts > 3 {
		middleware.ErrorResponse(w, errors.NewTooManyAttempts("Too many 2FA attempts", nil))
		return
	}

	// Validate code
	if !h.twoFAManager.ValidateCode(secretKey, req.Code) {
		middleware.ErrorResponse(w, errors.NewInvalidOTP("Invalid 2FA code", nil))
		return
	}

	// Reset attempts on successful verification
	if err := h.redisClient.ResetTwoFAAttempts(ctx, req.Phone); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to reset 2FA attempts", err))
		return
	}

	// Generate new token with TwoFAVerified set to true
	claims := &Claims{
		Phone:         req.Phone,
		TwoFAEnabled:  true,
		TwoFAVerified: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to generate token", err))
		return
	}

	// Set the new token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/api/v1",
	})

	response := &twofa.Verify2FAResponse{
		Status:  "success",
		Message: "2FA code verified successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}

func (h *TwoFAHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	var req twofa.Disable2FARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid request body", err))
		return
	}

	ctx := r.Context()

	// Check if 2FA is enabled
	enabled, err := h.redisClient.GetTwoFAEnabled(ctx, req.Phone)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to check 2FA status", err))
		return
	}
	if !enabled {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("2FA is not enabled", nil))
		return
	}

	// Get secret key
	secretKey, err := h.redisClient.GetTwoFASecret(ctx, req.Phone)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to get 2FA secret key", err))
		return
	}

	// Validate code
	if !h.twoFAManager.ValidateCode(secretKey, req.Code) {
		middleware.ErrorResponse(w, errors.NewInvalidOTP("Invalid 2FA code", nil))
		return
	}

	// Disable 2FA
	if err := h.redisClient.SetTwoFAEnabled(ctx, req.Phone, false); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to disable 2FA", err))
		return
	}

	// Delete secret key
	if err := h.redisClient.DeleteTwoFASecret(ctx, req.Phone); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to delete 2FA secret key", err))
		return
	}

	response := &twofa.Disable2FAResponse{
		Status:  "success",
		Message: "2FA disabled successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}
