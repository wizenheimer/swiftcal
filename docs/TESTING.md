# Testing Guide

## Test Structure

SwiftCal uses a comprehensive testing strategy with the following test types:

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test service interactions
- **API Tests**: Test HTTP endpoints
- **Database Tests**: Test data persistence

## Running Tests

### All Tests

```bash
make test
```

### Unit Tests Only

```bash
go test ./internal/...
```

### Integration Tests

```bash
go test -tags=integration ./...
```

### API Tests

```bash
go test ./internal/handlers/...
```

## Test Coverage

Generate coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Test Data

### Mock Data

Test data is stored in `testdata/` directory:

- `users.json` - Sample user data
- `events.json` - Sample calendar events
- `emails.json` - Sample email content

### Database Fixtures

Use `scripts/test_fixtures.sql` to populate test database.

## Writing Tests

### Unit Test Example

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    service := NewUserService(mockDB)
    user := &models.User{
        Email: "test@example.com",
        Name:  "Test User",
    }

    // Act
    result, err := service.CreateUser(user)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, user.Email, result.Email)
}
```

### Integration Test Example

```go
func TestCalendarService_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    service := NewCalendarService(db)

    // Test calendar operations
    // ...
}
```

## Performance Testing

### Load Testing

Use `scripts/load_test.sh` to run load tests:

```bash
./scripts/load_test.sh
```

### Benchmark Tests

Run benchmarks:

```bash
go test -bench=. ./...
```

## Continuous Integration

Tests are automatically run on:

- Pull requests
- Main branch commits
- Release tags

### CI Pipeline

1. Lint code
2. Run unit tests
3. Run integration tests
4. Generate coverage report
5. Build Docker image
6. Run end-to-end tests
