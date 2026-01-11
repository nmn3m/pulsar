# ============================================================================
# PULSAR - Incident Management Platform
# ============================================================================
# Usage: make [target]
# Run 'make help' for available commands
# ============================================================================

.DEFAULT_GOAL := help

# ----------------------------------------------------------------------------
# Variables
# ----------------------------------------------------------------------------
PROJECT_NAME    := pulsar
BACKEND_DIR     := backend
FRONTEND_DIR    := frontend
DOCKER_COMPOSE  := docker-compose
DOCKER_TEST     := docker-compose -f docker-compose.test.yml

# Database
DB_USER         := pulsar
DB_PASS         := pulsar_dev_password
DB_NAME         := pulsar
DB_HOST         := postgres
DB_PORT         := 5432
DB_URL          := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Colors
CYAN  := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED   := \033[31m
RESET := \033[0m
BOLD  := \033[1m

# ----------------------------------------------------------------------------
# Help
# ----------------------------------------------------------------------------
.PHONY: help
help:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)Pulsar$(RESET) - Incident Management Platform\n"
	@printf "\n"
	@printf "$(BOLD)Usage:$(RESET) make $(GREEN)<target>$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)Development:$(RESET)\n"
	@printf "  $(GREEN)up$(RESET)                Start all services\n"
	@printf "  $(GREEN)down$(RESET)              Stop all services\n"
	@printf "  $(GREEN)build$(RESET)             Build Docker images\n"
	@printf "  $(GREEN)logs$(RESET)              View container logs\n"
	@printf "  $(GREEN)restart$(RESET)           Restart all services\n"
	@printf "  $(GREEN)ps$(RESET)                Show running containers\n"
	@printf "\n"
	@printf "$(BOLD)Database:$(RESET)\n"
	@printf "  $(GREEN)migrate-up$(RESET)        Run database migrations\n"
	@printf "  $(GREEN)migrate-down$(RESET)      Rollback last migration\n"
	@printf "  $(GREEN)migrate-create$(RESET)    Create new migration (NAME=<name>)\n"
	@printf "  $(GREEN)db-reset$(RESET)          Reset database $(RED)(destructive)$(RESET)\n"
	@printf "  $(GREEN)seed$(RESET)              Seed demo data\n"
	@printf "\n"
	@printf "$(BOLD)Testing:$(RESET)\n"
	@printf "  $(GREEN)test$(RESET)              Run unit tests\n"
	@printf "  $(GREEN)test-integration$(RESET)  Run integration tests\n"
	@printf "  $(GREEN)test-coverage$(RESET)     Generate coverage report\n"
	@printf "\n"
	@printf "$(BOLD)Code Quality:$(RESET)\n"
	@printf "  $(GREEN)lint$(RESET)              Run all linters\n"
	@printf "  $(GREEN)fmt$(RESET)               Format all code\n"
	@printf "  $(GREEN)fmt-check$(RESET)         Check code formatting\n"
	@printf "  $(GREEN)vet$(RESET)               Run go vet\n"
	@printf "\n"
	@printf "$(BOLD)Utilities:$(RESET)\n"
	@printf "  $(GREEN)swagger$(RESET)           Generate Swagger docs\n"
	@printf "  $(GREEN)clean$(RESET)             Clean up containers and volumes\n"
	@printf "\n"

# ----------------------------------------------------------------------------
# Development
# ----------------------------------------------------------------------------
.PHONY: up down build logs restart ps

up:
	@printf "$(CYAN)Starting services...$(RESET)\n"
	@$(DOCKER_COMPOSE) up -d
	@printf "$(YELLOW)Waiting for services to be ready...$(RESET)\n"
	@sleep 5
	@printf "\n"
	@printf "$(GREEN)$(BOLD)Services started successfully!$(RESET)\n"
	@printf "\n"
	@printf "  Frontend:  $(CYAN)http://localhost:5173$(RESET)\n"
	@printf "  API:       $(CYAN)http://localhost:8081$(RESET)\n"
	@printf "  Swagger:   $(CYAN)http://localhost:8081/swagger/index.html$(RESET)\n"
	@printf "  Mailpit:   $(CYAN)http://localhost:8025$(RESET)\n"
	@printf "  Database:  $(CYAN)localhost:5433$(RESET)\n"
	@printf "\n"

down:
	@printf "$(CYAN)Stopping services...$(RESET)\n"
	@$(DOCKER_COMPOSE) down
	@printf "$(GREEN)Services stopped.$(RESET)\n"

build:
	@printf "$(CYAN)Building Docker images...$(RESET)\n"
	@$(DOCKER_COMPOSE) build
	@printf "$(GREEN)Build complete.$(RESET)\n"

logs:
	@$(DOCKER_COMPOSE) logs -f

restart: down up

ps:
	@$(DOCKER_COMPOSE) ps

# ----------------------------------------------------------------------------
# Database
# ----------------------------------------------------------------------------
.PHONY: migrate-up migrate-down migrate-create db-reset seed

migrate-up:
	@printf "$(CYAN)Running migrations...$(RESET)\n"
	@$(DOCKER_COMPOSE) exec -T backend migrate -path=/app/migrations -database "$(DB_URL)" up
	@printf "$(GREEN)Migrations complete.$(RESET)\n"

migrate-down:
	@printf "$(YELLOW)Rolling back last migration...$(RESET)\n"
	@$(DOCKER_COMPOSE) exec -T backend migrate -path=/app/migrations -database "$(DB_URL)" down 1
	@printf "$(GREEN)Rollback complete.$(RESET)\n"

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		printf "$(RED)Error: NAME is required$(RESET)\n"; \
		printf "Usage: make migrate-create NAME=<migration_name>\n"; \
		exit 1; \
	fi
	@printf "$(CYAN)Creating migration: $(NAME)$(RESET)\n"
	@$(DOCKER_COMPOSE) exec backend migrate create -ext sql -dir /app/migrations -seq $(NAME)

db-reset:
	@printf "$(RED)$(BOLD)WARNING: This will delete all data!$(RESET)\n"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		printf "$(CYAN)Resetting database...$(RESET)\n"; \
		$(DOCKER_COMPOSE) down -v; \
		$(DOCKER_COMPOSE) up -d postgres; \
		sleep 5; \
		$(MAKE) migrate-up; \
		printf "$(GREEN)Database reset complete.$(RESET)\n"; \
	else \
		printf "$(YELLOW)Cancelled.$(RESET)\n"; \
	fi

seed:
	@printf "$(CYAN)Seeding demo data...$(RESET)\n"
	@cd $(BACKEND_DIR) && \
		DATABASE_URL="postgres://$(DB_USER):$(DB_PASS)@localhost:5433/$(DB_NAME)?sslmode=disable" \
		JWT_SECRET="dev_jwt_secret_change_in_production_min_32_chars" \
		JWT_REFRESH_SECRET="dev_refresh_secret_change_in_production_min_32_chars" \
		go run ./cmd/seed/main.go
	@printf "$(GREEN)Seeding complete.$(RESET)\n"

# ----------------------------------------------------------------------------
# Testing
# ----------------------------------------------------------------------------
.PHONY: test test-db-up test-db-down test-integration test-integration-verbose test-coverage

test:
	@printf "$(CYAN)Running unit tests...$(RESET)\n"
	@cd $(BACKEND_DIR) && go test -v ./...

test-db-up:
	@printf "$(CYAN)Starting test database...$(RESET)\n"
	@$(DOCKER_TEST) up -d
	@sleep 3
	@printf "$(GREEN)Test database ready at localhost:5434$(RESET)\n"

test-db-down:
	@$(DOCKER_TEST) down -v

test-integration: test-db-up
	@printf "$(CYAN)Running integration tests...$(RESET)\n"
	@cd $(BACKEND_DIR) && go test -v ./tests/integration/... -count=1 || ($(MAKE) test-db-down && exit 1)
	@$(MAKE) test-db-down
	@printf "$(GREEN)Integration tests passed.$(RESET)\n"

test-integration-verbose: test-db-up
	@printf "$(CYAN)Running integration tests (verbose)...$(RESET)\n"
	@cd $(BACKEND_DIR) && go test -v -race ./tests/integration/... -count=1 || ($(MAKE) test-db-down && exit 1)
	@$(MAKE) test-db-down

test-coverage: test-db-up
	@printf "$(CYAN)Running tests with coverage...$(RESET)\n"
	@cd $(BACKEND_DIR) && go test -v -coverprofile=coverage.out -covermode=atomic ./tests/integration/... || ($(MAKE) test-db-down && exit 1)
	@cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@$(MAKE) test-db-down
	@printf "$(GREEN)Coverage report: $(BACKEND_DIR)/coverage.html$(RESET)\n"

# ----------------------------------------------------------------------------
# Code Quality
# ----------------------------------------------------------------------------
.PHONY: lint lint-backend lint-frontend fmt fmt-backend fmt-frontend fmt-check fmt-check-backend fmt-check-frontend vet

lint: lint-backend lint-frontend
	@printf "$(GREEN)Linting complete.$(RESET)\n"

lint-backend:
	@printf "$(CYAN)Linting backend...$(RESET)\n"
	@cd $(BACKEND_DIR) && golangci-lint run ./...

lint-frontend:
	@printf "$(CYAN)Linting frontend...$(RESET)\n"
	@cd $(FRONTEND_DIR) && npm run lint

fmt: fmt-backend fmt-frontend
	@printf "$(GREEN)Formatting complete.$(RESET)\n"

fmt-backend:
	@printf "$(CYAN)Formatting backend...$(RESET)\n"
	@cd $(BACKEND_DIR) && gofmt -w .
	@cd $(BACKEND_DIR) && goimports -w . 2>/dev/null || true

fmt-frontend:
	@printf "$(CYAN)Formatting frontend...$(RESET)\n"
	@cd $(FRONTEND_DIR) && npm run format

fmt-check: fmt-check-backend fmt-check-frontend

fmt-check-backend:
	@printf "$(CYAN)Checking backend formatting...$(RESET)\n"
	@cd $(BACKEND_DIR) && if [ -n "$$(gofmt -l .)" ]; then \
		printf "$(RED)Backend: Files not formatted:$(RESET)\n"; \
		gofmt -l .; \
		exit 1; \
	else \
		printf "$(GREEN)Backend: All files formatted$(RESET)\n"; \
	fi

fmt-check-frontend:
	@printf "$(CYAN)Checking frontend formatting...$(RESET)\n"
	@cd $(FRONTEND_DIR) && npm run format:check

vet:
	@printf "$(CYAN)Running go vet...$(RESET)\n"
	@cd $(BACKEND_DIR) && go vet ./...
	@printf "$(GREEN)Vet complete.$(RESET)\n"

# ----------------------------------------------------------------------------
# Utilities
# ----------------------------------------------------------------------------
.PHONY: swagger clean

swagger:
	@printf "$(CYAN)Generating Swagger documentation...$(RESET)\n"
	@cd $(BACKEND_DIR) && swag init -g cmd/api/main.go -o docs
	@printf "$(GREEN)Swagger docs generated.$(RESET)\n"

clean:
	@printf "$(CYAN)Cleaning up...$(RESET)\n"
	@$(DOCKER_COMPOSE) down -v 2>/dev/null || true
	@$(DOCKER_TEST) down -v 2>/dev/null || true
	@rm -rf $(BACKEND_DIR)/tmp
	@rm -f $(BACKEND_DIR)/build-errors.log
	@rm -f $(BACKEND_DIR)/coverage.out $(BACKEND_DIR)/coverage.html
	@printf "$(GREEN)Cleanup complete.$(RESET)\n"
