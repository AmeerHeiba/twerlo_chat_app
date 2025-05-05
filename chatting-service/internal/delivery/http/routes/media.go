package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupMediaRoutes(app *fiber.App, handler *handlers.MediaHandler, authMiddleware fiber.Handler) {
	media := app.Group("/api/media", authMiddleware)
	media.Post("/upload", handler.Upload)
}
