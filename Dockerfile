# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -o passless-auth ./cmd/server/main.go

# Build the encrypt utility
RUN CGO_ENABLED=0 GOOS=linux go build -o encrypt ./cmd/encrypt/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary tools
RUN apk add --no-cache bash

# Copy the binaries from builder
COPY --from=builder /app/passless-auth .
COPY --from=builder /app/encrypt .

# Copy config files
COPY config/ ./config/

# Copy and set up entrypoint script
COPY scripts/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Expose the application port
EXPOSE 8080 9090

# Set the entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]

# Run the application
CMD ["./passless-auth"] 