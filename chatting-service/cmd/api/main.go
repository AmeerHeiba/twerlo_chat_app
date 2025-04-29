package main

import (
	"log"
	"os"

	"github.com/AmeerHeiba/chatting-service/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := config.NewDBConnection(config.LoadDBConfig())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	app := fiber.New()

	// Health check with DB verification
	app.Get("/api/health", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "down",
				"message": "Database connection failed",
			})
		}

		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "down",
				"message": "Database ping failed",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Chatting Service is running 🚀",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
