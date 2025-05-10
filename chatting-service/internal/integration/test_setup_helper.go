package integration

import (
	"testing"

	"github.com/AmeerHeiba/chatting-service/internal/config"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test configuration
	cfg := config.DBConfig{
		Host:     "localhost",
		Port:     "5433", // Use a different port for tests
		User:     "postgres",
		Password: "postgres",
		DBName:   "chatting_service",
		SSLMode:  "disable",
	}

	// Connect to database
	db, err := config.NewDBConnection(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean database before each test
	err = db.Exec("DROP SCHEMA public CASCADE").Error
	if err != nil {
		t.Fatalf("Failed to drop schema: %v", err)
	}

	err = db.Exec("CREATE SCHEMA public").Error
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Run migrations
	err = runMigrations(db)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Return clean DB instance
	return db
}

func runMigrations(db *gorm.DB) error {
	// Same as your migrate/main.go but without logging
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	enums := []string{
		`CREATE TYPE message_status AS ENUM ('sent', 'delivered', 'read')`,
		`CREATE TYPE user_status AS ENUM ('online', 'offline', 'away')`,
		`CREATE TYPE message_type AS ENUM ('direct', 'broadcast')`,
	}

	for _, e := range enums {
		if err := db.Exec(e).Error; err != nil {
			return err
		}
	}

	models := []interface{}{
		&domain.User{},
		&domain.Message{},
		&domain.MessageRecipient{},
	}

	return db.AutoMigrate(models...)
}
