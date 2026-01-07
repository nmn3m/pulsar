.PHONY: help up down build logs migrate-up migrate-down migrate-create db-reset test test-db-up test-db-down test-integration test-integration-verbose test-coverage clean \
	lint lint-backend lint-frontend fmt fmt-backend fmt-frontend fmt-check fmt-check-backend fmt-check-frontend vet

help:
	@echo "Pulsar - Development Commands"
	@echo ""
	@echo "  make up            - Start all services"
	@echo "  make down          - Stop all services"
	@echo "  make build         - Build all Docker images"
	@echo "  make logs          - View logs"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback last migration"
	@echo "  make migrate-create NAME=<name> - Create new migration"
	@echo "  make db-reset      - Reset database (WARNING: destructive)"
	@echo "  make test          - Run unit tests"
	@echo "  make test-db-up    - Start test database"
	@echo "  make test-db-down  - Stop test database"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean up containers and volumes"
	@echo ""
	@echo "Linting & Formatting:"
	@echo "  make lint          - Run linters (backend + frontend)"
	@echo "  make lint-backend  - Run Go linters (golangci-lint)"
	@echo "  make lint-frontend - Run ESLint"
	@echo "  make fmt           - Format code (backend + frontend)"
	@echo "  make fmt-backend   - Format Go code (gofmt + goimports)"
	@echo "  make fmt-frontend  - Format frontend code (Prettier)"
	@echo "  make fmt-check     - Check formatting (backend + frontend)"
	@echo "  make vet           - Run go vet"

up:
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services started!"
	@echo "Backend: http://localhost:8081"
	@echo "Frontend: http://localhost:5173"
	@echo "Database: localhost:5433"

down:
	docker-compose down

build:
	docker-compose build

logs:
	docker-compose logs -f

migrate-up:
	docker-compose exec -T backend migrate -path=/app/migrations -database "postgres://pulsar:pulsar_dev_password@postgres:5432/pulsar?sslmode=disable" up

migrate-down:
	docker-compose exec -T backend migrate -path=/app/migrations -database "postgres://pulsar:pulsar_dev_password@postgres:5432/pulsar?sslmode=disable" down 1

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=<migration_name>"; \
		exit 1; \
	fi
	docker-compose exec backend migrate create -ext sql -dir /app/migrations -seq $(NAME)

db-reset:
	@echo "WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker-compose up -d postgres; \
		sleep 5; \
		make migrate-up; \
	fi

test:
	cd backend && go test -v ./...

# Integration test targets
test-db-up:
	docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for test database to be ready..."
	@sleep 3
	@echo "Test database ready at localhost:5434"

test-db-down:
	docker-compose -f docker-compose.test.yml down -v

test-integration: test-db-up
	cd backend && go test -v ./tests/integration/... -count=1 || ($(MAKE) test-db-down && exit 1)
	$(MAKE) test-db-down

test-integration-verbose: test-db-up
	cd backend && go test -v -race ./tests/integration/... -count=1 || ($(MAKE) test-db-down && exit 1)
	$(MAKE) test-db-down

test-coverage: test-db-up
	cd backend && go test -v -coverprofile=coverage.out -covermode=atomic ./tests/integration/... || ($(MAKE) test-db-down && exit 1)
	cd backend && go tool cover -html=coverage.out -o coverage.html
	$(MAKE) test-db-down
	@echo "Coverage report generated at backend/coverage.html"

clean:
	docker-compose down -v
	docker-compose -f docker-compose.test.yml down -v 2>/dev/null || true
	rm -rf backend/tmp
	rm -f backend/build-errors.log
	rm -f backend/coverage.out backend/coverage.html

# Linting targets
lint: lint-backend lint-frontend

lint-backend:
	cd backend && golangci-lint run ./...

lint-frontend:
	cd frontend && npm run lint

# Formatting targets
fmt: fmt-backend fmt-frontend

fmt-backend:
	cd backend && gofmt -w .
	cd backend && goimports -w .

fmt-frontend:
	cd frontend && npm run format

# Format check targets
fmt-check: fmt-check-backend fmt-check-frontend

fmt-check-backend:
	@cd backend && if [ -n "$$(gofmt -l .)" ]; then \
		echo "Backend: The following files are not formatted:"; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "Backend: All files are properly formatted"; \
	fi

fmt-check-frontend:
	cd frontend && npm run format:check

# Go vet
vet:
	cd backend && go vet ./...
