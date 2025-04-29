package main

import (
	"fmt"
	"log"

	"github.com/AmeerHeiba/chatting-service/internal/config"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"gorm.io/gorm"
)

func main() {
	dbConfig := config.LoadDBConfig()
	db, err := config.NewDBConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := runMigrations(db); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func runMigrations(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create UUID extension: %w", err)
	}

	// Create enum types first
	enums := []string{
		`CREATE TYPE message_status AS ENUM ('sent', 'delivered', 'read')`,
		`CREATE TYPE user_status AS ENUM ('online', 'offline', 'away')`,
		`CREATE TYPE message_type AS ENUM ('direct', 'broadcast')`,
	}

	for _, e := range enums {
		if err := db.Exec(e).Error; err != nil {
			log.Printf("Note: Enum creation error (might already exist): %v", err)
		}
	}

	// Migrate models
	models := []interface{}{
		&domain.User{},
		&domain.Message{},
		&domain.MessageRecipient{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}
