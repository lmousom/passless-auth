package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lmousom/passless-auth/internal/auth"
	"github.com/lmousom/passless-auth/internal/errors"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/lmousom/passless-auth/internal/storage"
	"github.com/lmousom/passless-auth/models/verifydata"
	"github.com/lmousom/passless-auth/utils"
)

var jwtKey = []byte("github.com/lmousom/passless-auth")

type Claims struct {
	Phone         string `json:"phone"`
	TwoFAEnabled  bool   `json:"twofa_enabled"`
	TwoFAVerified bool   `json:"twofa_verified"`
	jwt.RegisteredClaims
}

type VerifyOtpHandler struct {
	redisClient  *storage.RedisClient
	twoFAManager *auth.TwoFAManager
}

func NewVerifyOtpHandler(redisClient *storage.RedisClient, twoFAManager *auth.TwoFAManager) *VerifyOtpHandler {
	return &VerifyOtpHandler{
		redisClient:  redisClient,
		twoFAManager: twoFAManager,
	}
}

func (h *VerifyOtpHandler) VerifyOtp(verifyOtpRequest verifydata.VerifyOtpRequest) (*verifydata.VerifyOtpResponse, string, error) {
	if verifyOtpRequest.Phone == "" || verifyOtpRequest.Hash == "" || verifyOtpRequest.Otp == "" {
		return nil, "", errors.NewInvalidRequest("Phone, hash, and OTP are required", nil)
	}

	// First validate the OTP
	f := func(c rune) bool {
		return c == '.'
	}
	extValue := strings.FieldsFunc(verifyOtpRequest.Hash, f)
	if len(extValue) != 2 {
		return nil, "", errors.NewInvalidRequest("Invalid hash format", nil)
	}

	hashValue := extValue[0]
	expiresIn := extValue[1]

	now := time.Now()
	expiredInTime, err := utils.MsToTime(expiresIn)
	if err != nil {
		return nil, "", errors.NewInvalidRequest("Invalid expiry time", err)
	}

	if now.After(expiredInTime) {
		return nil, "", errors.NewOTPExpired("OTP has expired", nil)
	}

	if hashValue != utils.Encrypt([]byte(verifyOtpRequest.Phone+"."+verifyOtpRequest.Otp+"."+expiresIn)) {
		return nil, "", errors.NewInvalidOTP("Invalid OTP", nil)
	}

	// Check if 2FA is enabled
	ctx := context.Background()
	twoFAEnabled, err := h.redisClient.GetTwoFAEnabled(ctx, verifyOtpRequest.Phone)
	if err != nil {
		return nil, "", errors.NewInternalServer("Failed to check 2FA status", err)
	}

	// For now, we'll just set TwoFAVerified to false if 2FA is enabled
	// The actual 2FA verification should happen in a separate endpoint
	twoFAVerified := !twoFAEnabled

	claims := &Claims{
		Phone:         verifyOtpRequest.Phone,
		TwoFAEnabled:  twoFAEnabled,
		TwoFAVerified: twoFAVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredInTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, "", errors.NewInternalServer("Failed to generate token", err)
	}

	return &verifydata.VerifyOtpResponse{
		Status:  "success",
		Message: "OTP verified successfully",
	}, tokenString, nil
}

func (h *VerifyOtpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var verifyOtpRequest verifydata.VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&verifyOtpRequest); err != nil {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid request body", err))
		return
	}

	response, tokenString, err := h.VerifyOtp(verifyOtpRequest)
	if err != nil {
		middleware.ErrorResponse(w, err)
		return
	}

	// Set the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(24 * time.Hour), // Set a reasonable expiry time
		Path:    "/api/v1",
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}
