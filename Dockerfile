FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x && \
    go mod verify

COPY . .

# Generate Swagger documentation
RUN go tool swag init -g cmd/mockbuster/main.go -o docs

RUN go build -ldflags="-w -s" -o mockbuster-api cmd/mockbuster/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/mockbuster-api .

# Copy migrations directory
COPY --from=builder /app/migrations ./migrations

# Copy docs directory for Swagger documentation
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["./mockbuster-api"]
