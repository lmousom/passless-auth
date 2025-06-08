#!/bin/sh

# Check if encryption key is set in environment
if [ -z "$PASSLESS_ENCRYPTION_KEY" ]; then
    echo "Error: PASSLESS_ENCRYPTION_KEY environment variable is not set"
    echo "Please set the encryption key using:"
    echo "export PASSLESS_ENCRYPTION_KEY='your-encryption-key'"
    exit 1
fi

echo "Encryption key is set and ready"

# Execute the main application
exec "$@" 