# Pulsar - Opsgenie Replacement Implementation Plan

## Overview
Building a production-ready incident management platform with:
- **Backend**: Go (Golang)
- **Frontend**: Svelte + TypeScript
- **Database**: PostgreSQL
- **Features**: Alert Management, On-Call Schedules, Incident Management, Escalation Policies
- **Integrations**: Webhooks, Email, Slack/Teams, REST API

## Tech Stack

### Backend
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with sqlx
- **Auth**: JWT with refresh tokens
- **WebSocket**: gorilla/websocket
- **Migrations**: golang-migrate
- **Logging**: uber/zap

### Frontend
- **Framework**: SvelteKit + TypeScript
- **Styling**: Tailwind CSS
- **Build**: Vite
- **API Client**: Fetch with custom wrapper

### Infrastructure
- **Development**: Docker Compose
- **Containerization**: Docker multi-stage builds
- **Real-time**: WebSockets for live updates

## Architecture

```
Frontend (Svelte) → API Gateway (Go/Gin) → Services Layer → PostgreSQL
                  ↕                         ↓
              WebSocket Hub ← Background Workers (Alert Router, Escalation, Notifications)
                                            ↓
                              External Integrations (Email, Slack, Teams)
```

### Key Architectural Patterns
- **Clean Architecture**: Domain → Repository → Service → Handler layers
- **Multi-tenancy**: Organization-based data isolation
- **RBAC**: Owner, Admin, Member, Viewer roles
- **Real-time**: WebSocket hub for live updates

## Project Structure

```
/home/nour/workspace/github/pulsar/
├── backend/
│   ├── cmd/
│   │   ├── api/main.go              # API server entry point
│   │   ├── worker/main.go           # Background workers
│   │   └── migrate/main.go          # Migration tool
│   ├── internal/
│   │   ├── config/config.go         # Configuration management
│   │   ├── domain/                  # Domain models (alert.go, incident.go, etc.)
│   │   ├── repository/              # Data access layer
│   │   │   ├── interfaces.go
│   │   │   └── postgres/            # PostgreSQL implementations
│   │   ├── service/                 # Business logic
│   │   ├── handler/                 # HTTP handlers
│   │   │   ├── rest/
│   │   │   └── websocket/
│   │   ├── middleware/              # Auth, CORS, logging, rate limiting
│   │   ├── worker/                  # Background job processors
│   │   └── integration/             # External integrations (email, slack, teams)
│   ├── migrations/                  # SQL migrations
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/                 # API client
│   │   │   ├── stores/              # Svelte stores (state)
│   │   │   ├── components/          # Reusable components
│   │   │   ├── types/               # TypeScript types
│   │   │   └── utils/
│   │   └── routes/
│   │       ├── (auth)/              # Login, register
│   │       └── (app)/               # Main app (dashboard, alerts, incidents, etc.)
│   ├── package.json
│   ├── svelte.config.js
│   └── Dockerfile
└── docker-compose.yml               # Local development environment
```

## Database Schema

### Core Tables
1. **organizations** - Multi-tenant isolation
2. **users** - User accounts
3. **organization_users** - User-org mapping with roles
4. **teams** - Teams within organizations
5. **alerts** - Alert management with assignment, escalation
6. **incidents** - Incident tracking
7. **incident_alerts** - Link alerts to incidents
8. **incident_timeline** - Incident activity log
9. **schedules** - On-call schedules
10. **schedule_rotations** - Rotation members
11. **schedule_overrides** - Temporary schedule changes
12. **escalation_policies** - Escalation rules
13. **escalation_rules** - Steps in escalation
14. **notifications** - Notification log
15. **integrations** - External integrations config
16. **api_keys** - API key management
17. **audit_logs** - Audit trail

## API Design

### Authentication
- JWT-based with refresh tokens
- Access token: 15 minutes
- Refresh token: 7 days
- API key support for programmatic access

### REST Endpoints Structure
```
/api/v1/auth/*              # Authentication
/api/v1/alerts/*            # Alert management
/api/v1/incidents/*         # Incident management
/api/v1/schedules/*         # On-call schedules
/api/v1/escalation-policies/* # Escalation policies
/api/v1/teams/*             # Team management
/api/v1/users/*             # User management
/api/v1/integrations/*      # Integration management
/api/v1/webhooks/*          # Webhook endpoints
/api/v1/ws                  # WebSocket connection
```

## Implementation Phases

### Phase 1: Foundation & Authentication ⭐ START HERE
**Goal**: Basic app with user authentication working

**Critical Files to Create**:
1. `/home/nour/workspace/github/pulsar/docker-compose.yml` - Development environment
2. `/home/nour/workspace/github/pulsar/backend/go.mod` - Go dependencies
3. `/home/nour/workspace/github/pulsar/backend/cmd/api/main.go` - API server entry point
4. `/home/nour/workspace/github/pulsar/backend/internal/config/config.go` - Configuration
5. `/home/nour/workspace/github/pulsar/backend/migrations/001_initial_schema.up.sql` - Database schema
6. `/home/nour/workspace/github/pulsar/backend/internal/domain/user.go` - User model
7. `/home/nour/workspace/github/pulsar/backend/internal/domain/organization.go` - Organization model
8. `/home/nour/workspace/github/pulsar/backend/internal/middleware/auth.go` - JWT middleware
9. `/home/nour/workspace/github/pulsar/backend/internal/service/auth_service.go` - Auth logic
10. `/home/nour/workspace/github/pulsar/backend/internal/handler/rest/auth_handler.go` - Auth endpoints
11. `/home/nour/workspace/github/pulsar/frontend/package.json` - Frontend dependencies
12. `/home/nour/workspace/github/pulsar/frontend/svelte.config.js` - Svelte config
13. `/home/nour/workspace/github/pulsar/frontend/src/lib/api/client.ts` - API client
14. `/home/nour/workspace/github/pulsar/frontend/src/lib/stores/auth.ts` - Auth state
15. `/home/nour/workspace/github/pulsar/frontend/src/routes/(auth)/login/+page.svelte` - Login page

**Deliverable**: Users can register, login, and see authenticated dashboard

---

### Phase 2: Alert Management
**Goal**: Create and manage alerts

**Key Files**:
- `backend/internal/domain/alert.go`
- `backend/internal/service/alert_service.go`
- `backend/internal/handler/rest/alert_handler.go`
- `backend/migrations/002_alerts.up.sql`
- `frontend/src/lib/components/alerts/AlertList.svelte`
- `frontend/src/routes/(app)/alerts/+page.svelte`

**Deliverable**: Users can create, view, filter, acknowledge, and close alerts

---

### Phase 3: Team Management
**Goal**: Multi-user collaboration with teams

**Key Files**:
- `backend/internal/domain/team.go`
- `backend/internal/service/team_service.go`
- `backend/migrations/003_teams.up.sql`
- `frontend/src/routes/(app)/teams/+page.svelte`

**Deliverable**: Organizations can create teams and assign alerts to teams

---

### Phase 4: On-Call Schedules
**Goal**: Schedule on-call rotations

**Key Files**:
- `backend/internal/domain/schedule.go`
- `backend/internal/service/schedule_service.go`
- `backend/internal/worker/schedule_processor.go`
- `backend/migrations/004_schedules.up.sql`
- `frontend/src/lib/components/schedules/ScheduleCalendar.svelte`

**Deliverable**: Teams can set up on-call rotations and see who's on-call

---

### Phase 5: Escalation Policies
**Goal**: Auto-escalate unacknowledged alerts

**Key Files**:
- `backend/internal/domain/escalation.go`
- `backend/internal/service/escalation_service.go`
- `backend/internal/worker/escalation_processor.go`
- `backend/migrations/005_escalations.up.sql`
- `frontend/src/routes/(app)/escalations/+page.svelte`

**Deliverable**: Alerts automatically escalate based on configured policies

---

### Phase 6: Notifications
**Goal**: Multi-channel notifications (Email, Slack, Teams)

**Key Files**:
- `backend/internal/integration/email/smtp_client.go`
- `backend/internal/integration/slack/slack_client.go`
- `backend/internal/integration/teams/teams_client.go`
- `backend/internal/worker/notification_dispatcher.go`
- `frontend/src/routes/(app)/integrations/+page.svelte`

**Deliverable**: Users receive notifications via email, Slack, Teams

---

### Phase 7: Incident Management
**Goal**: Track and manage incidents

**Key Files**:
- `backend/internal/domain/incident.go`
- `backend/internal/service/incident_service.go`
- `backend/migrations/007_incidents.up.sql`
- `frontend/src/routes/(app)/incidents/+page.svelte`
- `frontend/src/lib/components/incidents/IncidentTimeline.svelte`

**Deliverable**: Users can create incidents, link alerts, track timeline

---

### Phase 8: Real-time Updates
**Goal**: Live updates via WebSocket

**Key Files**:
- `backend/internal/handler/websocket/ws_handler.go`
- `backend/internal/handler/websocket/hub.go`
- `frontend/src/lib/api/websocket.ts`
- `frontend/src/lib/stores/realtime.ts`

**Deliverable**: Real-time updates for alerts, incidents, on-call changes

---

### Phase 9: Webhooks & Integrations
**Goal**: External system integration

**Key Files**:
- `backend/internal/handler/rest/webhook_handler.go`
- `backend/internal/integration/webhook/webhook_client.go`
- `backend/migrations/009_integrations.up.sql`

**Deliverable**: External systems can send alerts via webhooks

---

### Phase 10: API Keys & Production Polish
**Goal**: Production readiness

**Key Files**:
- `backend/internal/middleware/rate_limit.go`
- `backend/migrations/010_api_keys.up.sql`
- `frontend/src/routes/(app)/settings/api-keys/+page.svelte`
- Production Dockerfiles
- CI/CD configuration

**Deliverable**: Production-ready Pulsar application

## Key Dependencies

### Backend (Go)
```go
github.com/gin-gonic/gin              // HTTP framework
github.com/golang-jwt/jwt/v5          // JWT authentication
github.com/lib/pq                     // PostgreSQL driver
github.com/jmoiron/sqlx               // SQL extensions
github.com/golang-migrate/migrate/v4  // Database migrations
github.com/google/uuid                // UUID generation
github.com/gorilla/websocket          // WebSocket support
github.com/go-playground/validator/v10 // Request validation
github.com/spf13/viper                // Configuration management
go.uber.org/zap                       // Structured logging
golang.org/x/crypto                   // Password hashing (bcrypt)
```

### Frontend (npm)
```json
@sveltejs/kit         // SvelteKit framework
svelte                // Svelte compiler
typescript            // TypeScript support
tailwindcss           // Utility-first CSS
dayjs                 // Date manipulation
zod                   // Schema validation
```

## Security Considerations

1. **Authentication**: JWT with HttpOnly cookies for refresh tokens
2. **Authorization**: RBAC with resource-level checks
3. **Password Security**: Bcrypt hashing with salt
4. **API Keys**: Hashed storage, prefix for identification
5. **SQL Injection**: Parameterized queries via sqlx
6. **XSS Protection**: Content Security Policy, proper escaping
7. **Rate Limiting**: Per user/IP rate limits
8. **CORS**: Configured allowed origins
9. **HTTPS**: Required in production
10. **Audit Logging**: Track all significant actions

## Development Workflow

### Getting Started
```bash
# 1. Set up environment
cd /home/nour/workspace/github/pulsar
cp .env.example .env

# 2. Start services
docker-compose up -d

# 3. Run migrations
make migrate-up

# 4. Access application
# Backend: http://localhost:8080
# Frontend: http://localhost:3000
# Database: localhost:5432
```

### Common Commands (via Makefile)
```bash
make run          # Run backend locally
make test         # Run tests
make migrate-up   # Apply migrations
make migrate-down # Rollback migrations
make build        # Build Docker images
make lint         # Run linters
```

## Next Steps

1. **Start with Phase 1**: Set up foundation and authentication
2. **Iterate incrementally**: Complete each phase before moving to next
3. **Test continuously**: Write tests alongside implementation
4. **Review regularly**: Ensure alignment with requirements
5. **Document as you go**: Update README and API docs

## Success Criteria

- ✅ Users can register and authenticate
- ✅ Alerts can be created, assigned, and managed
- ✅ Teams can be organized with on-call schedules
- ✅ Escalation policies automatically route alerts
- ✅ Notifications sent via multiple channels
- ✅ Incidents tracked with timeline
- ✅ Real-time updates via WebSocket
- ✅ External integrations via webhooks and API
- ✅ Production-ready with monitoring and security
