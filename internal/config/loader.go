package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// LoadConfig loads the configuration from multiple sources
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Set default values
	setDefaults(v)

	// Add configuration paths
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/passless-auth/")

	// Read environment variables
	v.SetEnvPrefix("PASSLESS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate config
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.allow_origins", "*")
	v.SetDefault("server.read_timeout", "5s")
	v.SetDefault("server.write_timeout", "10s")
	v.SetDefault("server.idle_timeout", "120s")

	// JWT defaults
	v.SetDefault("jwt.token_lifetime", "24h")
	v.SetDefault("jwt.issuer", "passless-auth")

	// Security defaults
	v.SetDefault("security.max_login_attempts", 3)
	v.SetDefault("security.lockout_duration", "15m")
	v.SetDefault("security.otp_length", 6)
	v.SetDefault("security.otp_expiry", "5m")
	v.SetDefault("security.rate_limit.requests_per_minute", 20)
	v.SetDefault("security.rate_limit.burst_size", 5)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Metrics defaults
	v.SetDefault("metrics.enabled", false)
	v.SetDefault("metrics.port", "9090")
	v.SetDefault("metrics.path", "/metrics")

	// Tracing defaults
	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.service_name", "passless-auth")
}

// validateConfig validates the configuration using go-playground/validator
func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return err
	}
	return nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() string {
	// Check environment variable
	if path := os.Getenv("PASSLESS_CONFIG_PATH"); path != "" {
		return path
	}

	// Check common locations
	paths := []string{
		"config.yaml",
		"config.yml",
		"./config/config.yaml",
		"./config/config.yml",
		"/etc/passless-auth/config.yaml",
		"/etc/passless-auth/config.yml",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default config
	v := viper.New()
	setDefaults(v)

	// Write config file
	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
