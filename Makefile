.PHONY: dev prod down db-up db-down psql tidy fmt vet lint test test-verbose test-coverage

MIGRATE_DIR := apps/backend
GOLANGCI_LINT_VERSION := v2.12.2

dev: ## Run the development stack with Air hot reload
	docker compose --profile dev up --build --watch --remove-orphans

prod: ## Build and run the production stack
	docker compose --profile prod up -d --build

down: ## Stop either Compose stack and remove its containers
	docker compose --profile dev --profile prod down --remove-orphans
	
db-up: ## Start local Postgres
	docker compose --profile dev up -d --wait --wait-timeout 30 postgres-dev

db-down: ## Stop local Postgres
	docker compose --profile dev stop postgres-dev

psql: ## Open a psql shell against local postgres
	docker compose --profile dev exec postgres-dev psql -U wordsearch -d wordsearch

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
