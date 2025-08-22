# Testing Documentation

This directory contains comprehensive tests for the Mockbuster Movie API, organized into unit tests and integration tests following Go testing best practices.

## Test Structure

```
tests/
├── unit/                     # Unit tests (no external dependencies)
│   ├── service/             # Service layer tests
│   │   ├── film_service_test.go
│   │   └── comment_service_test.go
│   └── handlers/            # Handler layer tests
│       └── film_handlers_test.go
├── integration/             # Integration tests (requires database)
│   └── api_test.go
├── config/                  # Test configuration utilities
│   └── test_config.go
└── README.md               # This file
```

## Test Types

### Unit Tests (`tests/unit/`)

Unit tests focus on testing individual components in isolation using mocks and stubs. They:

- **Service Tests**: Test business logic in the service layer
  - Input validation
  - Business rule enforcement
  - Error handling
  - Mocked repository dependencies

- **Handler Tests**: Test HTTP request/response handling
  - Request parsing
  - Response formatting
  - HTTP status codes
  - Mocked service dependencies

**Benefits:**
- Fast execution
- No external dependencies
- Easy to debug
- High coverage of edge cases

### Integration Tests (`tests/integration/`)

Integration tests verify the entire application stack working together with mocked dependencies:

- **API Tests**: End-to-end testing of REST endpoints
  - Mocked database interactions
  - Full request/response cycle
  - Mocked data persistence
  - Multi-component workflows

**Benefits:**
- Tests real-world scenarios
- Catches integration issues
- Fast execution (no database required)
- Ensures API contracts work
- Reliable and repeatable

## Running Tests

### Prerequisites

- **For Unit Tests**: No external dependencies required
- **For Integration Tests**: No external dependencies required (uses mocked database)
- **For Database Tests**: PostgreSQL test database on `localhost:5556`

### Commands

```bash
# Run all tests
make test

# Run only unit tests (fast, no dependencies)
make test-unit

# Run only integration tests (no database required)
make test-integration

# Generate coverage reports
# - coverage.html (all tests)
# - coverage-unit.html (unit tests only)
# - coverage-integration.html (integration tests only)
```

### Test Database Setup

For integration tests, you need a PostgreSQL test database:

```bash
# Start test database with Docker
docker run -d \
  --name dvdrental-postgres-test \
-e POSTGRES_DB=dvdrental_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5556:5432 \
  postgres:15

# Or use environment variables to configure
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5556
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=password
export TEST_DB_NAME=dvdrental_test
```

## Test Configuration

Tests use environment variables for configuration with sensible defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `TEST_DB_HOST` | `localhost` | Test database host |
| `TEST_DB_PORT` | `5556` | Test database port |
| `TEST_DB_USER` | `postgres` | Test database user |
| `TEST_DB_PASSWORD` | `password` | Test database password |
| `TEST_DB_NAME` | `dvdrental_test` | Test database name |

## Testing Patterns

### Dependency Injection Testing

Following the [dependency injection article](https://medium.com/avenue-tech/dependency-injection-in-go-35293ef7b6), our tests use:

- **Interface-based mocking**: Services depend on repository interfaces
- **Constructor injection**: Dependencies injected through constructors
- **Testify mocks**: `github.com/stretchr/testify/mock` for mock implementations

### Test Naming Convention

```go
func TestServiceName_MethodName(t *testing.T) {
    tests := []struct {
        name           string
        input          InputType
        expectedResult ExpectedType
        expectedError  string
    }{
        {
            name: "descriptive test case name",
            // test case data
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Coverage Goals

- **Unit Tests**: Aim for >90% coverage of business logic
- **Integration Tests**: Cover all major API workflows
- **Combined**: >80% overall coverage

## Continuous Integration

These tests are designed to run in CI/CD pipelines:

- **Unit tests**: Run on every commit (fast feedback)
- **Integration tests**: Run on merge requests (thorough validation)
- **Coverage reports**: Generated and uploaded as artifacts

## Adding New Tests

### For New Features

1. **Start with unit tests**: Test business logic in isolation
2. **Add integration tests**: Test the full feature workflow
3. **Update this README**: Document any new test patterns or requirements

### Test File Organization

- Place unit tests in `tests/unit/` matching the source structure
- Place integration tests in `tests/integration/` by feature area
- Use descriptive file names ending with `_test.go`

## Best Practices

1. **Test Isolation**: Each test should be independent
2. **Clear Assertions**: Use descriptive error messages
3. **Mock Management**: Clean up mocks after each test
4. **Test Data**: Use minimal, focused test data
5. **Error Cases**: Test both success and failure scenarios

## Troubleshooting

### Common Issues

- **"Database connection failed"**: Ensure test database is running on port 5556
- **"Mock expectations not met"**: Check mock setup and cleanup
- **"Import cycle"**: Ensure test packages don't import main packages circularly

### Debug Tips

- Use `go test -v` for verbose output
- Check coverage reports to find untested code paths
- Run individual tests: `go test -run TestSpecificTest`
