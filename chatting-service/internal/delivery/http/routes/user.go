package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, handler *handlers.UserHandler, authMiddleware fiber.Handler) {
	user := app.Group("/api/users", authMiddleware)
	user.Get("/", handler.GetUserProfile)
	// user.Get("/", handler.ListUsers)
	// user.Get("/:id", handler.GetUser)
}
