# Mockbuster Movie API 🎬

A production-ready RESTful API for managing movie rentals, built with Go and PostgreSQL. This API provides comprehensive film management capabilities with advanced filtering, customer comments, and full CRUD operations.

![Go](https://img.shields.io/badge/Go-1.25+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)
![Swagger](https://img.shields.io/badge/Swagger-Documented-orange.svg)

### Future improvements
- **Enable SSL**
- **Add structured logger to the service**
- **Move shared library code to a shared repo, or to pkg directory to be used elsewhere**
- **Update build command to use Orchestrion for automated instrumentation (`go tool orchestrion build`)**
- **Use Resty or similar HTTP framework for structured request/response logging middleware**
- **Implement HTTP response codes properly**
- **Consolidate all error responses and add to swagger**
- **Add helm charts and k8s config for deployment**
- **Authentication and Authorization**
- **Add rate limiting**

## 🚀 Features

### Core Functionality
- **🎭 Film Management**: Complete CRUD operations for films with rich metadata
- **💬 Customer Comments**: Add and retrieve customer reviews and comments
- **🔍 Advanced Search**: Search films by title, category, actor, and rating
- **📄 Pagination**: Efficient pagination for large datasets with customizable limits
- **🏷️ Category Management**: Browse and filter by film categories
- **👥 Actor Information**: View cast information for each film

### Technical
- **📚 OpenAPI Documentation**: Auto-generated interactive API documentation
- **🗄️ Database Migrations**: Version-controlled schema management with Goose
- **📊 Monitoring Ready**: Structured logging and health checks, instrumentation ready by swapping build command using Orchestrion
- **🐳 Containerized**: Full Docker support with optimized multi-stage pipeline for deployment image vs build image

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.25+** - For modern Go features and tooling
- **Docker & Docker Compose** - For containerization and database setup
- **Earthfile (Optional)** - Replacement for Dockerfile+Makefile, for consistent CI and local development

## ⚡ Quick Start

Get up and running in minutes:

```bash
# Clone the repository
# Start everything with Docker Compose
make run

# Verify the API is running
curl http://localhost:8080/

# Access the interactive API documentation
open http://localhost:8080/swagger/
```

## 📁 Project Architecture

```
misal-patel-rxbenefits/
├── cmd/mockbuster/          # Application entry point
│   └── main.go              # Main application file
├── internal/                # Private application code
│   ├── database/            # Database connection & migrations
│   ├── handlers/            # HTTP request handlers
│   ├── models/              # Data structures & validation
│   ├── repository/          # Data access layer (Repository pattern)
│   ├── service/             # Business logic layer
│   └── util/                # Configuration and future utilities
├── migrations/              # 📦 Database migrations (Goose)
├── tests/                   # 🧪 Tests
│   ├── integration/         # End-to-end tests
│   └── unit/                # Unit tests
├── docs/                    # 📚 Generated API documentation
├── assets/                  # 🎨 Static assets
└── test/data/               # 🗃️ Database sample data
```

## 🔌 API Endpoints

### Films Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/films` | List films with filtering and pagination |
| `GET` | `/api/v1/films/{id}` | Get detailed film information |
| `GET` | `/api/v1/categories` | List all available categories |

### Comments System
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/films/{id}/comments` | Add a customer comment |
| `GET` | `/api/v1/films/{id}/comments` | Get all comments for a film |

### General
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Welcome message and API status |

## 📖 API Examples

### Get Films with Filtering
```bash
# Get first 10 films
curl "http://localhost:8080/api/v1/films?page=1&limit=10"

# Search by title
curl "http://localhost:8080/api/v1/films?title=Academy"

# Filter by rating
curl "http://localhost:8080/api/v1/films?rating=PG"

# Filter by category
curl "http://localhost:8080/api/v1/films?category=Action"

# Combine filters
curl "http://localhost:8080/api/v1/films?title=Academy&rating=PG&page=1&limit=5"
```

### Add a Comment
```bash
curl -X POST "http://localhost:8080/api/v1/films/1/comments" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "comment": "Excellent movie! Highly recommended."
  }'
```

### Get comments
```bash
curl "http://localhost:8080/api/v1/films/1/comments"
```

### Get Film Details
```bash
curl "http://localhost:8080/api/v1/films/1"
```

## 🗄️ Database Schema

The API uses a PostgreSQL database with the following key tables:

| Table | Description |
|-------|-------------|
| `film` | Core movie information (title, description, rating, etc.) |
| `category` | Film categories (Action, Drama, Comedy, etc.) |
| `actor` | Actor information |
| `film_actor` | Many-to-many relationship between films and actors |
| `film_category` | Many-to-many relationship between films and categories |
| `film_comments` | Customer comments and reviews |

### Database Migrations

The application uses **Goose** for database migrations, providing:
- ✅ **Version Control**: Track schema changes over time
- ✅ **Rollback Capability**: Revert changes if needed (not recommended)
- ✅ **Team Collaboration**: Consistent schema across environments
- ✅ **Production Safety**: Controlled database changes

**Migration Commands:**
```bash
# Check migration status
make migrate-status

# Run pending migrations
make migrate-up

# Rollback last migration
make migrate-down
```

## ⚙️ Configuration

The API is configured through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database host address |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `dvdrental` | Database name |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `password` | Database password |
| `PORT` | `8080` | API server port |

## 🧪 Testing

### Running Tests
```bash
# Run all tests with coverage
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run specific test package
go test -v ./internal/handlers

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📚 Documentation

### Interactive API Documentation
- **Swagger UI**: http://localhost:8080/swagger/
- **OpenAPI JSON**: http://localhost:8080/swagger/doc.json

### Local Documentation Generation
```bash
# Generate documentation
make docs

# View generated docs
open docs/index.html
```
## Earthly Support (Alternative Build System)

### Why Earthly?
- **Consistent Builds**: Same commands work locally and in CI
- **Better Caching**: Intelligent layer caching for faster builds
- **Parallel Execution**: Build steps run in parallel when possible
- **Docker Compatible**: Uses Docker under the hood

### Essential Earthly Commands
```bash
# Run with Docker Compose (equivalent to 'make run')
earthly +run

# Run all tests with coverage (equivalent to 'make test')
earthly +test
```

### Installation
```bash
# Install Earthly
curl -sSL https://earthly.dev/get-earthly | bash

# Or with Homebrew
brew install earthly
```

## Design Patterns & Architecture

### Architectural Patterns
- **Repository Pattern**: Clean separation of data access logic
- **Service Layer Pattern**: Business logic encapsulation
- **Dependency Injection**: Loose coupling between components
- **Middleware Pattern**: Cross-cutting concerns (CORS, logging)
- **Interface Segregation**: Small, focused interfaces

## 🔒 Security Features

### Input Validation
- **Request Validation**: All inputs validated using `go-playground/validator`
- **SQL Injection Prevention**: Parameterized queries throughout

### API Security
- **CORS Configuration**: Configurable cross-origin resource sharing
- **Input Sanitization**: Proper handling of user inputs

## 📊 Performance Optimizations

### Database Performance
- **Connection Pooling**: Efficient database connection management
- **Indexed Queries**: Optimized database queries with proper indexing
- **Pagination**: Large result sets are paginated to prevent memory issues

### Application Performance
- **Structured Logging**: Efficient logging with `log/slog`
- **Context Management**: Proper timeout and cancellation handling
- **Memory Management**: Efficient memory usage with Go's garbage collector

## 🛠️ Development Tools

### Code Quality
```bash
# Lint code
make lint

# Format code
go fmt ./...

# Run security checks
go vet ./...
```

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch
3. **Commit** your changes (Follow conventional commit pattern)
4. **Push** to the branch
5. **Open** a Pull Request

### Development Guidelines
- Follow Go coding standards and conventions
- Use Linter and fix all lints
- Write tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting