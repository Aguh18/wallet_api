# Wallet API

Simple REST API wallet service built with Go and Clean Architecture principles.

[![Bruno Collection](https://img.shields.io/badge/API_Testing-Browser-blue?logo=usebruno)](docs/api/)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

> **🚀 Quick Start**: [Bruno API Collection](docs/api/) available for testing all endpoints!

## Features

- **User Authentication**
  - JWT-based authentication (access + refresh tokens)
  - HttpOnly cookie-based token storage (XSS protection)
  - Secure password hashing with bcrypt
  - User registration and login

- **Account Management**
  - Create multiple accounts per user
  - Support for multiple currencies
  - Account status management (active, inactive, frozen)

- **Transaction Processing**
  - Deposit and withdraw funds
  - Transfer funds between accounts
  - Transaction history with pagination
  - Balance tracking (before/after transaction)
  - Idempotency support with reference ID
  - Pessimistic locking (SELECT FOR UPDATE) to prevent race conditions
  - Atomic transactions for data consistency

- **Security**
  - Cookie-based authentication (HttpOnly, Secure, SameSite)
  - Password hashing with bcrypt (cost 12)
  - JWT access tokens (15 min expiry)
  - JWT refresh tokens (7 days expiry)
  - Input validation

- **Architecture**
  - Clean Architecture with pragmatic approach
  - Modular structure (user, account modules)
  - Interface-based dependency injection
  - Generic base repository pattern
  - Request/Response DTOs
  - Bruno API collection for testing

## Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Fiber v2.52.10
- **ORM**: GORM v1.31.1
- **Database**: PostgreSQL 18
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Migrations**: golang-migrate/migrate
- **Password Hashing**: bcrypt
- **Hot Reload**: Air (development)
- **API Testing**: Bruno collection included

## Project Structure

```
wallet_api/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go               # Application initialization
│   │   └── migrate.go           # Database migrations
│   ├── common/
│   │   ├── base/
│   │   │   └── base.repository.go   # Generic repository pattern
│   │   ├── consts/
│   │   │   └── consts.go            # Application constants
│   │   ├── errors/
│   │   │   └── error.go             # Error types
│   │   └── response/
│   │       └── response.go          # API response helpers
│   ├── entity/
│   │   ├── user.go               # User entity
│   │   ├── account.go            # Account entity
│   │   ├── transaction.go        # Transaction entity
│   │   ├── session.go            # Session entity
│   │   └── access_token.go       # Access token entity
│   ├── middleware/
│   │   ├── auth.go               # JWT authentication middleware
│   │   ├── logger.go             # HTTP request logging
│   │   └── recovery.go           # Panic recovery
│   ├── module/
│   │   ├── account/              # Account module
│   │   │   ├── account.module.go
│   │   │   ├── account.router.go
│   │   │   ├── handler/
│   │   │   │   └── account.handler.go
│   │   │   ├── usecase/
│   │   │   │   └── account.usecase.go
│   │   │   ├── repository/
│   │   │   │   ├── account.repository.go
│   │   │   │   └── transaction.repository.go
│   │   │   └── dto/
│   │   │       ├── request/
│   │   │       │   └── account.request.go
│   │   │       └── response/
│   │   │           └── account.response.go
│   │   └── user/                 # User module
│   │       ├── user.module.go
│   │       ├── user.router.go
│   │       ├── handler/
│   │       │   └── user.handler.go
│   │       ├── usecase/
│   │       │   └── user.usecase.go
│   │       ├── repository/
│   │       │   └── user.repository.go
│   │       └── dto/
│   │           ├── request/
│   │           │   └── user.request.go
│   │           └── response/
│   │               └── user.response.go
│   ├── router/
│   │   └── module.go             # Global router initialization
│   └── utils/
│       ├── jwt.go                # JWT utilities
│       ├── cookie.go             # Cookie management
│       └── password.go           # Password hashing
├── migrations/                   # Database migrations
├── .air.toml                    # Air hot reload config
├── docker-compose.yml           # Docker services
├── Dockerfile                    # Multi-stage Docker build
├── Makefile                     # Build & run commands
└── go.mod                       # Go dependencies
```

## Quick Start

### Test API with Bruno

Don't want to code? Test all endpoints immediately!

1. Install [Bruno](https://www.usebruno.com/)
2. Import [Bruno Collection](docs/api/)
3. Start testing all endpoints

[📖 Full Documentation](docs/api/README.md)

### Run the Application

**Using Docker Compose (Recommended)**

```bash
make compose-up
# App available at http://localhost:8000
```

**Manual Setup**

```bash
go mod download
make migrate-up
make run
```

**Hot Reload (Development)**

```bash
go install github.com/cosmtrek/air@latest
air
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `wallet_api` |
| `APP_VERSION` | Application version | `1.0.0` |
| `HTTP_PORT` | HTTP server port | `8000` |
| `LOG_LEVEL` | Logging level | `debug` |
| `PG_URL` | PostgreSQL connection string | - |
| `PG_POOL_MAX` | Database max connections | `2` |
| `JWT_SECRET` | JWT signing secret | - |
| `ACCESS_TOKEN_EXPIRY` | Access token expiry (minutes) | `15` |
| `REFRESH_TOKEN_EXPIRY` | Refresh token expiry (days) | `7` |

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v1/auth/register` | Register new user | No |
| POST | `/v1/auth/login` | Login user | No |
| POST | `/v1/auth/logout` | Logout user | Yes |
| POST | `/v1/auth/refresh` | Refresh access token | No |

### User Profile

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v1/users/profile` | Get user profile | Yes |
| PUT | `/v1/users/profile` | Update user profile | Yes |

### Accounts

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v1/accounts` | Create new account | Yes |
| GET | `/v1/accounts/:id` | Get account by ID | Yes |
| GET | `/v1/accounts` | Get all user accounts | Yes |
| POST | `/v1/accounts/:id/deposit` | Deposit to account | Yes |
| POST | `/v1/accounts/:id/withdraw` | Withdraw from account | Yes |
| POST | `/v1/accounts/:id/transfer` | Transfer to another account | Yes |
| GET | `/v1/accounts/:id/transactions` | Get account transactions | Yes |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/healthz` | Health check endpoint |

## Database Schema

### Entity Relationship Diagram

![Database ERD](docs/images/erd.png)

## Make Commands

```bash
# Build
make build              # Build application binary

# Run
make run                # Run application
make compose-up         # Start Docker services
make compose-down       # Stop Docker services
make compose-logs       # View Docker logs

# Database
make migrate-create     # Create new migration
make migrate-up         # Run migrations
make migrate-down       # Rollback last migration
make migrate-force      # Force version (usage: make migrate-force VERSION=001)

# Utilities
make nuke               # Remove all containers, volumes, and images
make clean              # Clean build artifacts
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

## Architecture Highlights

### Clean Architecture Diagram

![Clean Architecture](docs/images/cleanArchitect.webp)

### Clean Architecture Principles

- **Dependency Injection**: All dependencies injected through constructors
- **Interface-based Design**: Repository and UseCase defined as interfaces
- **Layer Separation**: Handler → UseCase → Repository → Entity
- **Encapsulation**: Private concrete types, public interfaces

### Module Structure

Each module follows this pattern:
```
module/
├── module.go         # Module initialization (DI)
├── router.go         # Route registration
├── handler/          # HTTP handlers
├── usecase/          # Business logic
├── repository/       # Data access
└── dto/              # Request/Response DTOs
    ├── request/      # Request DTOs
    └── response/     # Response DTOs + mappers
```

## Response Format

All API responses follow this format:

**Success Response**:
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response**:
```json
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Error message"
  }
}
```
