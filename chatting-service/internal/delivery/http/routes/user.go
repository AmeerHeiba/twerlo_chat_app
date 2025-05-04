package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, handler *handlers.UserHandler, authMiddleware fiber.Handler) {
	user := app.Group("/api/users", authMiddleware)
	user.Get("/profile", handler.GetUserProfile)
	user.Put("/profile", handler.UpdateProfile)
	user.Get("/messages", handler.GetMessageHistory)

}
