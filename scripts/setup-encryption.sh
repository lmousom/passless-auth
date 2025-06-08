#!/bin/bash

# Create .env directory if it doesn't exist
mkdir -p .env

# Check if encryption key exists
if [ ! -f .env/encryption.key ]; then
    echo "Generating new encryption key..."
    # Generate a new encryption key using the Go encrypt utility
    go run cmd/encrypt/main.go -generate-key > .env/encryption.key
    
    # Remove any extra output (like "Generated encryption key: ")
    sed -i '' 's/Generated encryption key: //' .env/encryption.key
    
    echo "Encryption key generated and stored in .env/encryption.key"
else
    echo "Using existing encryption key from .env/encryption.key"
fi

# Read the key
ENCRYPTION_KEY=$(cat .env/encryption.key)

# Create .env file for docker-compose
echo "PASSLESS_ENCRYPTION_KEY=$ENCRYPTION_KEY" > .env/docker.env

echo "Encryption key has been set up and is ready for use with docker-compose" 