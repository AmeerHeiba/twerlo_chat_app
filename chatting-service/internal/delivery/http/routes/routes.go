package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/middleware"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB          *gorm.DB
	UserHandler *handlers.UserHandler
	AuthHandler *handlers.AuthHandler
	JWTProvider domain.TokenProvider
}

func SetupRoutes(app *fiber.App, deps Dependencies) {
	// Health route (no auth)
	SetupHealthRoutes(app, deps.DB)

	// Auth routes (no auth)
	SetupAuthRoutes(app, deps.AuthHandler)

	// Profile routes (protected)
	SetupUserRoutes(app, deps.UserHandler, middleware.NewAuthMiddleware(deps.JWTProvider))

}
