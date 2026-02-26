# Pulsar Backend — Developer Guide

This guide explains the architecture, code organization, and development workflow for the Pulsar backend. Read this before contributing.

## Architecture Overview

The backend follows **Clean Architecture** with four layers. Each layer has a strict dependency rule — inner layers never import outer layers.

```
                    ┌──────────────────────┐
                    │       domain/        │  Entities, value objects, errors
                    │   (no dependencies)  │  No struct tags, no framework imports
                    └──────────┬───────────┘
                               │
                    ┌──────────▼───────────┐
                    │      usecase/        │  Business logic, repository interfaces
                    │  (imports domain/)   │  Request/response DTOs live here
                    └──────────┬───────────┘
                               │
              ┌────────────────┼────────────────┐
              │                                 │
   ┌──────────▼────────────┐         ┌──────────▼────────────┐
   │  repository/postgres/ │         │   delivery/rest/      │
   │ (implements usecase)  │         │ (calls usecase layer) │
   │  repository interfaces│         │  handlers, middleware │
   └───────────────────────┘         └───────────────────────┘
```

**Dependency rule**: `domain/ ← usecase/ ← repository/postgres/` and `delivery/rest/ → usecase/ → domain/`

## Project Structure

```
backend/
├── cmd/
│   ├── api/main.go              # Application entry point, wiring
│   └── seed/main.go             # Demo data seeder
├── internal/
│   ├── config/                  # Environment variable loading
│   ├── domain/                  # Layer 1: Pure entities
│   ├── usecase/                 # Layer 2: Business logic
│   │   ├── repository/          # Repository interface definitions
│   │   └── providers/           # Notification providers (email, slack, teams, webhook)
│   ├── repository/              # Layer 3: Data access
│   │   └── postgres/            # PostgreSQL implementations
│   ├── delivery/                # Layer 4: HTTP transport
│   │   └── rest/
│   │       ├── handler/         # HTTP handlers (one per domain area)
│   │       └── middleware/      # Auth, API key, CORS, logging, OpenTelemetry
│   └── pkg/
│       ├── logger/              # Zap logger setup
│       └── telemetry/           # OpenTelemetry initialization
├── migrations/                  # SQL migration files (up/down pairs)
├── tests/integration/           # Integration tests
│   └── testutils/               # Test server, fixtures, HTTP client
└── docs/                        # Generated Swagger docs
```

## The Four Layers

### Layer 1: `domain/` — Entities and Value Objects

Pure Go structs with no struct tags (`json`, `db`, `binding` are all absent). This layer defines:

- **Entities**: `Alert`, `Incident`, `Team`, `User`, `Schedule`, `EscalationPolicy`, etc.
- **Value objects**: Priority levels, statuses, roles, channel types
- **Errors**: `ErrNotFound`, `ErrUnauthorized`, `ErrValidation`, etc.
- **Validation methods**: Business rule checks on entities
- **Constants**: Enums like `AlertStatusOpen`, `PriorityP1`, `RotationTypeWeekly`

This layer imports only the standard library and `uuid`. Nothing else.

**Key files to start with:**
- `domain/alert.go` — Core alert entity and related types
- `domain/incident.go` — Incident management types
- `domain/errors.go` — Shared error definitions

### Layer 2: `usecase/` — Business Logic

Each domain area has a `*_usecase.go` file containing:

- A **usecase struct** (e.g., `AlertUsecase`) with repository dependencies injected via constructor
- **Request/response DTOs** with `json` and `binding` tags, defined locally in the same file
- **Business methods** that orchestrate repository calls, validation, and side effects

```go
// Example: usecase/alert_usecase.go
type CreateAlertRequest struct {
    Source   string `json:"source" binding:"required"`
    Priority string `json:"priority" binding:"required"`
    Message  string `json:"message" binding:"required"`
    // ...
}

type AlertUsecase struct {
    alertRepo      repository.AlertRepository
    alertNotifier  *AlertNotifier
    wsUsecase      *WebSocketUsecase
    webhookUsecase *WebhookUsecase
}

func NewAlertUsecase(...) *AlertUsecase { ... }
func (s *AlertUsecase) CreateAlert(ctx context.Context, orgID uuid.UUID, req *CreateAlertRequest) (*domain.Alert, error) { ... }
```

**`usecase/repository/`** — Interfaces that define what the usecase layer needs from persistence. One file per domain area. The postgres layer implements these.

```go
// usecase/repository/alert_repository.go
type AlertRepository interface {
    Create(ctx context.Context, alert *domain.Alert) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error)
    // ...
}
```

**Key files:**
- `usecase/auth_usecase.go` — Registration, login, JWT generation/validation
- `usecase/alert_usecase.go` — Alert CRUD, acknowledgment, assignment
- `usecase/alert_notifier.go` — Notification routing when alerts fire
- `usecase/notification_usecase.go` — Multi-channel notification dispatch

### Layer 3: `repository/postgres/` — Database Access

Each file implements one interface from `usecase/repository/`. Uses `sqlx` for queries with manual `rows.Scan()` (no ORM).

- `checks.go` — Compile-time interface satisfaction checks (`var _ repository.AlertRepository = (*AlertRepository)(nil)`)
- `db.go` — Database connection setup and migration runner
- One `*_repo.go` per domain area

### Layer 4: `delivery/rest/` — HTTP Transport

**`handler/`** — One handler per domain area. Each handler:
1. Binds the HTTP request to a usecase request DTO
2. Extracts auth context (user ID, org ID) from middleware
3. Calls the appropriate usecase method
4. Returns the JSON response

**`middleware/`** — Five middleware files:
- `auth.go` — JWT token validation, sets user claims in context
- `apikey.go` — API key validation, combined JWT+API key auth
- `cors.go` — CORS headers
- `logging.go` — Request/response logging with Zap
- `otel.go` — OpenTelemetry trace/metrics instrumentation

## How to Read the Code

Start from the outside in, following a single request:

1. **Pick an endpoint** — Open `cmd/api/main.go` and find a route, e.g., `alerts.POST("", alertHandler.Create)`
2. **Read the handler** — Open `delivery/rest/handler/alert_handler.go`, find the `Create` method. See how it binds JSON, extracts auth, calls the usecase.
3. **Read the usecase** — Open `usecase/alert_usecase.go`, find `CreateAlert`. See the business logic, validation, and repository calls.
4. **Read the repository** — Open `repository/postgres/alert_repo.go`, find `Create`. See the SQL query.
5. **Read the domain** — Open `domain/alert.go` to understand the entity fields and constants.

For understanding how the layers wire together, read `cmd/api/main.go` from the repository initialization section downward.

## Development Setup

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- Docker & Docker Compose

### Running with Docker Compose

```bash
# From the project root (not backend/)
docker-compose up -d

# Seed demo data
docker-compose exec backend go run cmd/seed/main.go

# View logs
docker-compose logs -f backend
```

This starts PostgreSQL (port 5433), Mailpit for email testing (port 8025), the backend (port 8081), and the frontend.

### Running Locally

```bash
cd backend

# Install dependencies
go mod download

# Set required environment variables
export DATABASE_URL="postgres://pulsar:pulsar_dev_password@localhost:5433/pulsar?sslmode=disable"
export JWT_SECRET="dev_jwt_secret_change_in_production_min_32_chars"
export JWT_REFRESH_SECRET="dev_refresh_secret_change_in_production_min_32_chars"

# Run database migrations (handled automatically on startup via db.go)
go run cmd/api/main.go

# Or use Air for live reload
air
```

### Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `JWT_SECRET` | Yes | — | JWT signing key (min 32 chars) |
| `JWT_REFRESH_SECRET` | Yes | — | Refresh token signing key (min 32 chars) |
| `SERVER_PORT` | No | `8080` | HTTP server port |
| `ENV` | No | `development` | Environment (`development`, `production`) |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000` | Comma-separated origins |
| `SMTP_ENABLED` | No | `false` | Enable SMTP email |
| `SMTP_HOST` | No | `localhost` | SMTP server host |
| `SMTP_PORT` | No | `587` | SMTP server port |
| `EMAIL_ENABLED` | No | `false` | Enable email provider |
| `EMAIL_PROVIDER` | No | `smtp` | `smtp` or `resend` |
| `RESEND_API_KEY` | No | — | Resend API key (production) |
| `OTEL_ENABLED` | No | `false` | Enable OpenTelemetry |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | No | `localhost:4317` | OTLP collector endpoint |

## Testing

### Integration Tests

Integration tests run against a real PostgreSQL database:

```bash
# Start test database
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
cd backend
DATABASE_URL="postgres://..." JWT_SECRET="..." JWT_REFRESH_SECRET="..." \
  go test -v ./tests/integration/...
```

Tests use `testutils/` helpers:
- `NewTestServer()` — Spins up a full server with all layers wired
- `NewTestFixtures()` — Factory methods for creating test users, teams, alerts, etc.
- Each test calls `cleanDatabase(t)` to truncate tables before running

### Adding a New Test

```go
func TestAlerts_Create_Success(t *testing.T) {
    cleanDatabase(t)
    ctx := context.Background()
    client := newTestClient(t)

    user, _ := testFixtures.CreateUniqueUser(ctx)
    client.SetAuthToken(user.AccessToken)

    body := map[string]interface{}{
        "source":   "test",
        "priority": "P3",
        "message":  "Test alert",
    }

    resp := client.Post("/api/v1/alerts", body)
    client.AssertStatus(resp, http.StatusCreated)
}
```

## Adding a New Feature

Follow this order when adding a new domain area (e.g., "status pages"):

1. **Domain entity** — Add `domain/statuspage.go` with the struct (no tags) and constants
2. **Repository interface** — Add `usecase/repository/statuspage_repository.go`
3. **Usecase** — Add `usecase/statuspage_usecase.go` with request DTOs, business logic, and constructor
4. **Postgres repo** — Add `repository/postgres/statuspage_repo.go` implementing the interface
5. **Compile-time check** — Add `var _ repository.StatusPageRepository = (*StatusPageRepository)(nil)` to `repository/postgres/checks.go`
6. **Handler** — Add `delivery/rest/handler/statuspage_handler.go`
7. **Wire it up** — Update `cmd/api/main.go`: instantiate repo → usecase → handler, add routes
8. **Migration** — Add a new SQL migration file in `migrations/`

## API Documentation

Swagger docs are auto-generated and available at `http://localhost:8080/swagger/index.html` when the server is running.

To regenerate after adding/modifying handler annotations:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs
```

## Key Design Decisions

- **No ORM** — Raw SQL with `sqlx`. Queries are explicit and optimizable.
- **Interfaces at the consumer** — Repository interfaces live in `usecase/repository/`, not next to the implementations. This follows the Dependency Inversion Principle.
- **DTOs in usecase, not domain** — Request/response types with `json`/`binding` tags live in the usecase files. Domain entities stay pure.
- **Narrow config injection** — Usecases receive only the config they need (e.g., `AuthConfig` with 4 fields) instead of the entire `*config.Config`.
- **WebSocket hub** — Real-time updates for alerts and incidents via a centralized WebSocket hub in `WebSocketUsecase`.
- **Background workers** — Escalation processing and webhook delivery run as background goroutines started in `main.go`.
