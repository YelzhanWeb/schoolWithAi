.PHONY: help install run test migrate-up migrate-down docker-up docker-down clean

# Variables
BACKEND_DIR=backend
ML_DIR=ml-service
FRONTEND_DIR=frontend

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(GREEN)Education Platform - Available Commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

install: ## Install all dependencies
	@echo "$(GREEN)Installing Go dependencies...$(NC)"
	cd $(BACKEND_DIR) && go mod download && go mod tidy
	@echo "$(GREEN)Installing Python dependencies...$(NC)"
	cd $(ML_DIR) && pip install -r requirements.txt
	@echo "$(GREEN)✅ All dependencies installed$(NC)"

run-backend: ## Run backend server
	@echo "$(GREEN)Starting backend server...$(NC)"
	cd $(BACKEND_DIR) && go run cmd/main.go

run-ml: ## Run ML service
	@echo "$(GREEN)Starting ML service...$(NC)"
	cd $(ML_DIR) && python app/main.py

run-frontend: ## Run frontend (simple HTTP server)
	@echo "$(GREEN)Starting frontend on http://localhost:3000$(NC)"
	cd $(FRONTEND_DIR) && python3 -m http.server 3000

test: ## Run all tests
	@echo "$(GREEN)Running backend tests...$(NC)"
	cd $(BACKEND_DIR) && go test -v ./...
	@echo "$(GREEN)Running ML tests...$(NC)"
	cd $(ML_DIR) && pytest

test-coverage: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	cd $(BACKEND_DIR) && go test -coverprofile=coverage.out ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report: backend/coverage.html$(NC)"

migrate-up: ## Run database migrations up
	@echo "$(GREEN)Running migrations UP...$(NC)"
	cd $(BACKEND_DIR) && go run cmd/main.go -migrate up
	@echo "$(GREEN)✅ Migrations completed$(NC)"

migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)Rolling back migrations...$(NC)"
	cd $(BACKEND_DIR) && go run cmd/main.go -migrate down
	@echo "$(GREEN)✅ Rollback completed$(NC)"

migrate-create: ## Create a new migration (usage: make migrate-create name=add_users_table)
	@echo "$(GREEN)Creating new migration: $(name)$(NC)"
	cd $(BACKEND_DIR)/migrations && \
		touch $(shell date +%Y%m%d%H%M%S)_$(name).up.sql && \
		touch $(shell date +%Y%m%d%H%M%S)_$(name).down.sql
	@echo "$(GREEN)✅ Migration files created$(NC)"

docker-up: ## Start all services with Docker Compose
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✅ Services started:$(NC)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - pgAdmin: http://localhost:5050"
	@docker-compose ps

docker-down: ## Stop all Docker containers
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)✅ Containers stopped$(NC)"

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-clean: ## Remove Docker containers and volumes
	@echo "$(RED)Removing Docker containers and volumes...$(NC)"
	docker-compose down -v
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

db-shell: ## Connect to PostgreSQL shell
	@echo "$(GREEN)Connecting to PostgreSQL...$(NC)"
	docker-compose exec postgres psql -U postgres -d education_platform

db-reset: docker-down docker-clean docker-up migrate-up seed ## Reset database (DESTRUCTIVE!)
	@echo "$(GREEN)✅ Database reset completed$(NC)"

seed: ## Seed database with sample data
	@echo "$(GREEN)Seeding database...$(NC)"
	cd $(BACKEND_DIR) && go run scripts/seed.go
	@echo "$(GREEN)✅ Database seeded$(NC)"

lint: ## Run linters
	@echo "$(GREEN)Running Go linters...$(NC)"
	cd $(BACKEND_DIR) && golangci-lint run
	@echo "$(GREEN)Running Python linters...$(NC)"
	cd $(ML_DIR) && pylint app/

format: ## Format code
	@echo "$(GREEN)Formatting Go code...$(NC)"
	cd $(BACKEND_DIR) && go fmt ./...
	@echo "$(GREEN)Formatting Python code...$(NC)"
	cd $(ML_DIR) && black app/

build: ## Build backend binary
	@echo "$(GREEN)Building backend...$(NC)"
	cd $(BACKEND_DIR) && go build -o bin/education-platform cmd/main.go
	@echo "$(GREEN)✅ Build completed: backend/bin/education-platform$(NC)"

build-docker: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t education-platform:latest .
	@echo "$(GREEN)✅ Docker image built$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(BACKEND_DIR)/coverage.out
	rm -rf $(BACKEND_DIR)/coverage.html
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

dev: docker-up migrate-up ## Setup development environment
	@echo "$(GREEN)✅ Development environment ready!$(NC)"
	@echo ""
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Run backend: make run-backend"
	@echo "  2. Run ML service: make run-ml"
	@echo "  3. Run frontend: make run-frontend"

status: ## Check services status
	@echo "$(GREEN)Checking services status...$(NC)"
	@echo ""
	@echo "$(YELLOW)Docker containers:$(NC)"
	@docker-compose ps
	@echo ""
	@echo "$(YELLOW)Database connection:$(NC)"
	@docker-compose exec postgres pg_isready -U postgres && echo "✅ PostgreSQL is ready" || echo "❌ PostgreSQL is not ready"