package config

import (
	"time"
)

type Config struct {
	JWT struct {
		Secret        string        `mapstructure:"secret"`
		TokenLifetime time.Duration `mapstructure:"token_lifetime"`
	}
	Server struct {
		Port         string `mapstructure:"port"`
		Environment  string `mapstructure:"environment"`
		AllowOrigins string `mapstructure:"allow_origins"`
	}
	Security struct {
		MaxLoginAttempts int           `mapstructure:"max_login_attempts"`
		LockoutDuration  time.Duration `mapstructure:"lockout_duration"`
		OTPLength        int           `mapstructure:"otp_length"`
		OTPExpiry        time.Duration `mapstructure:"otp_expiry"`
	}
}
