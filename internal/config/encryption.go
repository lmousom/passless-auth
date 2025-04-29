package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	// EncryptionKeyEnv is the environment variable name for the encryption key
	EncryptionKeyEnv = "PASSLESS_ENCRYPTION_KEY"
	// EncryptionKeysEnv is the environment variable name for multiple encryption keys
	EncryptionKeysEnv = "PASSLESS_ENCRYPTION_KEYS"
	// Minimum key length in bytes
	minKeyLength = 32
)

// KeyInfo represents information about an encryption key
type KeyInfo struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

// EncryptedValue represents an encrypted configuration value
type EncryptedValue struct {
	Value   string `mapstructure:"value"`
	KeyID   string `mapstructure:"key_id,omitempty"`
	Version int    `mapstructure:"version,omitempty"`
}

// Decrypt decrypts an encrypted value
func (ev *EncryptedValue) Decrypt() (string, error) {
	if ev.Value == "" {
		return "", nil
	}

	// Check if the value is already decrypted
	if !strings.HasPrefix(ev.Value, "ENC[") || !strings.HasSuffix(ev.Value, "]") {
		return ev.Value, nil
	}

	// Extract the encrypted value
	encrypted := strings.TrimPrefix(strings.TrimSuffix(ev.Value, "]"), "ENC[")

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted value: %w", err)
	}

	// Try to decrypt with the specified key ID first
	if ev.KeyID != "" {
		if key := getKeyByID(ev.KeyID); key != nil {
			if plaintext, err := decryptWithKey(decoded, key); err == nil {
				return plaintext, nil
			}
		}
	}

	// Try all active keys
	keys := getActiveKeys()
	for _, key := range keys {
		if plaintext, err := decryptWithKey(decoded, key); err == nil {
			return plaintext, nil
		}
	}

	return "", fmt.Errorf("failed to decrypt value with any key")
}

// Encrypt encrypts a value
func (ev *EncryptedValue) Encrypt(value string) error {
	if value == "" {
		ev.Value = ""
		return nil
	}

	// Get the primary key
	key := getPrimaryKey()
	if key == nil {
		return fmt.Errorf("no active encryption key found")
	}

	// Decode the key
	keyBytes, err := base64.StdEncoding.DecodeString(key.Key)
	if err != nil {
		return fmt.Errorf("failed to decode key: %w", err)
	}

	// Create cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	// Set encrypted value and key ID
	ev.Value = fmt.Sprintf("ENC[%s]", encoded)
	ev.KeyID = key.ID
	ev.Version = 1

	return nil
}

// decryptWithKey attempts to decrypt a value with a specific key
func decryptWithKey(decoded []byte, key *KeyInfo) (string, error) {
	// Decode the key
	keyBytes, err := base64.StdEncoding.DecodeString(key.Key)
	if err != nil {
		return "", fmt.Errorf("failed to decode key: %w", err)
	}

	// Create cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", fmt.Errorf("encrypted value too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt value: %w", err)
	}

	return string(plaintext), nil
}

// getKeyByID returns a key by its ID
func getKeyByID(id string) *KeyInfo {
	keys := getActiveKeys()
	for _, key := range keys {
		if key.ID == id {
			return key
		}
	}
	return nil
}

// getPrimaryKey returns the primary (most recent) active key
func getPrimaryKey() *KeyInfo {
	keys := getActiveKeys()
	if len(keys) == 0 {
		return nil
	}
	return keys[0]
}

// getActiveKeys returns all active keys sorted by creation date (newest first)
func getActiveKeys() []*KeyInfo {
	// Try to get keys from environment variable first
	if keysStr := os.Getenv(EncryptionKeysEnv); keysStr != "" {
		var keys []*KeyInfo
		if err := json.Unmarshal([]byte(keysStr), &keys); err == nil {
			// Filter active keys and sort by creation date
			var activeKeys []*KeyInfo
			for _, key := range keys {
				if key.Active {
					activeKeys = append(activeKeys, key)
				}
			}
			return activeKeys
		}
	}

	// Fall back to single key
	if keyStr := os.Getenv(EncryptionKeyEnv); keyStr != "" {
		key := &KeyInfo{
			ID:        "default",
			Key:       keyStr,
			CreatedAt: time.Now(),
			Active:    true,
		}
		return []*KeyInfo{key}
	}

	return nil
}

// GenerateEncryptionKey generates a new encryption key
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, minKeyLength)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateKeyID generates a unique key ID
func GenerateKeyID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("key_%x", b)
}

// IsEncrypted checks if a value is encrypted
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, "ENC[") && strings.HasSuffix(value, "]")
}
