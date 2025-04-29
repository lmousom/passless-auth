package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lmousom/passless-auth/internal/errors"
	"github.com/lmousom/passless-auth/internal/middleware"
)

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			middleware.ErrorResponse(w, errors.NewUnauthorized("No authentication token provided", nil))
			return
		}
		middleware.ErrorResponse(w, errors.NewInvalidRequest("Invalid cookie", err))
		return
	}

	tokenStr := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
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

	if !token.Valid {
		middleware.ErrorResponse(w, errors.NewUnauthorized("Invalid token", nil))
		return
	}

	response := map[string]string{
		"message": "Welcome " + claims.Phone + "!",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.ErrorResponse(w, errors.NewInternalServer("Failed to encode response", err))
		return
	}
}
