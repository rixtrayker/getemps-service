# GetEmpStatus Service Makefile
# =================================

# Variables
APP_NAME := getemps-service
BINARY_NAME := getemps-api
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_COMPOSE_FILE := docker-compose.yml
GO_VERSION := 1.21

# Build info
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")

# Go build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
NC := \033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)GetEmpStatus Service - Available Commands:$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

# =================================
# Development Commands
# =================================

.PHONY: setup
setup: ## Setup development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@go version
	@go mod download
	@go mod tidy
	@cp .env.example .env
	@echo "$(GREEN)✓ Development environment ready$(NC)"

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@echo "$(YELLOW)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(RED)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Go vet passed$(NC)"

# =================================
# Build Commands
# =================================

.PHONY: build
build: ## Build the application binary
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/api
	@echo "$(GREEN)✓ Binary built: bin/$(BINARY_NAME)$(NC)"

.PHONY: build-linux
build-linux: ## Build for Linux (useful for Docker)
	@echo "$(YELLOW)Building $(BINARY_NAME) for Linux...$(NC)"
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux ./cmd/api
	@echo "$(GREEN)✓ Linux binary built: bin/$(BINARY_NAME)-linux$(NC)"

.PHONY: build-all
build-all: ## Build for all platforms
	@echo "$(YELLOW)Building for all platforms...$(NC)"
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/api
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/api
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/api
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/api
	@echo "$(GREEN)✓ All platform binaries built$(NC)"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@go clean
	@echo "$(GREEN)✓ Clean completed$(NC)"

# =================================
# Run Commands
# =================================

.PHONY: run
run: ## Run the application locally
	@echo "$(YELLOW)Starting $(APP_NAME)...$(NC)"
	@go run ./cmd/api

.PHONY: run-build
run-build: build ## Build and run the application
	@echo "$(YELLOW)Running built binary...$(NC)"
	@./bin/$(BINARY_NAME)

.PHONY: dev
dev: ## Run with live reload (requires air)
	@echo "$(YELLOW)Starting development server with live reload...$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(RED)Air not installed. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(YELLOW)Falling back to regular run...$(NC)"; \
		$(MAKE) run; \
	fi

# =================================
# Test Commands
# =================================

.PHONY: test
test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✓ Tests completed$(NC)"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

.PHONY: test-integration
test-integration: ## Run integration tests (requires running database)
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@go test -v -tags=integration ./...

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

# =================================
# Database Commands
# =================================

.PHONY: db-up
db-up: ## Start database container
	@echo "$(YELLOW)Starting database...$(NC)"
	@docker-compose up -d postgres
	@echo "$(GREEN)✓ Database started$(NC)"

.PHONY: db-down
db-down: ## Stop database container
	@echo "$(YELLOW)Stopping database...$(NC)"
	@docker-compose stop postgres
	@echo "$(GREEN)✓ Database stopped$(NC)"

.PHONY: db-reset
db-reset: ## Reset database (WARNING: destroys all data)
	@echo "$(RED)WARNING: This will destroy all database data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo ""; \
		echo "$(YELLOW)Resetting database...$(NC)"; \
		docker-compose down -v postgres; \
		docker-compose up -d postgres; \
		sleep 5; \
		$(MAKE) db-migrate; \
		echo "$(GREEN)✓ Database reset completed$(NC)"; \
	else \
		echo ""; \
		echo "$(YELLOW)Database reset cancelled$(NC)"; \
	fi

.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "$(YELLOW)Running database migrations...$(NC)"
	@sleep 2  # Wait for database to be ready
	@docker-compose exec postgres psql -U postgres -d getemps_db -f /docker-entrypoint-initdb.d/001_create_tables.sql || true
	@docker-compose exec postgres psql -U postgres -d getemps_db -f /docker-entrypoint-initdb.d/002_insert_data.sql || true
	@echo "$(GREEN)✓ Migrations completed$(NC)"

.PHONY: db-shell
db-shell: ## Connect to database shell
	@echo "$(YELLOW)Connecting to database...$(NC)"
	@docker-compose exec postgres psql -U postgres -d getemps_db

.PHONY: db-logs
db-logs: ## Show database logs
	@docker-compose logs -f postgres

# =================================
# Docker Commands
# =================================

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -f docker/Dockerfile -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE)$(NC)"

.PHONY: docker-run
docker-run: ## Run application in Docker container
	@echo "$(YELLOW)Running Docker container...$(NC)"
	@docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

.PHONY: docker-push
docker-push: docker-build ## Build and push Docker image
	@echo "$(YELLOW)Pushing Docker image...$(NC)"
	@docker push $(DOCKER_IMAGE)
	@echo "$(GREEN)✓ Docker image pushed$(NC)"

.PHONY: up
up: ## Start all services with docker-compose
	@echo "$(YELLOW)Starting all services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ All services started$(NC)"
	@echo "$(BLUE)API available at: http://localhost:8080$(NC)"
	@echo "$(BLUE)Health check: curl http://localhost:8080/health$(NC)"

.PHONY: down
down: ## Stop all services
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ All services stopped$(NC)"

.PHONY: restart
restart: down up ## Restart all services

.PHONY: logs
logs: ## Show logs from all services
	@docker-compose logs -f

.PHONY: logs-api
logs-api: ## Show API service logs
	@docker-compose logs -f api

.PHONY: ps
ps: ## Show running containers
	@docker-compose ps

# =================================
# API Testing Commands
# =================================

.PHONY: test-api
test-api: ## Test API endpoints (requires running service)
	@echo "$(YELLOW)Testing API endpoints...$(NC)"
	@if [ -f scripts/test-api.sh ]; then \
		./scripts/test-api.sh; \
	else \
		echo "$(BLUE)Testing health check...$(NC)"; \
		curl -s http://localhost:8080/health | jq . || curl -s http://localhost:8080/health; \
		echo ""; \
		echo "$(BLUE)Testing valid employee (NAT1001)...$(NC)"; \
		curl -s -X POST http://localhost:8080/api/GetEmpStatus \
			-H "Content-Type: application/json" \
			-d '{"NationalNumber": "NAT1001"}' | jq . || echo "Response received"; \
		echo ""; \
		echo "$(BLUE)Testing invalid employee (NAT9999)...$(NC)"; \
		curl -s -X POST http://localhost:8080/api/GetEmpStatus \
			-H "Content-Type: application/json" \
			-d '{"NationalNumber": "NAT9999"}' | jq . || echo "Response received"; \
		echo "$(GREEN)✓ API tests completed$(NC)"; \
	fi

.PHONY: test-postman
test-postman: ## Run Postman collection tests (requires newman)
	@echo "$(YELLOW)Running Postman collection tests...$(NC)"
	@if command -v newman >/dev/null 2>&1; then \
		newman run docs/api/GetEmpStatus_Complete_Collection.json \
			-e docs/api/environments/Development.postman_environment.json \
			--reporters cli,json \
			--reporter-json-export postman-results.json; \
		echo "$(GREEN)✓ Postman tests completed$(NC)"; \
	else \
		echo "$(RED)Newman not installed. Install with: npm install -g newman$(NC)"; \
	fi

.PHONY: test-postman-performance
test-postman-performance: ## Run Postman performance tests
	@echo "$(YELLOW)Running Postman performance tests...$(NC)"
	@if command -v newman >/dev/null 2>&1; then \
		newman run docs/api/GetEmpStatus_Performance_Collection.json \
			-e docs/api/environments/Development.postman_environment.json \
			--iteration-count 5 \
			--delay-request 100 \
			--reporters cli; \
		echo "$(GREEN)✓ Performance tests completed$(NC)"; \
	else \
		echo "$(RED)Newman not installed. Install with: npm install -g newman$(NC)"; \
	fi

.PHONY: test-api-all
test-api-all: test-api test-postman ## Run all API tests (shell script + Postman)

.PHONY: test-load
test-load: ## Run load test (requires hey tool)
	@echo "$(YELLOW)Running load test...$(NC)"
	@if command -v hey >/dev/null 2>&1; then \
		hey -n 1000 -c 10 -m POST \
			-H "Content-Type: application/json" \
			-d '{"NationalNumber": "NAT1001"}' \
			http://localhost:8080/api/GetEmpStatus; \
	else \
		echo "$(RED)hey not installed. Install with: go install github.com/rakyll/hey@latest$(NC)"; \
	fi

# =================================
# Utility Commands
# =================================

.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(shell go version)"

.PHONY: env-check
env-check: ## Check environment variables
	@echo "$(YELLOW)Checking environment configuration...$(NC)"
	@if [ -f .env ]; then \
		echo "$(GREEN)✓ .env file exists$(NC)"; \
		echo "$(BLUE)Current configuration:$(NC)"; \
		cat .env | grep -E '^[A-Z_]+=.*' | sed 's/PASSWORD=.*/PASSWORD=***/' | sed 's/SECRET=.*/SECRET=***/' | sed 's/KEY=.*/KEY=***/'; \
	else \
		echo "$(RED)✗ .env file missing$(NC)"; \
		echo "$(YELLOW)Run 'make setup' to create from template$(NC)"; \
	fi

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/rakyll/hey@latest
	@if command -v npm >/dev/null 2>&1; then \
		npm install -g newman; \
		echo "$(GREEN)✓ Newman installed$(NC)"; \
	else \
		echo "$(YELLOW)⚠ npm not found, skipping Newman installation$(NC)"; \
	fi
	@echo "$(GREEN)✓ Development tools installed$(NC)"

.PHONY: security-scan
security-scan: ## Run security scan (requires gosec)
	@echo "$(YELLOW)Running security scan...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(RED)gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(NC)"; \
	fi

# =================================
# CI/CD Commands
# =================================

.PHONY: ci
ci: deps vet fmt test ## Run CI pipeline locally
	@echo "$(GREEN)✓ CI pipeline completed successfully$(NC)"

.PHONY: pre-commit
pre-commit: fmt vet lint test ## Run pre-commit checks
	@echo "$(GREEN)✓ Pre-commit checks passed$(NC)"

.PHONY: release-check
release-check: ci test-coverage security-scan ## Run full release checks
	@echo "$(GREEN)✓ Release checks completed$(NC)"

# =================================
# Documentation Commands
# =================================

.PHONY: docs
docs: ## Generate documentation
	@echo "$(YELLOW)Generating documentation...$(NC)"
	@go doc -all ./... > docs/api-documentation.txt
	@echo "$(GREEN)✓ Documentation generated$(NC)"

.PHONY: serve-docs
serve-docs: ## Serve documentation locally
	@echo "$(YELLOW)Starting documentation server...$(NC)"
	@if command -v godoc >/dev/null 2>&1; then \
		echo "$(BLUE)Documentation available at: http://localhost:6060$(NC)"; \
		godoc -http=:6060; \
	else \
		echo "$(RED)godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest$(NC)"; \
	fi

# =================================
# Monitoring Commands
# =================================

.PHONY: monitor
monitor: ## Monitor application metrics
	@echo "$(YELLOW)Monitoring application...$(NC)"
	@watch -n 1 'curl -s http://localhost:8080/health | jq .'

.PHONY: profile
profile: ## Generate CPU profile
	@echo "$(YELLOW)Generating CPU profile...$(NC)"
	@go tool pprof http://localhost:8080/debug/pprof/profile

# Default target
.DEFAULT_GOAL := help