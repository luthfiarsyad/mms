package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/luthfiarsyad/mms/config"
	_ "github.com/go-sql-driver/mysql"
)

// TestConfig holds configuration for testing
type TestConfig struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
}

// GetTestConfig returns test configuration from environment variables or defaults
func GetTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseHost:     getEnv("TEST_DB_HOST", "127.0.0.1"),
		DatabasePort:     getEnvInt("TEST_DB_PORT", 3306),
		DatabaseUser:     getEnv("TEST_DB_USER", "root"),
		DatabasePassword: getEnv("TEST_DB_PASSWORD", ""),
		DatabaseName:     getEnv("TEST_DB_NAME", "mms_test"),
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as int or returns default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue := parseInt(value); intValue != 0 {
			return intValue
		}
	}
	return defaultValue
}

// parseInt parses string to int (simple implementation)
func parseInt(s string) int {
	var result int
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		} else {
			return 0
		}
	}
	return result
}

// SetupTestDatabase creates a test database connection and runs migrations
func SetupTestDatabase(t *testing.T) *sql.DB {
	// Setup test configuration
	testConfig := GetTestConfig()
	
	// Create a test config struct and set it globally
	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode:    "test",
			Address: ":8080",
		},
		Database: config.DatabaseConfig{
			Host:     testConfig.DatabaseHost,
			Port:     testConfig.DatabasePort,
			User:     testConfig.DatabaseUser,
			Password: testConfig.DatabasePassword,
			Name:     testConfig.DatabaseName,
		},
		Paseto: config.PasetoConfig{
			SymmetricKey:  "JIusUnwiN236xiEVXMtRaJTNLyz6e0BD4U4pfX0CQNQ=",
			ExpireMinutes: 60,
		},
		Log: config.LogConfig{
			Level: "info",
		},
	}
	
	// Set the global config
	config.Cfg = cfg
	
	// First connect to MySQL server without specifying database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&loc=Local",
		testConfig.DatabaseUser,
		testConfig.DatabasePassword,
		testConfig.DatabaseHost,
		testConfig.DatabasePort,
	)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to MySQL server: %v", err)
	}

	// Create test database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", testConfig.DatabaseName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Close the connection and reconnect to the test database
	db.Close()

	// Connect to the test database
	testDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		testConfig.DatabaseUser,
		testConfig.DatabasePassword,
		testConfig.DatabaseHost,
		testConfig.DatabasePort,
		testConfig.DatabaseName,
	)

	db, err = sql.Open("mysql", testDSN)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := runTestMigrations(db); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	return db
}

// CleanupTestDatabase cleans up test data after tests
func CleanupTestDatabase(t *testing.T, db *sql.DB) {
	// Clean up test data
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		t.Logf("Warning: Failed to clean up transactions: %v", err)
	}
	
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Logf("Warning: Failed to clean up users: %v", err)
	}

	// Close database connection
	if err := db.Close(); err != nil {
		t.Logf("Warning: Failed to close test database: %v", err)
	}
}

// runTestMigrations runs the database schema for testing
func runTestMigrations(db *sql.DB) error {
	log.Println("[TEST-DB] Running test migrations...")
	
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create transactions table
	_, err = db.Exec(`
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
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	
	log.Println("[TEST-DB] Test migrations completed successfully")
	return nil
}

// TestHelper provides common test utilities
type TestHelper struct {
	DB *sql.DB
	T  *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	db := SetupTestDatabase(t)
	return &TestHelper{
		DB: db,
		T:  t,
	}
}

// Cleanup cleans up test resources
func (th *TestHelper) Cleanup() {
	CleanupTestDatabase(th.T, th.DB)
}

// CreateTestUser creates a test user in the database
func (th *TestHelper) CreateTestUser(name, email, password string) int64 {
	result, err := th.DB.Exec(
		"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		name, email, password,
	)
	if err != nil {
		th.T.Fatalf("Failed to create test user: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		th.T.Fatalf("Failed to get test user ID: %v", err)
	}
	
	return id
}

// CreateTestTransaction creates a test transaction in the database
func (th *TestHelper) CreateTestTransaction(userID int64, amount float64, description, txType string) int64 {
	result, err := th.DB.Exec(
		"INSERT INTO transactions (user_id, amount, description, type) VALUES (?, ?, ?, ?)",
		userID, amount, description, txType,
	)
	if err != nil {
		th.T.Fatalf("Failed to create test transaction: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		th.T.Fatalf("Failed to get test transaction ID: %v", err)
	}
	
	return id
}