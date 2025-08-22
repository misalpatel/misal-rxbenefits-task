# Variables
APP_NAME = mockbuster-api
IMAGE_NAME = mockbuster-api
IMAGE_TAG = latest
DB_NAME = mockbuster-postgres
DB_PORT = 5555
API_PORT = 8080

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build        - Build Docker image"
	@echo "  make run          - Run with Docker Compose"
	@echo "  make test         - Run all tests with coverage (excludes docs, assets, tests dirs)"
	@echo "  make test-unit    - Run unit tests only"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make lint         - Lint code"
	@echo "  make docs         - Generate OpenAPI docs"
	@echo "  make migrate-up   - Run database migrations up"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make migrate-status - Show migration status"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make cleanup      - Full cleanup (containers + images)"



# Deps
.PHONY: deps
deps:
	go mod download
	go mod tidy
	@echo "Dependencies installed!"


# Build Docker image
.PHONY: build
build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

# Run the application locally
.PHONY: run
run: build
	docker-compose up --build

# Run all tests with coverage (excluding docs, assets, test directories)
.PHONY: test
test:
	@echo "Running all tests with coverage..."
	go test -v -coverprofile=coverage.out -coverpkg=$$(go list ./... | grep -v -E '/(docs|assets|test)' | grep -v '/docs$$' | tr '\n' ',' | sed 's/,$$//') ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@rm -f coverage.out

# Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "Running unit tests with coverage..."
	go test -v -coverprofile=coverage-unit.out -coverpkg=$$(go list ./... | grep -v -E '/(docs|assets|test)' | grep -v '/docs$$' | tr '\n' ',' | sed 's/,$$//') ./tests/unit/...
	go tool cover -func=coverage-unit.out
	go tool cover -html=coverage-unit.out -o coverage-unit.html
	@echo "Unit test coverage report generated: coverage-unit.html"
	@rm -f coverage-unit.out

# Run integration tests only (uses mocked database)
.PHONY: test-integration
test-integration:
	@echo "Running integration tests with mocked database..."
	go test -v ./tests/integration/...
	@echo "Integration tests completed"

# Lint code
.PHONY: lint
lint:
	go tool golangci-lint run

# Database migrations
.PHONY: migrate-up
migrate-up:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "host=localhost port=5555 user=postgres password=password dbname=dvdrental sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "host=localhost port=5555 user=postgres password=password dbname=dvdrental sslmode=disable" down

.PHONY: migrate-status
migrate-status:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "host=localhost port=5555 user=postgres password=password dbname=dvdrental sslmode=disable" status

# Generate OpenAPI docs
.PHONY: docs
docs: deps
	go tool swag init -g cmd/mockbuster/main.go -o docs

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(APP_NAME)
	rm -f coverage*.out coverage*.html
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) 2>/dev/null || true


# Full cleanup
.PHONY: cleanup
cleanup: clean db-clean
	docker-compose down -v 2>/dev/null || true
	@echo "Cleanup complete!"
