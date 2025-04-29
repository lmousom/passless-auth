package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lmousom/passless-auth/internal/errors"
	"github.com/lmousom/passless-auth/internal/middleware"
	"github.com/lmousom/passless-auth/models/verifydata"
	"github.com/lmousom/passless-auth/utils"
)

var jwtKey = []byte("github.com/lmousom/passless-auth")

type Claims struct {
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

func VerifyOtp(verifyOtpRequest verifydata.VerifyOtpRequest) (*verifydata.VerifyOtpResponse, string, error) {
	if verifyOtpRequest.Phone == "" || verifyOtpRequest.Hash == "" || verifyOtpRequest.Otp == "" {
		return nil, "", errors.NewInvalidRequest("Phone, hash, and OTP are required", nil)
	}

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

	claims := &Claims{
		Phone: verifyOtpRequest.Phone,
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

func VerifyOtpHandler(w http.ResponseWriter, r *http.Request) {
	var verifyOtpRequest verifydata.VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&verifyOtpRequest); err != nil {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid request body", err))
		return
	}

	response, tokenString, err := VerifyOtp(verifyOtpRequest)
	if err != nil {
		middleware.ErrorResponse(w, err)
		return
	}

	// Set the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(24 * time.Hour), // Set a reasonable expiry time
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}
