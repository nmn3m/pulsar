# Pulsar Backend

Go backend API for Pulsar incident management platform.

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Logging**: Uber Zap

## Project Structure

```
backend/
├── cmd/
│   └── api/              # Main application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── domain/           # Domain models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   └── pkg/              # Internal packages
└── migrations/           # Database migrations
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 16 or higher
- Docker & Docker Compose (optional)

### Running with Docker Compose

```bash
# From project root
docker-compose up -d

# Check logs
docker-compose logs -f backend
```

### Running Locally

```bash
# Install dependencies
go mod download

# Set up environment variables
cp ../.env.example ../.env
# Edit .env with your configuration

# Run migrations
make migrate-up

# Run the server
go run cmd/api/main.go
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user (protected)

### Health Check

- `GET /health` - Health check endpoint

## Environment Variables

See `.env.example` for required environment variables.

## Database Migrations

```bash
# Create new migration
make migrate-create NAME=add_something

# Apply migrations
make migrate-up

# Rollback last migration
make migrate-down
```

## Development

```bash
# Run tests
go test -v ./...

# Run with live reload (using air)
air
```
