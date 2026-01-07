<p align="center">
  <img src="pulsar.svg" alt="Pulsar Logo" width="180" height="180">
</p>

<h1 align="center">Pulsar</h1>

<p align="center">
  <b>A modern, open-source incident management platform.</b><br>
  Alerting, on-call scheduling, escalations, and real-time notifications — an Opsgenie alternative built with Go and Svelte.
</p>

---

<p align="center">
  <a href="https://github.com/nmn3m/pulsar/releases"><img src="https://img.shields.io/badge/release-v0.1.0-blue?style=flat-square" alt="Release"></a>
  <a href="#license"><img src="https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square" alt="License"></a>
  <a href="#"><img src="https://img.shields.io/badge/build-passing-brightgreen?style=flat-square" alt="Build"></a>
  <a href="#"><img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
  <a href="#"><img src="https://img.shields.io/badge/Svelte-4-FF3E00?style=flat-square&logo=svelte&logoColor=white" alt="Svelte"></a>
  <a href="#"><img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL"></a>
</p>

<p align="center">
  <a href="#features">Features</a> &bull;
  <a href="#quick-start">Quick Start</a> &bull;
  <a href="#installation">Installation</a> &bull;
  <a href="#api-reference">API Reference</a> &bull;
  <a href="#development">Development</a> &bull;
  <a href="#architecture">Architecture</a>
</p>

---

## Features

### Core Capabilities

| Feature | Description | Status |
|---------|-------------|--------|
| **Alert Management** | Create, acknowledge, snooze, assign, and resolve alerts with full lifecycle tracking | Done |
| **Incident Management** | Track incidents with responders, timeline, notes, and linked alerts | Done |
| **On-Call Schedules** | Flexible scheduling with rotations, participants, and schedule overrides | Done |
| **Escalation Policies** | Multi-level escalation rules with configurable targets and delays | Done |
| **Team Management** | Organize users into teams with role-based member management | Done |
| **Notifications** | Channel configuration, user preferences, and delivery tracking | Done |
| **Webhooks** | Outgoing webhook endpoints and incoming webhook tokens for integrations | Done |
| **Real-time Updates** | WebSocket support for live dashboard updates | Done |

### Platform Features

- **JWT Authentication** — Secure access and refresh token authentication
- **API Key Authentication** — Scoped API keys for programmatic access with fine-grained permissions
- **Multi-Tenancy** — Full organization isolation with scoped data access
- **Role-Based Access Control** — Admin, member, and viewer permission levels
- **Dark/Light Theme** — Beautiful UI with theme switching support
- **Swagger Documentation** — Interactive API documentation
- **Background Workers** — Automated escalation and webhook delivery processing

---

## Quick Start

```bash
# Clone and enter the repository
git clone https://github.com/nmn3m/pulsar.git && cd pulsar

# Copy environment configuration
cp .env.example .env

# Start all services
make up

# Run database migrations
make migrate-up
```

**Access the application:**
| Service | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| Backend API | http://localhost:8081 |
| Swagger Docs | http://localhost:8081/swagger/index.html |
| PostgreSQL | localhost:5433 |

---

## Installation

### Prerequisites

- Docker and Docker Compose
- (Optional) Go 1.25+ and Node.js 20+ for local development

### Docker Compose (Recommended)

```bash
# Start all services in detached mode
make up

# View logs
make logs

# Stop services
make down
```

### Manual Setup

<details>
<summary><b>Backend Setup</b></summary>

```bash
cd backend

# Install Go dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://pulsar:pulsar_dev_password@localhost:5433/pulsar?sslmode=disable"
export JWT_SECRET="your-secret-key-min-32-characters"
export JWT_REFRESH_SECRET="your-refresh-secret-min-32-characters"
export SERVER_PORT="8080"

# Run with live reload
air

# Or run directly
go run cmd/api/main.go
```

</details>

<details>
<summary><b>Frontend Setup</b></summary>

```bash
cd frontend

# Install dependencies
npm install

# Set environment variables
export VITE_API_URL="http://localhost:8081"

# Run development server
npm run dev
```

</details>

---

## Tech Stack

| Layer | Technologies |
|-------|--------------|
| **Backend** | Go 1.25, Gin, PostgreSQL 16, JWT, Zap Logger |
| **Frontend** | SvelteKit, TypeScript, Tailwind CSS, Vite |
| **Infrastructure** | Docker, Docker Compose |
| **Documentation** | Swagger / OpenAPI 3.0 |

---

## API Reference

Base URL: `http://localhost:8081/api/v1`

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth/register` | Register a new user and organization |
| `POST` | `/auth/login` | Authenticate and receive tokens |
| `POST` | `/auth/refresh` | Refresh access token |
| `POST` | `/auth/logout` | Invalidate refresh token |
| `GET` | `/auth/me` | Get current user info |

### API Keys

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api-keys` | List your API keys |
| `POST` | `/api-keys` | Create a new API key |
| `GET` | `/api-keys/scopes` | List available scopes |
| `GET` | `/api-keys/all` | List all org API keys (admin) |
| `GET` | `/api-keys/:id` | Get API key details |
| `PATCH` | `/api-keys/:id` | Update an API key |
| `DELETE` | `/api-keys/:id` | Delete an API key |
| `POST` | `/api-keys/:id/revoke` | Revoke an API key |

**Available Scopes:**
- `alerts:read`, `alerts:write` — Alert management
- `incidents:read`, `incidents:write` — Incident management
- `teams:read`, `teams:write` — Team management
- `schedules:read`, `schedules:write` — Schedule management
- `webhooks:read`, `webhooks:write` — Webhook management
- `notifications:read`, `notifications:write` — Notification management
- `users:read` — User listing
- `*` — Full access

**Usage:** Include the API key in the `X-API-Key` header or as `Authorization: ApiKey <key>`.

### Alerts

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/alerts` | List alerts with filters |
| `POST` | `/alerts` | Create a new alert |
| `GET` | `/alerts/:id` | Get alert details |
| `PATCH` | `/alerts/:id` | Update alert |
| `DELETE` | `/alerts/:id` | Delete alert |
| `POST` | `/alerts/:id/acknowledge` | Acknowledge alert |
| `POST` | `/alerts/:id/close` | Close alert with reason |
| `POST` | `/alerts/:id/snooze` | Snooze alert until time |
| `POST` | `/alerts/:id/assign` | Assign to user or team |

### Incidents

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/incidents` | List incidents |
| `POST` | `/incidents` | Create incident |
| `GET` | `/incidents/:id` | Get incident with details |
| `PATCH` | `/incidents/:id` | Update incident |
| `POST` | `/incidents/:id/responders` | Add responder |
| `GET` | `/incidents/:id/timeline` | Get incident timeline |
| `POST` | `/incidents/:id/notes` | Add note to timeline |
| `POST` | `/incidents/:id/alerts` | Link alert to incident |

### Teams, Schedules, Escalations

<details>
<summary>View all endpoints</summary>

**Teams**
- `GET/POST` `/teams` — List/create teams
- `GET/PATCH/DELETE` `/teams/:id` — Manage team
- `GET/POST` `/teams/:id/members` — Manage members

**Schedules**
- `GET/POST` `/schedules` — List/create schedules
- `GET/PATCH/DELETE` `/schedules/:id` — Manage schedule
- `GET` `/schedules/:id/oncall` — Get current on-call
- `*/schedules/:id/rotations/*` — Manage rotations
- `*/schedules/:id/overrides/*` — Manage overrides

**Escalation Policies**
- `GET/POST` `/escalation-policies` — List/create policies
- `GET/PATCH/DELETE` `/escalation-policies/:id` — Manage policy
- `*/escalation-policies/:id/rules/*` — Manage rules and targets

**Notifications**
- `*/notifications/channels/*` — Notification channels
- `*/notifications/preferences/*` — User preferences
- `POST` `/notifications/send` — Send notification
- `GET` `/notifications/logs/*` — Delivery logs

**Webhooks**
- `*/webhooks/endpoints/*` — Outgoing webhooks
- `*/webhooks/incoming/*` — Incoming webhook tokens
- `GET` `/webhooks/deliveries` — Delivery history
- `POST` `/webhook/:token` — Receive incoming webhook (public)

**Metrics**
- `GET` `/metrics/dashboard` — Dashboard metrics (alerts, incidents, notifications, trends)
- `GET` `/metrics/alerts` — Alert metrics (by status, priority, source)
- `GET` `/metrics/alerts/trend` — Alert time-series data (hourly/daily/weekly)
- `GET` `/metrics/incidents` — Incident metrics (by status, severity)
- `GET` `/metrics/notifications` — Notification metrics (by status, channel)
- `GET` `/metrics/teams` — Team performance metrics

</details>

---

## Development

### Available Commands

```bash
make up                 # Start all services
make down               # Stop all services
make logs               # View container logs
make build              # Rebuild Docker images

make migrate-up         # Apply database migrations
make migrate-down       # Rollback last migration
make migrate-create NAME=xyz  # Create new migration

make test               # Run unit tests
make test-integration   # Run integration tests (150+ tests)
make test-coverage      # Generate coverage report

make clean              # Remove containers and volumes
```

### Integration Testing

The project includes a comprehensive integration test suite covering all API endpoints:

```bash
# Run all integration tests
make test-integration

# Run specific test file
cd backend && go test -v ./tests/integration/... -run TestAlerts
```

**Test Coverage:**

| Domain | Tests | Endpoints Covered |
|--------|-------|-------------------|
| Auth | 14 | Register, Login, Refresh, Me, Logout |
| Alerts | 23 | CRUD, Acknowledge, Close, Snooze, Assign |
| Teams | 16 | CRUD, Members management |
| Schedules | 19 | CRUD, Rotations, Overrides, On-call |
| Escalations | 17 | CRUD, Rules, Targets |
| Incidents | 22 | CRUD, Responders, Timeline, Notes, Alerts |
| Notifications | 19 | Channels, Preferences, Send, Logs |
| Webhooks | 18 | Endpoints, Incoming tokens, Deliveries |
| Metrics | 17 | Dashboard, Alerts, Incidents, Notifications, Teams |
| Users | 2 | List users |

The tests use an isolated PostgreSQL database (`docker-compose.test.yml`) and clean the database between test runs for isolation.

### Backend Development

```bash
cd backend

# Run with hot reload
air

# Run tests
go test -v ./...

# Generate Swagger docs
swag init -g cmd/api/main.go
```

### Frontend Development

```bash
cd frontend

# Development server with HMR
npm run dev

# Type checking
npm run check

# Production build
npm run build
```

---

## Architecture

```
pulsar/
├── backend/
│   ├── cmd/api/              # Application entrypoint
│   ├── internal/
│   │   ├── config/           # Configuration loading
│   │   ├── domain/           # Business entities
│   │   ├── handler/rest/     # HTTP handlers
│   │   ├── middleware/       # Auth, CORS, logging
│   │   ├── repository/       # Data access layer
│   │   └── service/          # Business logic
│   ├── migrations/           # SQL migrations
│   └── docs/                 # Swagger specs
│
├── frontend/
│   └── src/
│       ├── lib/
│       │   ├── api/          # API client
│       │   ├── components/   # UI components
│       │   └── stores/       # State management
│       └── routes/           # SvelteKit pages
│           ├── (auth)/       # Login, register
│           └── (app)/        # Dashboard, alerts, etc.
│
├── docker-compose.yml        # Development environment
└── Makefile                  # Build automation
```

### Design Principles

- **Clean Architecture** — Separation of concerns with domain, service, and handler layers
- **Repository Pattern** — Abstracted data access for testability
- **Multi-Tenancy** — Organization-scoped data isolation
- **Event-Driven** — WebSocket notifications and background workers

---

## Security

- JWT authentication with short-lived access tokens (15 min) and refresh tokens (7 days)
- Password hashing with bcrypt (cost factor 10)
- Role-based access control at organization level
- SQL injection prevention via parameterized queries
- XSS protection through proper encoding
- CORS configuration for allowed origins

---

## Roadmap

### ✅ Completed
- [x] User authentication and multi-tenancy
- [x] Alert management with full lifecycle
- [x] Team management and RBAC
- [x] On-call schedules with rotations
- [x] Escalation policies with background processing
- [x] Incident management with timeline
- [x] Notification system with preferences
- [x] Webhooks (outgoing and incoming)
- [x] Real-time WebSocket updates
- [x] Dark/Light theme support
- [x] Email notification delivery
- [x] Slack integration (webhook-based)
- [x] Microsoft Teams integration (webhook-based)
- [x] Integration test suite (150+ tests)

### 🚧 In Progress
- [x] API key authentication for programmatic access
- [x] Metrics and reporting dashboard (API complete)
- [ ] Mobile push notifications (Firebase/APNS)

### 📋 Planned
- [ ] Slack App (OAuth, interactive messages)
- [ ] Microsoft Teams App (Bot framework)
- [ ] PagerDuty/Opsgenie import tools
- [ ] Terraform provider
- [ ] Prometheus/Grafana integration
- [ ] Audit logging

---

## Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the **Apache License 2.0** — see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <sub>Built with Go and Svelte</sub>
</p>
