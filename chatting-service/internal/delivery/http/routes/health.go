package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupHealthRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/api/health", func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*zap.Logger)
		start := time.Now()

		sqlDB, err := db.DB()
		if err != nil {
			logger.Error("Database connection failed",
				zap.Error(err),
				zap.Duration("latency", time.Since(start)))
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "error"})
		}

		if err := sqlDB.Ping(); err != nil {
			logger.Error("Database ping failed",
				zap.Error(err),
				zap.Duration("latency", time.Since(start)))
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "error"})
		}

		logger.Info("Health check succeeded",
			zap.Duration("latency", time.Since(start)))
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
