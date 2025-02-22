package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lmousom/passless-auth/internal/config"
)

type TokenManager struct {
	config *config.Config
}

func (tm *TokenManager) GenerateToken(phone string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   phone,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.config.JWT.TokenLifetime)),
		Issuer:    "passless-auth",
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.config.JWT.Secret))
}

func (tm *TokenManager) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(tm.config.JWT.Secret), nil
	})
}
