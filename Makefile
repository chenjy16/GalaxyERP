# Makefile for GalaxyERP

# Variables
APP_NAME = galaxyerp
BINARY_NAME = galaxyerp
MAIN_FILE = cmd/server/main.go
MIGRATE_FILE = cmd/migrate/main.go

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Default target
all: build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_FILE)

# Build the migration tool
build-migrate:
	$(GOBUILD) -o migrate -v $(MIGRATE_FILE)

# Run the application
run:
	GALAXYERP_ENV=dev $(GOCMD) run $(MAIN_FILE)

# Run the application in test environment
run-test:
	GALAXYERP_ENV=test $(GOCMD) run $(MAIN_FILE)

# Run the application in production environment
run-prod:
	GALAXYERP_ENV=prod $(GOCMD) run $(MAIN_FILE)

# Run the migration tool
migrate:
	GALAXYERP_ENV=dev $(GOCMD) run $(MIGRATE_FILE)

# Run the migration tool in test environment
migrate-test:
	GALAXYERP_ENV=test $(GOCMD) run $(MIGRATE_FILE)

# Run the migration tool in production environment
migrate-prod:
	GALAXYERP_ENV=prod $(GOCMD) run $(MIGRATE_FILE)

# Run tests
test:
	GALAXYERP_ENV=dev $(GOTEST) -v ./...

# Run tests in test environment
test-test:
	GALAXYERP_ENV=test $(GOTEST) -v ./...

# Run tests in production environment
test-prod:
	GALAXYERP_ENV=prod $(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f migrate

# Install dependencies
deps:
	$(GOGET) ./...

# Run tests with coverage
coverage:
	GALAXYERP_ENV=dev $(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Run development environment
dev:
	./scripts/run-dev.sh

# Start development database
dev-db-start:
	./scripts/start-dev-db.sh

# Stop development database
dev-db-stop:
	./scripts/stop-dev-db.sh

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  build        - Build the application"
	@echo "  build-migrate - Build the migration tool"
	@echo "  run          - Run the application (dev environment)"
	@echo "  run-test     - Run the application (test environment)"
	@echo "  run-prod     - Run the application (prod environment)"
	@echo "  migrate      - Run the migration tool (dev environment)"
	@echo "  migrate-test - Run the migration tool (test environment)"
	@echo "  migrate-prod - Run the migration tool (prod environment)"
	@echo "  test         - Run tests (dev environment)"
	@echo "  test-test    - Run tests (test environment)"
	@echo "  test-prod    - Run tests (prod environment)"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  coverage     - Run tests with coverage"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  dev          - Run development environment"
	@echo "  dev-db-start - Start development database"
	@echo "  dev-db-stop  - Stop development database"
	@echo "  help         - Show this help"

.PHONY: all build build-migrate run run-test run-prod migrate migrate-test migrate-prod test test-test test-prod clean deps coverage fmt vet dev dev-db-start dev-db-stop help