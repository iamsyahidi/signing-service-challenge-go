package main

import (
	"log"
	"os"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

const (
	ListenAddress = ":8080"
)

func main() {
	// Initialize dependencies
	// Database configuration
	cfg := persistence.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "signing_service"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database connection and run migrations
	db, err := persistence.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create repository
	repository := persistence.NewPostgresDeviceRepositoryFromDB(db)
	deviceService := domain.NewDeviceService(repository)
	server := api.NewServer(ListenAddress, deviceService)

	log.Printf("Starting signature service on %s", ListenAddress)
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
