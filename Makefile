.PHONY: help up down build logs migrate-up migrate-down migrate-create db-reset test clean

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
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean up containers and volumes"

up:
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services started!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Database: localhost:5432"

down:
	docker-compose down

build:
	docker-compose build

logs:
	docker-compose logs -f

migrate-up:
	docker-compose exec backend migrate -path=/app/migrations -database "$${DATABASE_URL}" up

migrate-down:
	docker-compose exec backend migrate -path=/app/migrations -database "$${DATABASE_URL}" down 1

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

clean:
	docker-compose down -v
	rm -rf backend/tmp
	rm -f backend/build-errors.log
