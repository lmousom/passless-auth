package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if time.Until(time.Unix(claims.ExpiresAt.Unix(), 0)) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ttl := 2 * 60 * 1000
	expirationTime := time.Now().UTC().UnixNano()/1000000 + int64(ttl)
	expiredInTime := time.Unix(expirationTime, 0)
	// Now, create a new token for the current use, with a renewed expiration time

	claims.ExpiresAt = jwt.NewNumericDate(expiredInTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expiredInTime,
	})
}
