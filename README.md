# Pulsar - Incident Management Platform

An open-source Opsgenie replacement built with Go and Svelte.

## Features

### Phase 1 (Completed)
- ✅ User authentication (JWT-based)
- ✅ User registration and login
- ✅ Organization management
- ✅ Multi-tenancy support
- ✅ Role-based access control (RBAC)

### Upcoming Features
- Alert Management
- On-Call Schedules
- Incident Management
- Escalation Policies
- Multi-channel Notifications (Email, Slack, Teams)
- Webhooks & Integrations
- Real-time Updates (WebSocket)
- API Keys

## Tech Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL 16
- **Authentication**: JWT
- **Logging**: Uber Zap

### Frontend
- **Framework**: SvelteKit
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Build Tool**: Vite

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Development**: Live reload with Air (backend) and Vite (frontend)

## Quick Start

### Prerequisites
- Docker and Docker Compose
- (Optional) Go 1.21+ and Node.js 20+ for local development

### Getting Started

1. **Clone the repository**
   ```bash
   cd /home/nour/workspace/github/pulsar
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env if needed (defaults work for development)
   ```

3. **Start all services**
   ```bash
   make up
   ```

   This will start:
   - PostgreSQL database on `localhost:5432`
   - Backend API on `http://localhost:8080`
   - Frontend on `http://localhost:3000`

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Access the application**
   - Open your browser to `http://localhost:3000`
   - Create a new account via the registration page
   - Log in and explore the dashboard

### Common Commands

```bash
make up            # Start all services
make down          # Stop all services
make logs          # View logs
make migrate-up    # Run database migrations
make migrate-down  # Rollback last migration
make clean         # Clean up containers and volumes
```

## Project Structure

```
pulsar/
├── backend/              # Go backend
│   ├── cmd/
│   │   └── api/         # API server entry point
│   ├── internal/
│   │   ├── config/      # Configuration
│   │   ├── domain/      # Domain models
│   │   ├── repository/  # Data access
│   │   ├── service/     # Business logic
│   │   ├── handler/     # HTTP handlers
│   │   ├── middleware/  # HTTP middleware
│   │   └── pkg/         # Internal packages
│   └── migrations/      # Database migrations
│
├── frontend/            # Svelte frontend
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/     # API client
│   │   │   ├── stores/  # State management
│   │   │   ├── components/ # UI components
│   │   │   └── types/   # TypeScript types
│   │   └── routes/      # SvelteKit routes
│   │       ├── (auth)/  # Login, register
│   │       └── (app)/   # Dashboard, alerts, etc.
│
└── docker-compose.yml   # Development environment
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user (protected)

### Health
- `GET /health` - Health check

## Development

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run with live reload
air

# Run tests
go test -v ./...
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## Database Migrations

Create a new migration:
```bash
make migrate-create NAME=add_something
```

Apply migrations:
```bash
make migrate-up
```

Rollback last migration:
```bash
make migrate-down
```

## Architecture

Pulsar follows clean architecture principles:

### Backend Layers
1. **Domain Layer**: Core business entities (User, Organization, Alert, etc.)
2. **Repository Layer**: Data access abstraction
3. **Service Layer**: Business logic
4. **Handler Layer**: HTTP request/response handling
5. **Middleware**: Cross-cutting concerns (auth, logging, CORS)

### Frontend Structure
1. **API Client**: Centralized API communication
2. **Stores**: Svelte stores for state management
3. **Components**: Reusable UI components
4. **Routes**: SvelteKit file-based routing

## Security

- JWT-based authentication with access and refresh tokens
- Password hashing with bcrypt
- CORS configuration
- Role-based access control (RBAC)
- SQL injection prevention via parameterized queries
- XSS protection

## Next Steps (Phase 2+)

See [plan.md](./plan.md) for the complete implementation roadmap.

Phase 2 will focus on Alert Management:
- Create, view, and manage alerts
- Alert filtering and search
- Alert assignment to users/teams
- Alert acknowledgment and closure

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - feel free to use this project for any purpose.

## Support

For questions and support, please open an issue on GitHub.

---

Built with ❤️ using Go and Svelte
