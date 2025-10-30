package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luthfiarsyad/mms/config"
)

var DB *sql.DB

// Connect membuka koneksi ke database MySQL menggunakan konfigurasi dari config package.
func Connect() (*sql.DB, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("config is not loaded")
	}

	var dsn string
	if cfg.Database.DSN != "" {
		dsn = cfg.Database.DSN
	} else {
		// format DSN: user:password@tcp(host:port)/dbname?parseTime=true&loc=Local
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Optional: set connection pool config
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	log.Println("[DB] Connected to MySQL successfully")

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	DB = db
	return db, nil
}

// Get returns the global *sql.DB instance (after Connect() is called)
func Get() *sql.DB {
	return DB
}

// Close menutup koneksi database
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// runMigrations executes the database schema migration
func runMigrations(db *sql.DB) error {
	log.Println("[DB] Running migrations...")

	// Read schema file
	schemaPath := filepath.Join("internal", "infrastructure", "persistence", "migration", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema
	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("[DB] Migrations completed successfully")
	return nil
}
