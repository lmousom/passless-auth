#!/bin/bash

# Make scripts executable
chmod +x scripts/setup-encryption.sh

# Run encryption key setup
./scripts/setup-encryption.sh

# Start the services
docker-compose up --build 