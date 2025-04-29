package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lmousom/passless-auth/internal/errors"
	"github.com/lmousom/passless-auth/internal/middleware"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			middleware.ErrorResponse(w, errors.NewUnauthorized("No authentication token provided", nil))
			return
		}
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid cookie", err))
		return
	}

	tknStr := c.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			middleware.ErrorResponse(w, errors.NewUnauthorized("Invalid token signature", err))
			return
		}
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid token", err))
		return
	}

	if !tkn.Valid {
		middleware.ErrorResponse(w, errors.NewUnauthorized("Invalid token", nil))
		return
	}

	if time.Until(time.Unix(claims.ExpiresAt.Unix(), 0)) > 30*time.Second {
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Token not expired yet", nil))
		return
	}

	// Create new token with extended expiry
	expirationTime := time.Now().Add(24 * time.Hour)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to generate new token", err))
		return
	}

	// Set the new token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	response := map[string]string{
		"message": "Token refreshed successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}
