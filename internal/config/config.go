package config

import (
	"time"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Port         string        `mapstructure:"port" validate:"required"`
		Environment  string        `mapstructure:"environment" validate:"required,oneof=development staging production"`
		AllowOrigins string        `mapstructure:"allow_origins" validate:"required"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout" validate:"required"`
		WriteTimeout time.Duration `mapstructure:"write_timeout" validate:"required"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout" validate:"required"`
	}

	// JWT configuration
	JWT struct {
		Secret        EncryptedValue `mapstructure:"secret" validate:"required"`
		TokenLifetime time.Duration  `mapstructure:"token_lifetime" validate:"required"`
		Issuer        string         `mapstructure:"issuer" validate:"required"`
	}

	// Security configuration
	Security struct {
		MaxLoginAttempts int           `mapstructure:"max_login_attempts" validate:"required,min=1"`
		LockoutDuration  time.Duration `mapstructure:"lockout_duration" validate:"required"`
		OTPLength        int           `mapstructure:"otp_length" validate:"required,min=4,max=8"`
		OTPExpiry        time.Duration `mapstructure:"otp_expiry" validate:"required"`
		RateLimit        struct {
			RequestsPerMinute int `mapstructure:"requests_per_minute" validate:"required,min=1"`
			BurstSize         int `mapstructure:"burst_size" validate:"required,min=1"`
		} `mapstructure:"rate_limit"`
		TwoFactor struct {
			Enabled   bool   `mapstructure:"enabled" validate:"required"`
			Issuer    string `mapstructure:"issuer" validate:"required"`
			Algorithm string `mapstructure:"algorithm" validate:"required,oneof=SHA1 SHA256 SHA512"`
			Digits    int    `mapstructure:"digits" validate:"required,min=6,max=8"`
			Period    int    `mapstructure:"period" validate:"required,min=30"`
			Skew      int    `mapstructure:"skew" validate:"required,min=1"`
		} `mapstructure:"two_factor"`
	}

	// SMS configuration
	SMS struct {
		Provider   string         `mapstructure:"provider" validate:"required,oneof=twilio"`
		AccountSID EncryptedValue `mapstructure:"account_sid" validate:"required"`
		AuthToken  EncryptedValue `mapstructure:"auth_token" validate:"required"`
		FromNumber string         `mapstructure:"from_number" validate:"required"`
		TemplateID string         `mapstructure:"template_id"`
	}

	// Logging configuration
	Logging struct {
		Level      string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
		Format     string `mapstructure:"format" validate:"required,oneof=json text"`
		OutputPath string `mapstructure:"output_path"`
	}

	// Metrics configuration
	Metrics struct {
		Enabled bool   `mapstructure:"enabled"`
		Port    string `mapstructure:"port" validate:"required_if=Enabled true"`
		Path    string `mapstructure:"path" validate:"required_if=Enabled true"`
	}

	// Tracing configuration
	Tracing struct {
		Enabled     bool   `mapstructure:"enabled"`
		ServiceName string `mapstructure:"service_name" validate:"required_if=Enabled true"`
		Endpoint    string `mapstructure:"endpoint" validate:"required_if=Enabled true"`
	}

	// Redis configuration
	Redis struct {
		Host         string         `mapstructure:"host"`
		Port         string         `mapstructure:"port"`
		Password     string         `mapstructure:"password"`
		DB           int            `mapstructure:"db"`
		PoolSize     int            `mapstructure:"pool_size"`
		MinIdleConns int            `mapstructure:"min_idle_conns"`
		MaxRetries   int            `mapstructure:"max_retries"`
		KeyPrefix    string         `mapstructure:"key_prefix"`
		TTL          RedisTTLConfig `mapstructure:"ttl"`
	}
}

type RedisTTLConfig struct {
	TwoFASecret   time.Duration `mapstructure:"twofa_secret"`
	TwoFAAttempts time.Duration `mapstructure:"twofa_attempts"`
}

// GetDecryptedJWTSecret returns the decrypted JWT secret
func (c *Config) GetDecryptedJWTSecret() (string, error) {
	return c.JWT.Secret.Decrypt()
}

// GetDecryptedSMSAccountSID returns the decrypted SMS account SID
func (c *Config) GetDecryptedSMSAccountSID() (string, error) {
	return c.SMS.AccountSID.Decrypt()
}

// GetDecryptedSMSAuthToken returns the decrypted SMS auth token
func (c *Config) GetDecryptedSMSAuthToken() (string, error) {
	return c.SMS.AuthToken.Decrypt()
}
