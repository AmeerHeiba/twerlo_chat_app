package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, handler *handlers.AuthHandler) {
	auth := app.Group("/api/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
}
