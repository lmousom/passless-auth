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
│   ├── server/          # Application entry point
│   └── encrypt/         # Configuration encryption utility
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

## Configuration Management

### Configuration File

The application uses YAML configuration files. A default configuration file is provided at `config/config.yaml`:

```yaml
# Server configuration
server:
  port: "8080"
  environment: "development"
  allow_origins: "*"
  read_timeout: "5s"
  write_timeout: "10s"
  idle_timeout: "120s"

# JWT configuration
jwt:
  secret:
    value: "ENC[your-encrypted-jwt-secret-here]"
  token_lifetime: "24h"
  issuer: "passless-auth"

# Security configuration
security:
  max_login_attempts: 3
  lockout_duration: "15m"
  otp_length: 6
  otp_expiry: "5m"
  rate_limit:
    requests_per_minute: 20
    burst_size: 5

# SMS configuration
sms:
  provider: "twilio"
  account_sid:
    value: "ENC[your-encrypted-account-sid-here]"
  auth_token:
    value: "ENC[your-encrypted-auth-token-here]"
  from_number: "+1234567890"
  template_id: ""

# Logging configuration
logging:
  level: "info"
  format: "json"
  output_path: ""

# Metrics configuration
metrics:
  enabled: false
  port: "9090"
  path: "/metrics"

# Tracing configuration
tracing:
  enabled: false
  service_name: "passless-auth"
  endpoint: "http://localhost:4317"
```

### Environment Variables

Configuration can be overridden using environment variables:

```bash
export PASSLESS_SERVER_PORT=8080
export PASSLESS_JWT_SECRET=your-secret
export PASSLESS_SMS_ACCOUNT_SID=your-sid
```

### Configuration Hot-Reloading

The application supports hot-reloading of configuration changes:

1. Changes to the configuration file are automatically detected
2. The server gracefully restarts with the new configuration
3. Existing connections are preserved
4. Invalid configurations are rejected

## Encryption System

### Key Management

The application uses AES-GCM encryption for sensitive values. Keys can be managed using the encryption utility:

1. Generate a new key:
```bash
go run cmd/encrypt/main.go -generate-key
```

2. Rotate keys:
```bash
go run cmd/encrypt/main.go -rotate-key
```

3. Encrypt a value:
```bash
go run cmd/encrypt/main.go -value "your-secret-value"
```

### Key Rotation

The system supports multiple active keys for smooth rotation:

1. Set multiple keys using environment variable:
```bash
export PASSLESS_ENCRYPTION_KEYS='[
  {
    "id": "key_123",
    "key": "base64-encoded-key",
    "created_at": "2024-03-14T12:00:00Z",
    "active": true
  }
]'
```

2. Or use a single key:
```bash
export PASSLESS_ENCRYPTION_KEY="your-base64-encoded-key"
```

### Encrypted Values

Sensitive values in the configuration are encrypted using the format:
```
ENC[base64-encoded-encrypted-value]
```

The system automatically:
- Encrypts new values with the primary key
- Decrypts values using the appropriate key
- Falls back to older keys if needed
- Tracks key usage with versioning

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

