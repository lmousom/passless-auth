# Passless Auth

A Go-based passwordless authentication system using OTP (One-Time Password) with enterprise-grade security features.

## Features

### Core Features
- OTP Generation and Verification
- JWT-based Authentication
- Session Management
- Refresh Token Support
- Rate Limiting
- Security Headers
- Request Logging
- Configurable Security Settings

### Enterprise Features
- Prometheus Metrics Collection
- Distributed Tracing with OpenTelemetry
- Circuit Breaker for External Services
- Health Checks and Service Status
- Configuration Management with Viper
- Graceful Shutdown
- Metrics Dashboard Integration
- Advanced Logging with Correlation IDs
- Service Discovery Ready
- Container Health Probes

## Project Structure

```
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/    # Request handlers
│   │   └── routes/      # Router setup
│   ├── auth/            # Authentication logic
│   ├── config/          # Configuration
│   ├── middleware/      # Security middleware
│   └── models/          # Data models
├── pkg/                 # Public packages
└── README.md
```

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/lmousom/passless-auth.git
cd passless-auth
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure the application:
Create a configuration file or set environment variables for:
- JWT secret and lifetime
- Server port and environment
- Rate limiting settings
- OTP length and expiry
- Security parameters

4. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Security Features

- Rate limiting: 20 requests per minute with burst of 5
- Secure headers including HSTS, CSP, and XSS protection
- Request logging with timing information
- JWT-based session management
- Cryptographically secure OTP generation

## API Endpoints

- POST `/api/v1/sendOtp` - Send OTP to phone number
- POST `/api/v1/verifyOtp` - Verify OTP
- GET `/api/v1/login` - Check authentication status
- POST `/api/v1/refreshToken` - Refresh JWT token
- POST `/api/v1/logout` - Logout user

## License

MIT License - See [LICENSE](LICENSE) for details





[![Golang API](https://i.imgur.com/jnr7kBu.png)](https://youtu.be/I5WBgYVA8-I)

