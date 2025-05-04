package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupHealthRoutes(app *fiber.App, db *gorm.DB) {
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
}
