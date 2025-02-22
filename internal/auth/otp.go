package auth

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/lmousom/passless-auth/internal/config"
)

type OTPManager struct {
	config *config.Config
}

func (om *OTPManager) GenerateOTP() (string, error) {
	buffer := make([]byte, om.config.Security.OTPLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(buffer)[:om.config.Security.OTPLength], nil
}
