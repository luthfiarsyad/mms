# Testing Guide for MMS (Money Management System)

This directory contains comprehensive unit tests for the MMS application, covering user authentication flow and transaction recording flow.

## Test Structure

### Test Files

- **`user_service_test.go`** - Tests for user authentication and registration
- **`transaction_usecase_test.go`** - Tests for transaction management operations
- **`test_setup.go`** - Test utilities and database setup helpers

### Test Coverage

#### User Authentication Flow Tests
- ✅ User registration with valid data
- ✅ User registration with duplicate email (should fail)
- ✅ User authentication with correct credentials
- ✅ User authentication with wrong email (should fail)
- ✅ User authentication with wrong password (should fail)
- ✅ AuthUsecase registration flow
- ✅ AuthUsecase login flow with PASETO token generation

#### Transaction Recording Flow Tests
- ✅ Create valid income transaction
- ✅ Create valid expense transaction
- ✅ Create transaction with invalid amount (should fail)
- ✅ Create transaction with invalid type (should fail)
- ✅ Create transaction with invalid user ID (should fail)
- ✅ Get transaction by ID
- ✅ Get transactions by user ID
- ✅ Update existing transaction
- ✅ Delete transaction
- ✅ TransactionUsecase operations

## Running Tests

### Prerequisites

1. **MySQL Server**: Make sure MySQL is running on your system
2. **Test Database**: Tests will automatically create a test database named `mms_test`
3. **Go Dependencies**: Ensure all required Go packages are installed

### Environment Variables (Optional)

You can customize the test database configuration using these environment variables:

```bash
export TEST_DB_HOST=127.0.0.1
export TEST_DB_PORT=3306
export TEST_DB_USER=root
export TEST_DB_PASSWORD=
export TEST_DB_NAME=mms_test
```

### Running All Tests

```bash
go test ./internal/test/...
```

### Running Specific Test Files

```bash
# Run user authentication tests
go test ./internal/test/ -run TestUserService

# Run transaction tests
go test ./internal/test/ -run TestTransactionService

# Run usecase tests
go test ./internal/test/ -run TestAuthUsecase
go test ./internal/test/ -run TestTransactionUsecase
```

### Running Tests with Verbose Output

```bash
go test ./internal/test/... -v
```

### Running Tests with Coverage

```bash
go test ./internal/test/... -cover
```

### Running Tests with Coverage Report

```bash
go test ./internal/test/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Database Setup

The test suite automatically:

1. **Creates Test Database**: Creates `mms_test` database if it doesn't exist
2. **Runs Migrations**: Sets up the required database schema
3. **Cleans Data**: Cleans up test data between test runs
4. **Tears Down**: Cleans up resources after tests complete

### Manual Database Setup (Optional)

If you prefer to set up the test database manually:

```sql
CREATE DATABASE mms_test;

USE mms_test;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    description VARCHAR(500) NOT NULL,
    type ENUM('income', 'expense') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
);
```

## Test Utilities

### TestHelper Class

The `TestHelper` class provides convenient methods for test setup:

```go
helper := NewTestHelper(t)
defer helper.Cleanup()

// Create test user
userID := helper.CreateTestUser("Test User", "test@example.com", "password")

// Create test transaction
txID := helper.CreateTestTransaction(userID, 100.0, "Test transaction", "income")
```

### Mock Services

Tests use mock implementations for external dependencies:

- **MockPasetoService**: Mocks PASETO token generation and verification
- **Test Database**: Uses isolated test database to avoid affecting production data

## Best Practices

### Test Isolation

- Each test runs in isolation with clean database state
- Test data is cleaned up after each test run
- Tests don't depend on each other's state

### Test Data Management

- Use descriptive test data names
- Create realistic test scenarios
- Test both success and failure cases

### Error Handling

- Test error conditions and edge cases
- Verify proper error messages and types
- Test validation logic

### Assertions

- Use specific assertions for better error messages
- Test all important fields and properties
- Verify both positive and negative conditions

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Ensure MySQL server is running
   - Check connection parameters in environment variables
   - Verify MySQL user permissions

2. **Permission Errors**
   - Ensure MySQL user has CREATE DATABASE privileges
   - Check file permissions for test directories

3. **Port Conflicts**
   - Make sure MySQL is running on the expected port
   - Update TEST_DB_PORT environment variable if needed

4. **Test Failures**
   - Check test output for specific error messages
   - Verify test database schema is up to date
   - Ensure all dependencies are properly installed

### Debug Mode

Run tests with debug output:

```bash
go test ./internal/test/... -v -args -test.v
```

## Contributing

When adding new tests:

1. Follow the existing test structure and naming conventions
2. Use the provided test utilities and helpers
3. Ensure tests are isolated and don't depend on each other
4. Test both success and failure scenarios
5. Add appropriate cleanup code
6. Update this README if adding new test categories

## Performance Considerations

- Tests use in-memory operations where possible
- Database connections are reused within test runs
- Test data is minimized to reduce execution time
- Parallel test execution is supported where safe

## Continuous Integration

These tests are designed to run in CI/CD environments:

- Uses environment variables for configuration
- Handles missing dependencies gracefully
- Provides clear error messages for debugging
- Supports parallel execution