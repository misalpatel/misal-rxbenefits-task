VERSION 0.7

# Base target for Go builds
build-base:
    FROM golang:1.25-alpine
    RUN apk add --no-cache git ca-certificates tzdata && \
        update-ca-certificates
    
    # Set Go environment variables
    ENV GO111MODULE=on
    ENV GOOS=linux
    ENV GOARCH=amd64
    ENV CGO_ENABLED=0
    
    WORKDIR /app
    
    # Copy go mod files for dependency caching
    COPY go.mod go.sum ./
    
    # Download dependencies
    RUN go mod download -x && \
        go mod verify
    
    # Copy source code
    COPY . .
    
    # Generate Swagger documentation
    RUN go tool swag init -g cmd/mockbuster/main.go -o docs
    
    # Build the application with optimizations
    RUN go build -ldflags="-w -s" -o mockbuster-api cmd/mockbuster/main.go

# Run target - equivalent to 'make run'
run:
    FROM +build-base
    RUN apk add --no-cache docker-compose
    COPY docker-compose.yml .
    COPY test/data ./test/data
    EXPOSE 8080
    ENV DB_HOST=postgres
    ENV DB_PORT=5432
    ENV DB_NAME=dvdrental
    ENV DB_USER=postgres
    ENV DB_PASSWORD=password
    CMD ["docker-compose", "up", "--build"]

# Test target - equivalent to 'make test'
test:
    FROM +build-base
    RUN go test -v -coverprofile=coverage.out -coverpkg=$(go list ./... | grep -v -E '/(docs|assets|test)' | grep -v '/docs$$' | tr '\n' ',' | sed 's/,$//') ./...
    RUN go tool cover -func=coverage.out
    RUN go tool cover -html=coverage.out -o coverage.html
    SAVE ARTIFACT coverage.html AS LOCAL coverage.html

# Unit tests only - equivalent to 'make test-unit'
test-unit:
    FROM +build-base
    RUN go test -v -coverprofile=coverage-unit.out -coverpkg=$(go list ./... | grep -v -E '/(docs|assets|test)' | grep -v '/docs$$' | tr '\n' ',' | sed 's/,$//') ./tests/unit/...
    RUN go tool cover -func=coverage-unit.out
    RUN go tool cover -html=coverage-unit.out -o coverage-unit.html
    SAVE ARTIFACT coverage-unit.html AS LOCAL coverage-unit.html

# Database migrations - equivalent to 'make migrate-up'
migrate-up:
    FROM +build-base
    RUN go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "host=localhost port=5555 user=postgres password=password dbname=dvdrental sslmode=disable" up

# Database migrations down - equivalent to 'make migrate-down'
migrate-down:
    FROM +build-base
    RUN go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "host=localhost port=5555 user=postgres password=password dbname=dvdrental sslmode=disable" down

# Database migration status - equivalent to 'make migrate-status'
migrate-status:
    FROM +build-base
    RUN go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "host=localhost port=5555 user=postgres password=password dbname=dvdrental sslmode=disable" status
