#!/bin/bash

# Function to generate a new key
generate_key() {
    echo "Generating new encryption key..."
    go run cmd/encrypt/main.go -generate-key
}

# Function to encrypt a value
encrypt_value() {
    if [ -z "$1" ]; then
        echo "Usage: $0 encrypt 'value-to-encrypt'"
        exit 1
    fi
    go run cmd/encrypt/main.go -key "$PASSLESS_ENCRYPTION_KEY" -value "$1"
}

# Function to show help
show_help() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  generate    Generate a new encryption key"
    echo "  encrypt     Encrypt a value using the current key"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 generate"
    echo "  $0 encrypt 'my-secret-value'"
}

# Main script logic
case "$1" in
    "generate")
        generate_key
        ;;
    "encrypt")
        encrypt_value "$2"
        ;;
    "help"|"")
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        show_help
        exit 1
        ;;
esac 