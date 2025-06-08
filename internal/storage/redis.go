package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/lmousom/passless-auth/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	config *config.Config
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
		config: cfg,
	}, nil
}

// TwoFA operations
func (r *RedisClient) SetTwoFASecret(ctx context.Context, phone, secretKey string) error {
	key := fmt.Sprintf("%stwofa:secret:%s", r.config.Redis.KeyPrefix, phone)
	return r.client.Set(ctx, key, secretKey, r.config.Redis.TTL.TwoFASecret).Err()
}

func (r *RedisClient) GetTwoFASecret(ctx context.Context, phone string) (string, error) {
	key := fmt.Sprintf("%stwofa:secret:%s", r.config.Redis.KeyPrefix, phone)
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) DeleteTwoFASecret(ctx context.Context, phone string) error {
	key := fmt.Sprintf("%stwofa:secret:%s", r.config.Redis.KeyPrefix, phone)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) SetTwoFAEnabled(ctx context.Context, phone string, enabled bool) error {
	key := fmt.Sprintf("%stwofa:enabled:%s", r.config.Redis.KeyPrefix, phone)
	return r.client.Set(ctx, key, fmt.Sprintf("%t", enabled), 0).Err() // No expiration for enabled status
}

func (r *RedisClient) GetTwoFAEnabled(ctx context.Context, phone string) (bool, error) {
	key := fmt.Sprintf("%stwofa:enabled:%s", r.config.Redis.KeyPrefix, phone)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "true", nil
}

func (r *RedisClient) IncrementTwoFAAttempts(ctx context.Context, phone string) (int64, error) {
	key := fmt.Sprintf("%stwofa:attempts:%s", r.config.Redis.KeyPrefix, phone)
	attempts, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration on first attempt
	if attempts == 1 {
		r.client.Expire(ctx, key, r.config.Redis.TTL.TwoFAAttempts)
	}

	return attempts, nil
}

func (r *RedisClient) ResetTwoFAAttempts(ctx context.Context, phone string) error {
	key := fmt.Sprintf("%stwofa:attempts:%s", r.config.Redis.KeyPrefix, phone)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
