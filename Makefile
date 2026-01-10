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
	@echo ""
	@echo "$(BOLD)$(CYAN)Pulsar$(RESET) - Incident Management Platform"
	@echo ""
	@echo "$(BOLD)Usage:$(RESET) make $(GREEN)<target>$(RESET)"
	@echo ""
	@echo "$(BOLD)Development:$(RESET)"
	@echo "  $(GREEN)up$(RESET)                Start all services"
	@echo "  $(GREEN)down$(RESET)              Stop all services"
	@echo "  $(GREEN)build$(RESET)             Build Docker images"
	@echo "  $(GREEN)logs$(RESET)              View container logs"
	@echo "  $(GREEN)restart$(RESET)           Restart all services"
	@echo "  $(GREEN)ps$(RESET)                Show running containers"
	@echo ""
	@echo "$(BOLD)Database:$(RESET)"
	@echo "  $(GREEN)migrate-up$(RESET)        Run database migrations"
	@echo "  $(GREEN)migrate-down$(RESET)      Rollback last migration"
	@echo "  $(GREEN)migrate-create$(RESET)    Create new migration (NAME=<name>)"
	@echo "  $(GREEN)db-reset$(RESET)          Reset database $(RED)(destructive)$(RESET)"
	@echo "  $(GREEN)seed$(RESET)              Seed demo data"
	@echo ""
	@echo "$(BOLD)Testing:$(RESET)"
	@echo "  $(GREEN)test$(RESET)              Run unit tests"
	@echo "  $(GREEN)test-integration$(RESET)  Run integration tests"
	@echo "  $(GREEN)test-coverage$(RESET)     Generate coverage report"
	@echo ""
	@echo "$(BOLD)Code Quality:$(RESET)"
	@echo "  $(GREEN)lint$(RESET)              Run all linters"
	@echo "  $(GREEN)fmt$(RESET)               Format all code"
	@echo "  $(GREEN)fmt-check$(RESET)         Check code formatting"
	@echo "  $(GREEN)vet$(RESET)               Run go vet"
	@echo ""
	@echo "$(BOLD)Utilities:$(RESET)"
	@echo "  $(GREEN)swagger$(RESET)           Generate Swagger docs"
	@echo "  $(GREEN)clean$(RESET)             Clean up containers and volumes"
	@echo ""

# ----------------------------------------------------------------------------
# Development
# ----------------------------------------------------------------------------
.PHONY: up down build logs restart ps

up:
	@echo "$(CYAN)Starting services...$(RESET)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(YELLOW)Waiting for services to be ready...$(RESET)"
	@sleep 5
	@echo ""
	@echo "$(GREEN)$(BOLD)Services started successfully!$(RESET)"
	@echo ""
	@echo "  Frontend:  $(CYAN)http://localhost:5173$(RESET)"
	@echo "  API:       $(CYAN)http://localhost:8081$(RESET)"
	@echo "  Swagger:   $(CYAN)http://localhost:8081/swagger/index.html$(RESET)"
	@echo "  Database:  $(CYAN)localhost:5433$(RESET)"
	@echo ""

down:
	@echo "$(CYAN)Stopping services...$(RESET)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Services stopped.$(RESET)"

build:
	@echo "$(CYAN)Building Docker images...$(RESET)"
	@$(DOCKER_COMPOSE) build
	@echo "$(GREEN)Build complete.$(RESET)"

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
	@echo "$(CYAN)Running migrations...$(RESET)"
	@$(DOCKER_COMPOSE) exec -T backend migrate -path=/app/migrations -database "$(DB_URL)" up
	@echo "$(GREEN)Migrations complete.$(RESET)"

migrate-down:
	@echo "$(YELLOW)Rolling back last migration...$(RESET)"
	@$(DOCKER_COMPOSE) exec -T backend migrate -path=/app/migrations -database "$(DB_URL)" down 1
	@echo "$(GREEN)Rollback complete.$(RESET)"

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required$(RESET)"; \
		echo "Usage: make migrate-create NAME=<migration_name>"; \
		exit 1; \
	fi
	@echo "$(CYAN)Creating migration: $(NAME)$(RESET)"
	@$(DOCKER_COMPOSE) exec backend migrate create -ext sql -dir /app/migrations -seq $(NAME)

db-reset:
	@echo "$(RED)$(BOLD)WARNING: This will delete all data!$(RESET)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(CYAN)Resetting database...$(RESET)"; \
		$(DOCKER_COMPOSE) down -v; \
		$(DOCKER_COMPOSE) up -d postgres; \
		sleep 5; \
		$(MAKE) migrate-up; \
		echo "$(GREEN)Database reset complete.$(RESET)"; \
	else \
		echo "$(YELLOW)Cancelled.$(RESET)"; \
	fi

seed:
	@echo "$(CYAN)Seeding demo data...$(RESET)"
	@cd $(BACKEND_DIR) && \
		DATABASE_URL="postgres://$(DB_USER):$(DB_PASS)@localhost:5433/$(DB_NAME)?sslmode=disable" \
		JWT_SECRET="dev_jwt_secret_change_in_production_min_32_chars" \
		JWT_REFRESH_SECRET="dev_refresh_secret_change_in_production_min_32_chars" \
		go run ./cmd/seed/main.go
	@echo "$(GREEN)Seeding complete.$(RESET)"

# ----------------------------------------------------------------------------
# Testing
# ----------------------------------------------------------------------------
.PHONY: test test-db-up test-db-down test-integration test-integration-verbose test-coverage

test:
	@echo "$(CYAN)Running unit tests...$(RESET)"
	@cd $(BACKEND_DIR) && go test -v ./...

test-db-up:
	@echo "$(CYAN)Starting test database...$(RESET)"
	@$(DOCKER_TEST) up -d
	@sleep 3
	@echo "$(GREEN)Test database ready at localhost:5434$(RESET)"

test-db-down:
	@$(DOCKER_TEST) down -v

test-integration: test-db-up
	@echo "$(CYAN)Running integration tests...$(RESET)"
	@cd $(BACKEND_DIR) && go test -v ./tests/integration/... -count=1 || ($(MAKE) test-db-down && exit 1)
	@$(MAKE) test-db-down
	@echo "$(GREEN)Integration tests passed.$(RESET)"

test-integration-verbose: test-db-up
	@echo "$(CYAN)Running integration tests (verbose)...$(RESET)"
	@cd $(BACKEND_DIR) && go test -v -race ./tests/integration/... -count=1 || ($(MAKE) test-db-down && exit 1)
	@$(MAKE) test-db-down

test-coverage: test-db-up
	@echo "$(CYAN)Running tests with coverage...$(RESET)"
	@cd $(BACKEND_DIR) && go test -v -coverprofile=coverage.out -covermode=atomic ./tests/integration/... || ($(MAKE) test-db-down && exit 1)
	@cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@$(MAKE) test-db-down
	@echo "$(GREEN)Coverage report: $(BACKEND_DIR)/coverage.html$(RESET)"

# ----------------------------------------------------------------------------
# Code Quality
# ----------------------------------------------------------------------------
.PHONY: lint lint-backend lint-frontend fmt fmt-backend fmt-frontend fmt-check fmt-check-backend fmt-check-frontend vet

lint: lint-backend lint-frontend
	@echo "$(GREEN)Linting complete.$(RESET)"

lint-backend:
	@echo "$(CYAN)Linting backend...$(RESET)"
	@cd $(BACKEND_DIR) && golangci-lint run ./...

lint-frontend:
	@echo "$(CYAN)Linting frontend...$(RESET)"
	@cd $(FRONTEND_DIR) && npm run lint

fmt: fmt-backend fmt-frontend
	@echo "$(GREEN)Formatting complete.$(RESET)"

fmt-backend:
	@echo "$(CYAN)Formatting backend...$(RESET)"
	@cd $(BACKEND_DIR) && gofmt -w .
	@cd $(BACKEND_DIR) && goimports -w . 2>/dev/null || true

fmt-frontend:
	@echo "$(CYAN)Formatting frontend...$(RESET)"
	@cd $(FRONTEND_DIR) && npm run format

fmt-check: fmt-check-backend fmt-check-frontend

fmt-check-backend:
	@echo "$(CYAN)Checking backend formatting...$(RESET)"
	@cd $(BACKEND_DIR) && if [ -n "$$(gofmt -l .)" ]; then \
		echo "$(RED)Backend: Files not formatted:$(RESET)"; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "$(GREEN)Backend: All files formatted$(RESET)"; \
	fi

fmt-check-frontend:
	@echo "$(CYAN)Checking frontend formatting...$(RESET)"
	@cd $(FRONTEND_DIR) && npm run format:check

vet:
	@echo "$(CYAN)Running go vet...$(RESET)"
	@cd $(BACKEND_DIR) && go vet ./...
	@echo "$(GREEN)Vet complete.$(RESET)"

# ----------------------------------------------------------------------------
# Utilities
# ----------------------------------------------------------------------------
.PHONY: swagger clean

swagger:
	@echo "$(CYAN)Generating Swagger documentation...$(RESET)"
	@cd $(BACKEND_DIR) && swag init -g cmd/api/main.go -o docs
	@echo "$(GREEN)Swagger docs generated.$(RESET)"

clean:
	@echo "$(CYAN)Cleaning up...$(RESET)"
	@$(DOCKER_COMPOSE) down -v 2>/dev/null || true
	@$(DOCKER_TEST) down -v 2>/dev/null || true
	@rm -rf $(BACKEND_DIR)/tmp
	@rm -f $(BACKEND_DIR)/build-errors.log
	@rm -f $(BACKEND_DIR)/coverage.out $(BACKEND_DIR)/coverage.html
	@echo "$(GREEN)Cleanup complete.$(RESET)"
