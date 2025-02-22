package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	viper.SetDefault("jwt.lifetime", "24h")
	viper.SetDefault("security.otp_length", 6)
	viper.SetDefault("security.otp_expiry", "5m")
	viper.SetDefault("security.max_attempts", 3)
	viper.SetDefault("server.port", "8080")

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
