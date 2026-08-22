.PHONY: dev prod down db-up dev-up db-down psql tidy fmt vet lint test test-verbose test-coverage \
	migrate-diff migrate-up migrate-down migrate-status migrate-prod-up migrate-prod-down migrate-prod-status

MIGRATE_DIR := apps/backend
GOLANGCI_LINT_VERSION := v2.12.2
POSTGRES_USER ?= wordsearch
POSTGRES_DEV_DB ?= wordsearch_dev

dev: ## Run the development stack with Air hot reload
	docker compose --profile dev up --build --watch --remove-orphans

prod: ## Build and run the production stack
	docker compose --profile prod up -d --build

down: ## Stop either Compose stack and remove its containers
	docker compose --profile dev --profile prod down --remove-orphans
	
db-up: ## Start local Postgres
	docker compose --profile dev up -d --wait --wait-timeout 30 postgres-dev

dev-up: db-up ## Start Postgres and create the Atlas development database if missing
	@docker compose --profile dev exec -T postgres-dev \
	  psql -U $(POSTGRES_USER) -tc \
	  "SELECT 1 FROM pg_database WHERE datname='$(POSTGRES_DEV_DB)'" | grep -q 1 \
	  || docker compose --profile dev exec -T postgres-dev \
	  psql -U $(POSTGRES_USER) -c "CREATE DATABASE $(POSTGRES_DEV_DB);"

db-down: ## Stop local Postgres
	docker compose --profile dev stop postgres-dev

psql: ## Open a psql shell against local postgres
	docker compose --profile dev exec postgres-dev psql -U $${POSTGRES_USER:-wordsearch} -d $${POSTGRES_DB:-wordsearch}

tidy: ## Sync go.mod/go.sum
	cd $(MIGRATE_DIR) && go mod tidy

fmt: ## Format all Go source
	cd $(MIGRATE_DIR) && gofmt -l -w .

vet: ## Run go vet
	cd $(MIGRATE_DIR) && go vet ./...

lint: ## Run golangci-lint (installs the pinned version if needed)
	@golangci-lint version 2>/dev/null | grep -q "$(GOLANGCI_LINT_VERSION)" || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	cd $(MIGRATE_DIR) && golangci-lint run ./...

# ── Tests ──────────────────────────────────────────────────

test:                                                                   ## Run all tests
	cd $(MIGRATE_DIR) && go test ./... -count=1

test-verbose:                                                           ## Run all tests with verbose output
	cd $(MIGRATE_DIR) && go test ./... -v -count=1

test-coverage:                                                          ## Run all tests with coverage report
	cd $(MIGRATE_DIR) && go test ./... -cover -count=1

# ── Database Migrations (Atlas) ──────────────────────────

migrate-diff: dev-up                                                    ## Generate a new migration from GORM models
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-diff NAME=<description>"; exit 1; fi
	cd $(MIGRATE_DIR) && atlas migrate diff $(NAME) --env local

migrate-up: dev-up                                                      ## Apply pending migrations (local)
	cd $(MIGRATE_DIR) && atlas migrate apply --env local

migrate-down: dev-up                                                    ## Rollback the last migration (local)
	cd $(MIGRATE_DIR) && atlas migrate down --env local 1

migrate-status: dev-up                                                  ## Show migration status (local)
	cd $(MIGRATE_DIR) && atlas migrate status --env local

migrate-prod-up:                                                        ## Apply pending migrations (production)
	cd $(MIGRATE_DIR) && atlas migrate apply --env production

migrate-prod-down:         
	@if [ "$(CONFIRM_PRODUCTION_DOWN)" != "1" ]; then \
	  echo "Set CONFIRM_PRODUCTION_DOWN=1 to roll back production"; exit 1; \
	fi                                             ## Rollback the last migration (production)
	cd $(MIGRATE_DIR) && atlas migrate down --env production 1

migrate-prod-status:                                                    ## Show migration status (production)
	cd $(MIGRATE_DIR) && atlas migrate status --env production
