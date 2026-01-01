.PHONY: help build run test clean docker-build docker-up docker-down docker-logs docker-restart swagger install deps migrate

# Variables
APP_NAME=goapi
DOCKER_COMPOSE=docker-compose
GO=go
SWAG=swag

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(GREEN)Available commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# Local Development
install: ## Install Go dependencies
	@echo "$(GREEN)Installing dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

deps: install ## Alias for install

run: ## Run the application locally
	@echo "$(GREEN)Running application...$(NC)"
	$(GO) run main.go

build: ## Build the application binary
	@echo "$(GREEN)Building application...$(NC)"
	$(GO) build -o $(APP_NAME) main.go
	@echo "$(GREEN)Build complete: $(APP_NAME)$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -f $(APP_NAME)
	rm -f $(APP_NAME).exe
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean complete$(NC)"

# Swagger Documentation
swagger: ## Generate Swagger documentation
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@which $(SWAG) > /dev/null || ($(GO) install github.com/swaggo/swag/cmd/swag@latest && echo "$(GREEN)swag installed$(NC)")
	$(SWAG) init
	@echo "$(GREEN)Swagger docs generated in docs/$(NC)"

swag: ## Install swag CLI tool
	@echo "$(GREEN)Installing swag CLI...$(NC)"
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)swag installed$(NC)"

# Docker Commands
docker-build: ## Build Docker images
	@echo "$(GREEN)Building Docker images...$(NC)"
	$(DOCKER_COMPOSE) build
	@echo "$(GREEN)Docker build complete$(NC)"

docker-up: ## Start Docker containers
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Containers started$(NC)"
	@echo "$(GREEN)API: http://localhost:8080$(NC)"
	@echo "$(GREEN)Swagger: http://localhost:8080/swagger/index.html$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Containers stopped$(NC)"

docker-down-volumes: ## Stop Docker containers and remove volumes
	@echo "$(YELLOW)Stopping containers and removing volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)Containers stopped and volumes removed$(NC)"

docker-logs: ## View Docker container logs
	@echo "$(GREEN)Viewing Docker logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f

docker-logs-api: ## View API container logs only
	@echo "$(GREEN)Viewing API logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f api

docker-logs-db: ## View database container logs only
	@echo "$(GREEN)Viewing database logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f db

docker-logs-redis: ## View Redis container logs only
	@echo "$(GREEN)Viewing Redis logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f redis

docker-restart: ## Restart Docker containers
	@echo "$(GREEN)Restarting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) restart
	@echo "$(GREEN)Containers restarted$(NC)"

docker-rebuild: ## Rebuild and restart Docker containers
	@echo "$(GREEN)Rebuilding and restarting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) up -d --build
	@echo "$(GREEN)Containers rebuilt and restarted$(NC)"

docker-ps: ## Show running Docker containers
	@echo "$(GREEN)Running containers:$(NC)"
	$(DOCKER_COMPOSE) ps

docker-shell-api: ## Open shell in API container
	@echo "$(GREEN)Opening shell in API container...$(NC)"
	$(DOCKER_COMPOSE) exec api sh

docker-shell-db: ## Open PostgreSQL shell in database container
	@echo "$(GREEN)Opening PostgreSQL shell...$(NC)"
	$(DOCKER_COMPOSE) exec db psql -U goapi_user -d goapi

docker-shell-redis: ## Open Redis CLI in Redis container
	@echo "$(GREEN)Opening Redis CLI...$(NC)"
	$(DOCKER_COMPOSE) exec redis redis-cli

docker-logs-redis-commander: ## View Redis Commander logs
	@echo "$(GREEN)Viewing Redis Commander logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f redis-commander

docker-logs-pgadmin: ## View pgAdmin logs
	@echo "$(GREEN)Viewing pgAdmin logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f pgadmin

docker-open-redis-commander: ## Open Redis Commander in browser
	@echo "$(GREEN)Opening Redis Commander at http://localhost:8081$(NC)"
	@echo "$(YELLOW)Username: admin$(NC)"
	@echo "$(YELLOW)Password: admin$(NC)"
	@$(if $(shell which start 2>/dev/null),start http://localhost:8081,echo "Please open http://localhost:8081 in your browser")

docker-open-pgadmin: ## Open pgAdmin in browser
	@echo "$(GREEN)Opening pgAdmin at http://localhost:5050$(NC)"
	@echo "$(YELLOW)Email: admin@goapi.com$(NC)"
	@echo "$(YELLOW)Password: admin$(NC)"
	@$(if $(shell which start 2>/dev/null),start http://localhost:5050,echo "Please open http://localhost:5050 in your browser")

# Database Commands
db-migrate: ## Run database migrations (local)
	@echo "$(GREEN)Running database migrations...$(NC)"
	$(GO) run main.go migrate

db-seed: ## Seed database with sample data (if implemented)
	@echo "$(GREEN)Seeding database...$(NC)"
	@echo "$(YELLOW)Not implemented yet$(NC)"

# Development Workflow
dev: install run ## Install dependencies and run locally

dev-docker: docker-up docker-logs-api ## Start Docker and follow API logs

# Full Setup
setup: install swagger ## Full setup: install deps and generate Swagger docs
	@echo "$(GREEN)Setup complete!$(NC)"
	@echo "$(GREEN)Run 'make run' to start the application$(NC)"
	@echo "$(GREEN)Or 'make docker-up' to start with Docker$(NC)"

# Production Build
prod-build: clean build ## Production build: clean and build
	@echo "$(GREEN)Production build complete: $(APP_NAME)$(NC)"

# All-in-one commands
all: clean install swagger build ## Clean, install, generate docs, and build
	@echo "$(GREEN)All tasks complete!$(NC)"

docker-all: docker-down docker-build docker-up ## Full Docker rebuild: down, build, up
	@echo "$(GREEN)Docker stack ready!$(NC)"

