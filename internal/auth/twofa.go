package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/lmousom/passless-auth/internal/config"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TwoFAManager struct {
	config *config.Config
}

func NewTwoFAManager(cfg *config.Config) *TwoFAManager {
	return &TwoFAManager{
		config: cfg,
	}
}

func (tm *TwoFAManager) GenerateSecretKey(phone string) (string, error) {
	// Generate a random secret key
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate secret key: %w", err)
	}

	// Encode the secret key in base32
	secretKey := base32.StdEncoding.EncodeToString(secret)

	// Generate TOTP configuration
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      tm.config.Security.TwoFactor.Issuer,
		AccountName: phone,
		Secret:      []byte(secretKey),
		Algorithm:   otp.AlgorithmSHA1,
		Digits:      otp.Digits(tm.config.Security.TwoFactor.Digits),
		Period:      uint(tm.config.Security.TwoFactor.Period),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	return key.Secret(), nil
}

func (tm *TwoFAManager) GenerateQRCode(phone, secretKey string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      tm.config.Security.TwoFactor.Issuer,
		AccountName: phone,
		Secret:      []byte(secretKey),
		Algorithm:   otp.AlgorithmSHA1,
		Digits:      otp.Digits(tm.config.Security.TwoFactor.Digits),
		Period:      uint(tm.config.Security.TwoFactor.Period),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	return key.URL(), nil
}

func (tm *TwoFAManager) ValidateCode(secretKey, code string) bool {
	return totp.Validate(code, secretKey)
}

func (tm *TwoFAManager) GenerateCode(secretKey string) (string, error) {
	code, err := totp.GenerateCode(secretKey, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP code: %w", err)
	}
	return code, nil
}
