package config

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ConfigManager manages configuration and supports hot-reloading
type ConfigManager struct {
	config      *Config
	viper       *viper.Viper
	mu          sync.RWMutex
	subscribers []chan *Config
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() (*ConfigManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cm := &ConfigManager{
		viper:       viper.New(),
		subscribers: make([]chan *Config, 0),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Set up viper
	cm.viper.SetConfigName("config")
	cm.viper.SetConfigType("yaml")
	cm.viper.AddConfigPath(".")
	cm.viper.AddConfigPath("./config")
	cm.viper.AddConfigPath("/etc/passless-auth/")

	// Set up environment variables
	cm.viper.SetEnvPrefix("PASSLESS")
	cm.viper.AutomaticEnv()

	// Set defaults
	setDefaults(cm.viper)

	// Load initial configuration
	if err := cm.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load initial configuration: %w", err)
	}

	// Set up file watcher
	cm.viper.WatchConfig()
	cm.viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuration file changed: %s", e.Name)
		if err := cm.loadConfig(); err != nil {
			log.Printf("Failed to reload configuration: %v", err)
			return
		}
		cm.notifySubscribers()
	})

	return cm, nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// Subscribe returns a channel that will receive configuration updates
func (cm *ConfigManager) Subscribe() <-chan *Config {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	ch := make(chan *Config, 1)
	cm.subscribers = append(cm.subscribers, ch)
	return ch
}

// Unsubscribe removes a subscriber
func (cm *ConfigManager) Unsubscribe(ch <-chan *Config) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for i, subscriber := range cm.subscribers {
		if subscriber == ch {
			cm.subscribers = append(cm.subscribers[:i], cm.subscribers[i+1:]...)
			close(subscriber)
			return
		}
	}
}

// Close stops the configuration manager
func (cm *ConfigManager) Close() {
	cm.cancel()
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, ch := range cm.subscribers {
		close(ch)
	}
	cm.subscribers = nil
}

// loadConfig loads the configuration from the current source
func (cm *ConfigManager) loadConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Read config file
	if err := cm.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	var cfg Config
	if err := cm.viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate config
	if err := validateConfig(&cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	cm.config = &cfg
	return nil
}

// notifySubscribers notifies all subscribers of configuration changes
func (cm *ConfigManager) notifySubscribers() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, ch := range cm.subscribers {
		select {
		case ch <- cm.config:
		default:
			// Skip if channel is full
		}
	}
}
