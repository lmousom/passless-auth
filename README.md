# Passless Auth

A Go-based passwordless authentication system using OTP (One-Time Password) with enterprise-grade security features.

[![Go Report Card](https://goreportcard.com/badge/github.com/lmousom/passless-auth)](https://goreportcard.com/report/github.com/lmousom/passless-auth)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/lmousom/passless-auth?status.svg)](https://godoc.org/github.com/lmousom/passless-auth)

## 📋 Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Security](#security)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Deployment](#deployment)
- [License](#license)

## ✨ Features

### Core Features
- 🔐 OTP Generation and Verification
- 🔑 JWT-based Authentication
- 📱 Session Management
- 🔄 Refresh Token Support
- 🛡️ Rate Limiting
- 🚦 Security Headers
- 📝 Request Logging
- ⚙️ Configurable Security Settings

### Enterprise Features
- 📊 Prometheus Metrics Collection
- 🔍 Distributed Tracing with OpenTelemetry
- ⚡ Circuit Breaker for External Services
- 🏥 Health Checks and Service Status
- ⚙️ Configuration Management with Viper
- 🛑 Graceful Shutdown
- 📈 Metrics Dashboard Integration
- 📝 Advanced Logging with Correlation IDs
- 🔄 Service Discovery Ready
- 🏥 Container Health Probes

## 🏗️ Architecture

### Project Structure
```
├── cmd/
│   ├── server/          # Application entry point
│   └── encrypt/         # Configuration encryption utility
├── internal/
│   ├── api/
│   │   ├── handlers/    # Request handlers
│   │   └── routes/      # Router setup
│   ├── auth/            # Authentication logic
│   ├── config/          # Configuration
│   ├── middleware/      # Security middleware
│   ├── services/        # External services (SMS, etc.)
│   └── models/          # Data models
├── pkg/                 # Public packages
└── README.md
```

### Component Overview
- **API Layer**: HTTP handlers and route definitions
- **Authentication**: OTP and JWT management
- **Services**: External service integrations
- **Middleware**: Security and logging middleware
- **Configuration**: App configuration and encryption
- **Storage**: Data persistence layer

## 🚀 Getting Started

### Prerequisites
- Go 1.24 or later
- Docker and Docker Compose
- Redis (for session management)
- Twilio account (for SMS)

### Quick Start
1. Clone the repository:
```bash
git clone https://github.com/lmousom/passless-auth.git
cd passless-auth
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up encryption:
```bash
# Generate encryption key
./scripts/manage-keys.sh generate

# Set the encryption key
export PASSLESS_ENCRYPTION_KEY='your-generated-key'
```

4. Start the services:
```bash
docker-compose up --build
```

## ⚙️ Configuration

### Environment Variables
```bash
export PASSLESS_SERVER_PORT=8080
export PASSLESS_JWT_SECRET=your-secret
export PASSLESS_SMS_ACCOUNT_SID=your-sid
export PASSLESS_ENCRYPTION_KEY=your-encryption-key
```

### Configuration File
See `config/config.yaml` for detailed configuration options.

## 🔒 Security

### Key Features
- AES-GCM encryption for sensitive values
- Rate limiting (20 requests/minute)
- Secure headers (HSTS, CSP, XSS)
- JWT-based session management
- Secure OTP generation
- Encrypted configuration
- SMS-based OTP delivery

### Best Practices
1. Never commit encryption keys
2. Use different keys per environment
3. Regular key rotation
4. Monitor access attempts
5. Keep dependencies updated

## 📚 API Documentation

### Endpoints
- `POST /api/v1/sendOtp` - Send OTP
- `POST /api/v1/verifyOtp` - Verify OTP
- `POST /api/v1/2fa/enable` - Enable 2FA
- `POST /api/v1/2fa/verify` - Verify 2FA
- `GET /api/v1/login` - Check auth status
- `POST /api/v1/refreshToken` - Refresh token
- `POST /api/v1/logout` - Logout

### Postman Collection
Import `passless-auth.postman_collection.json` for API testing.

## 💻 Development

### Local Development
1. Start dependencies:
```bash
docker-compose up -d redis prometheus otel-collector
```

2. Run the application:
```bash
go run cmd/server/main.go
```


## 🚢 Deployment

### Docker Deployment
```bash
docker-compose up --build
```


## 📄 License

MIT License - See [LICENSE](LICENSE) for details

---

[![Golang API](https://i.imgur.com/jnr7kBu.png)](https://youtu.be/I5WBgYVA8-I)

