# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

EXECUTOR := docker compose

# if podman is available, use it instead of docker
ifneq (, $(shell which podman 2>/dev/null))
$(info ❗❗ Using podman for compose commands❗❗)
EXECUTOR := podman compose
endif

# Get absolute path to project root
ROOT_DIR := $(shell pwd)
COMPOSE := $(EXECUTOR) -f $(ROOT_DIR)/Ops/docker-compose.yml --env-file $(ROOT_DIR)/Ops/.docker.env

DEV_SERVICE := postgres

# =============================================================================
# Check .env file exists
# =============================================================================

.PHONY: check-env
check-env: ## Check if .env or Ops/.docker.env file exists
	@if [ ! -f .env ] || [ ! -f Ops/.docker.env ]; then \
		echo -e "$(RED)❌ .env or Ops/.docker.env file not found!$(RESET)"; \
		echo -e "$(YELLOW)💡 Copy .env.sample to .env and Ops/.docker.env.sample to Ops/.docker.env and configure it:$(RESET)"; \
		echo "   cp .env.sample .env"; \
		echo "   cp Ops/.docker.env.sample Ops/.docker.env"; \
		exit 1; \
	fi

# =============================================================================
# Dev Profile - Development commands
# =============================================================================

.PHONY: up
up: check-env ## Start only Postgres for development (no app)
	@echo -e "$(CYAN)🔧 Starting development environment...$(RESET)"
	@$(COMPOSE) up $(DEV_SERVICE) -d
	@echo -e "$(GREEN)✅ Development environment ready!$(RESET)"

.PHONY: down
down: ## Stop development environment
	@echo -e "$(CYAN)🛑 Stopping development environment...$(RESET)"
	@$(COMPOSE) down

# =============================================================================
# Database commands
# =============================================================================

.PHONY: db-shell
db-shell: check-env ## Open psql shell in PostgreSQL
	@echo -e "$(CYAN)🐘 Opening PostgreSQL shell...$(RESET)"
	@$(COMPOSE) exec postgres psql -h localhost -p 5432 -U $(shell grep POSTGRES_USER .env | cut -d '=' -f2) -d $(shell grep POSTGRES_DB .env | cut -d '=' -f2)

.PHONY: db-logs
db-logs: ## Show PostgreSQL logs
	@$(COMPOSE) logs -f postgres

.PHONY: db-migrate-create
db-migrate-create: check-env ## Create a new database migration
	@echo -e "$(CYAN)🔄 Creating a new database migration...$(RESET)"
	@read -p "Enter the migration name: " name; \
	goose -dir db/migration create $$name sql
	@echo -e "$(GREEN)✅ Migration created$(RESET)"

.PHONY: db-migrate-status
db-migrate-status: check-env ## Check database migration status
	@echo -e "$(CYAN)🔄 Checking database migration status...$(RESET)"
	goose -dir db/migration postgres "$(shell grep DATABASE_URL .env | cut -d '=' -f2)" status
	@echo -e "$(GREEN)✅ Migrations status checked$(RESET)"

.PHONY: db-migrate-up
db-migrate-up: check-env ## Apply database migrations
	@echo -e "$(CYAN)🔄 Applying database migrations...$(RESET)"
	goose -dir db/migration postgres "$(shell grep DATABASE_URL .env | cut -d '=' -f2)" up
	@echo -e "$(GREEN)✅ Migrations applied$(RESET)"

.PHONY: db-migrate-down
db-migrate-down: check-env ## Rollback database migrations
	@echo -e "$(CYAN)🔄 Rolling back database migrations...$(RESET)"
	goose -dir db/migration postgres "$(shell grep DATABASE_URL .env | cut -d '=' -f2)" down
	@echo -e "$(GREEN)✅ Migrations rolled back$(RESET)"

.PHONY: db-generate
db-generate: check-env ## Generate a database models
	@echo -e "$(CYAN)🔄 Generating database models...$(RESET)"
	sqlc generate
	@echo -e "$(GREEN)✅ Database models generated$(RESET)"

# =============================================================================
# Utility commands
# =============================================================================

.PHONY: logs
logs: ## Show logs for all services
	@$(COMPOSE) logs -f

.PHONY: env-check
env-check: check-env ## Verify .env configuration
	@echo -e "$(CYAN)🔍 Environment Configuration:$(RESET)"
	@echo -e "$(YELLOW)POSTGRES_PORT:$(RESET)       $(shell grep POSTGRES_PORT .env | cut -d '=' -f2)"
	@echo -e "$(YELLOW)POSTGRES_USER:$(RESET)       $(shell grep POSTGRES_USER .env | cut -d '=' -f2)"
	@echo -e "$(YELLOW)POSTGRES_PASSWORD:$(RESET)   $(shell grep POSTGRES_PASSWORD .env | cut -d '=' -f2)"
	@echo -e "$(YELLOW)POSTGRES_DB:$(RESET)         $(shell grep POSTGRES_DB .env | cut -d '=' -f2)"
	@echo -e "$(YELLOW)DATABASE_URL:$(RESET)        $(shell grep DATABASE_URL .env | cut -d '=' -f2)"

# =============================================================================
# Help
# =============================================================================

.PHONY: help
help: ## Show this help message
	@echo -e "$(CYAN)Search Engine$(RESET)" - Available Commands:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Default target
.DEFAULT_GOAL := help

# Allow arguments to be passed to make commands
%:
	@: