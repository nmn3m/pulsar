<div align="center">

<img src="pulsar.svg" alt="Pulsar" width="100" height="100">

# Pulsar

**Modern Incident Management Platform**

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Svelte](https://img.shields.io/badge/Svelte-4-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)](https://svelte.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://postgresql.org/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue?style=for-the-badge)](LICENSE)

[Features](#features) • [Screenshots](#screenshots) • [Getting Started](#getting-started) • [Documentation](#documentation)

---

*Open-source incident management platform with alerting, on-call scheduling, and escalation policies. Built for teams who want full control over their workflow.*

</div>

<br>

## Features

<table>
<tr>
<td width="50%">

### Alerting & Response
- **Alert Management** — Full lifecycle tracking with acknowledge, snooze, escalate, and resolve actions
- **Routing Rules** — Automatically route alerts based on priority, source, tags, or content
- **Escalation Policies** — Multi-level escalation chains with configurable delays

</td>
<td width="50%">

### On-Call & Scheduling
- **On-Call Schedules** — Flexible rotations with daily, weekly, or custom patterns
- **Schedule Overrides** — Easy shift swaps without modifying the rotation
- **Who's On-Call** — Instant visibility into current responders

</td>
</tr>
<tr>
<td width="50%">

### Incident Management
- **Incident Tracking** — Coordinate response with responders and timeline
- **Alert Correlation** — Link related alerts to incidents
- **Post-Incident Notes** — Document findings and actions taken

</td>
<td width="50%">

### Integrations & API
- **Notification Channels** — Email, Slack, Microsoft Teams, webhooks
- **Incoming Webhooks** — Accept alerts from any monitoring tool
- **REST API** — Full API with scoped API key authentication

</td>
</tr>
</table>

<br>

## Screenshots

<div align="center">
<table>
<tr>
<td><img src="docs/screenshots/dashboard.png" alt="Dashboard" width="400"></td>
<td><img src="docs/screenshots/alerts.png" alt="Alerts" width="400"></td>
</tr>
<tr>
<td align="center"><b>Dashboard</b><br><sub>Real-time metrics and alert overview</sub></td>
<td align="center"><b>Alerts</b><br><sub>Manage alerts with priority filtering</sub></td>
</tr>
<tr>
<td><img src="docs/screenshots/incidents.png" alt="Incidents" width="400"></td>
<td><img src="docs/screenshots/teams.png" alt="Teams" width="400"></td>
</tr>
<tr>
<td align="center"><b>Incidents</b><br><sub>Track and coordinate incident response</sub></td>
<td align="center"><b>Teams</b><br><sub>Organize responders into teams</sub></td>
</tr>
<tr>
<td><img src="docs/screenshots/schedules.png" alt="Schedules" width="400"></td>
<td><img src="docs/screenshots/escalations.png" alt="Escalations" width="400"></td>
</tr>
<tr>
<td align="center"><b>On-Call Schedules</b><br><sub>Flexible rotation management</sub></td>
<td align="center"><b>Escalation Policies</b><br><sub>Multi-level escalation rules</sub></td>
</tr>
<tr>
<td><img src="docs/screenshots/notifications.png" alt="Notifications" width="400"></td>
<td><img src="docs/screenshots/webhooks.png" alt="Webhooks" width="400"></td>
</tr>
<tr>
<td align="center"><b>Notifications</b><br><sub>Email, Slack, Teams channels</sub></td>
<td align="center"><b>Webhooks</b><br><sub>Outgoing integrations</sub></td>
</tr>
</table>
</div>

<br>

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/nmn3m/pulsar.git
cd pulsar

# Start all services
make up

# Apply database migrations
make migrate-up
```

### Access

| Service | URL |
|:--------|:----|
| Web UI | [http://localhost:5173](http://localhost:5173) |
| API | [http://localhost:8081/api/v1](http://localhost:8081/api/v1) |
| API Docs | [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html) |

### Demo Data

Optionally seed the database with sample data for testing:

```bash
cd backend
DATABASE_URL="postgres://pulsar:pulsar_dev_password@localhost:5433/pulsar?sslmode=disable" \
JWT_SECRET="dev_jwt_secret_change_in_production_min_32_chars" \
JWT_REFRESH_SECRET="dev_refresh_secret_change_in_production_min_32_chars" \
go run ./cmd/seed/main.go
```

> **Demo credentials:** `admin@acme-corp.com` / `DemoPass123!`

<br>

## Documentation

### API Authentication

Pulsar supports two authentication methods:

```bash
# JWT Token Authentication
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password"}'

# Then use the token
curl http://localhost:8081/api/v1/alerts \
  -H "Authorization: Bearer <access_token>"
```

```bash
# API Key Authentication
curl http://localhost:8081/api/v1/alerts \
  -H "X-API-Key: <your_api_key>"
```

### API Endpoints

| Resource | Description |
|:---------|:------------|
| `/auth/*` | Authentication (login, register, refresh) |
| `/alerts/*` | Alert CRUD, acknowledge, close, snooze, assign |
| `/incidents/*` | Incident management with responders and timeline |
| `/teams/*` | Team and member management |
| `/schedules/*` | On-call schedules, rotations, overrides |
| `/escalation-policies/*` | Escalation rules and targets |
| `/webhooks/*` | Incoming and outgoing webhook configuration |
| `/metrics/*` | Dashboard and reporting metrics |

> Full interactive documentation available at `/swagger/index.html`

<br>

## Architecture

```
                              ┌─────────────────────┐
                              │    Notification     │
                              │      Channels       │
                              │  (Email, Slack,     │
                              │   Teams, Webhook)   │
                              └──────────▲──────────┘
                                         │
┌─────────────────┐           ┌──────────┴──────────┐             ┌─────────────────┐
│   Monitoring    │           │                     │             │    Frontend     │
│   Systems       │─────────▶│   Pulsar Backend    │◀────────── │   (SvelteKit)   │
│  (Prometheus,   │  Webhook  │      (Go/Gin)       │    API      │                 │
│ Datadog, etc.)  │           │                     │             │                 │
└─────────────────┘           └──────────┬──────────┘             └─────────────────┘
                                         │
                              ┌──────────▼──────────┐
                              │     PostgreSQL      │
                              │      Database       │
                              └─────────────────────┘
```

### Tech Stack

| Component | Technology |
|:----------|:-----------|
| Backend | Go 1.25, Gin, sqlx |
| Frontend | SvelteKit, TypeScript, Tailwind CSS |
| Database | PostgreSQL 16 |
| Auth | JWT, bcrypt |
| Real-time | WebSocket |
| Deployment | Docker, Docker Compose |

<br>

## Development

```bash
make up                 # Start services
make down               # Stop services
make logs               # View logs
make migrate-up         # Run migrations
make test-integration   # Run test suite
```

<details>
<summary><b>Local Development Setup</b></summary>

**Backend:**
```bash
cd backend
go mod download
export DATABASE_URL="postgres://pulsar:pulsar_dev_password@localhost:5433/pulsar?sslmode=disable"
export JWT_SECRET="dev_jwt_secret_change_in_production_min_32_chars"
air  # Hot reload development server
```

**Frontend:**
```bash
cd frontend
npm install
export VITE_API_URL="http://localhost:8081"
npm run dev
```

</details>

<br>

## License

This project is licensed under the [Apache License 2.0](LICENSE).

---

<div align="center">
<sub>Built with Go and Svelte</sub>
</div>
