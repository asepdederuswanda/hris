# HRIS Platform Makefile
# ============ Commands ============
# make build       - Build Go binary
# make run         - Run server locally
# make test        - Run all tests
# make lint        - Run linter
# make docker      - Build Docker image
# make clean       - Clean build artifacts
# make migrate     - Run platform migrations
# make seed        - Run seeders
# make coverage    - Run tests with coverage

.PHONY: build build-installer run run-hot test test-verbose test-short test-pkg test-employee test-organization coverage cover-view cover-func lint vet migrate seed docker docker-compose-up docker-compose-down tidy clean help

BUILD_DIR := ./bin
BINARY := $(BUILD_DIR)/server

# ============ Build ============

build:
	@echo "Building HRIS Platform..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) ./cmd/server
	@echo "Build complete: $(BINARY)"

build-installer:
	@echo "Building installer..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/installer ./cmd/installer
	@echo "Build complete: $(BUILD_DIR)/installer"

# ============ Run ============

run:
	@echo "Starting server..."
	@go run ./cmd/server --config ./config/config.yaml

run-hot:
	@echo "Starting server with hot reload..."
	@which air || go install github.com/air-verse/air@latest
	@air

# ============ Test ============

test:
	@echo "Running all tests..."
	@go test -v -race -count=1 ./...

test-verbose:
	@echo "Running all tests (verbose)..."
	@go test -v -race -count=1 ./... 2>&1 | tee test-output.log

test-short:
	@echo "Running short tests (non-integration)..."
	@go test -short -count=1 ./...

test-pkg:
	@echo "Running tests for package: $(pkg)"
	@go test -v -race -count=1 ./internal/$(pkg)/...

test-employee:
	@echo "Running employee module tests..."
	@go test -v -count=1 ./internal/modules/employee/...

test-organization:
	@echo "Running organization module tests..."
	@go test -v -count=1 ./internal/modules/organization/...

coverage:
	@echo "Running tests with coverage..."
	@go test -v -count=1 -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

cover-view:
	@echo "Opening coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

cover-func:
	@echo "Coverage by function:"
	@go tool cover -func=coverage.out

# ============ Lint ============

lint:
	@echo "Running linter..."
	@which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.0
	@golangci-lint run ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

# ============ Database ============

migrate:
	@echo "Running platform migrations..."
	@go run ./cmd/server migrate --config ./config/config.yaml

seed:
	@echo "Running seeders..."
	@go run ./cmd/server seed --config ./config/config.yaml

# ============ Docker ============

docker:
	@echo "Building Docker image..."
	@docker build -t hris-platform:latest -f docker/Dockerfile .

docker-compose-up:
	@echo "Starting Docker Compose environment..."
	@docker compose -f ../docker/docker-compose.yml up -d

docker-compose-down:
	@echo "Stopping Docker Compose environment..."
	@docker compose -f ../docker/docker-compose.yml down

# ============ Utilities ============

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html test-output.log
	@go clean

tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify

help:
	@echo "HRIS Platform Makefile"
	@echo "======================"
	@echo "=== Build ==="
	@echo "build          - Build Go binary"
	@echo "build-installer - Build CLI installer"
	@echo ""
	@echo "=== Run ==="
	@echo "run            - Run server locally"
	@echo "run-hot        - Run with hot reload (air)"
	@echo ""
	@echo "=== Test ==="
	@echo "test           - Run all tests (verbose + race + count=1)"
	@echo "test-verbose   - Run all tests, output to test-output.log"
	@echo "test-short     - Run short tests (non-integration)"
	@echo "test-pkg pkg=   - Run tests for a specific package (e.g. pkg=modules/employee)"
	@echo "test-employee  - Run employee module tests"
	@echo "test-organization - Run organization module tests"
	@echo "coverage       - Run tests with coverage report (HTML)"
	@echo "cover-view     - Open existing coverage report"
	@echo "cover-func     - Show coverage by function"
	@echo ""
	@echo "=== Lint ==="
	@echo "lint           - Run linter (golangci-lint)"
	@echo "vet            - Run go vet"
	@echo ""
	@echo "=== Database ==="
	@echo "migrate        - Run platform migrations"
	@echo "seed           - Run seeders"
	@echo ""
	@echo "=== Docker ==="
	@echo "docker         - Build Docker image"
	@echo "docker-compose-up   - Start Docker Compose environment"
	@echo "docker-compose-down - Stop Docker Compose environment"
	@echo ""
	@echo "=== Utilities ==="
	@echo "tidy           - Tidy and verify dependencies"
	@echo "clean          - Clean build artifacts"
	@echo "help           - Show this help"