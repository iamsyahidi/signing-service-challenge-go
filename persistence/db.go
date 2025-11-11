package persistence

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// DBConfig holds the database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection and runs migrations
func NewPostgresDB(cfg DBConfig) (*sql.DB, error) {
	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Get the absolute path to the migrations directory
	// This uses the module path to find the migrations directory
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "migrations")

	// Run migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}
