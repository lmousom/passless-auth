package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lmousom/passless-auth/internal/config"
)

func main() {
	// Parse command line flags
	value := flag.String("value", "", "Value to encrypt")
	key := flag.String("key", "", "Encryption key (base64 encoded)")
	generateKey := flag.Bool("generate-key", false, "Generate a new encryption key")
	rotateKey := flag.Bool("rotate-key", false, "Rotate encryption keys")
	flag.Parse()

	// Generate new key if requested
	if *generateKey {
		newKey, err := config.GenerateEncryptionKey()
		if err != nil {
			fmt.Printf("Failed to generate key: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated encryption key: %s\n", newKey)
		fmt.Printf("Set this key as the %s environment variable\n", config.EncryptionKeyEnv)
		return
	}

	// Rotate keys if requested
	if *rotateKey {
		if err := rotateKeys(); err != nil {
			fmt.Printf("Failed to rotate keys: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Check if value is provided
	if *value == "" {
		fmt.Println("Error: value is required")
		flag.Usage()
		os.Exit(1)
	}

	// Set encryption key if provided
	if *key != "" {
		os.Setenv(config.EncryptionKeyEnv, *key)
	}

	// Create encrypted value
	ev := &config.EncryptedValue{}
	if err := ev.Encrypt(*value); err != nil {
		fmt.Printf("Failed to encrypt value: %v\n", err)
		os.Exit(1)
	}

	// Print encrypted value
	fmt.Printf("Encrypted value: %s\n", ev.Value)
	if ev.KeyID != "" {
		fmt.Printf("Key ID: %s\n", ev.KeyID)
	}
}

// rotateKeys rotates the encryption keys
func rotateKeys() error {
	// Generate new key
	newKey, err := config.GenerateEncryptionKey()
	if err != nil {
		return fmt.Errorf("failed to generate new key: %w", err)
	}

	// Create new key info
	newKeyInfo := &config.KeyInfo{
		ID:        config.GenerateKeyID(),
		Key:       newKey,
		CreatedAt: time.Now(),
		Active:    true,
	}

	// Get existing keys
	var keys []*config.KeyInfo
	if keysStr := os.Getenv(config.EncryptionKeysEnv); keysStr != "" {
		if err := json.Unmarshal([]byte(keysStr), &keys); err != nil {
			return fmt.Errorf("failed to parse existing keys: %w", err)
		}
	}

	// Add new key
	keys = append([]*config.KeyInfo{newKeyInfo}, keys...)

	// Keep only active keys
	var activeKeys []*config.KeyInfo
	for _, key := range keys {
		if key.Active {
			activeKeys = append(activeKeys, key)
		}
	}

	// Marshal keys
	keysJSON, err := json.Marshal(activeKeys)
	if err != nil {
		return fmt.Errorf("failed to marshal keys: %w", err)
	}

	// Print new key info
	fmt.Printf("Generated new key with ID: %s\n", newKeyInfo.ID)
	fmt.Printf("Set the following as the %s environment variable:\n", config.EncryptionKeysEnv)
	fmt.Println(string(keysJSON))

	return nil
}
