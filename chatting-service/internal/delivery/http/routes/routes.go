package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/middleware"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB             *gorm.DB
	UserHandler    *handlers.UserHandler
	AuthHandler    *handlers.AuthHandler
	MessageHandler *handlers.MessageHandler
	MediaHandler   *handlers.MediaHandler
	JWTProvider    domain.TokenProvider
}

func SetupRoutes(app *fiber.App, deps Dependencies) {
	// Health route (no auth)
	SetupHealthRoutes(app, deps.DB)

	// Auth routes (no auth)
	SetupAuthRoutes(app, deps.AuthHandler, middleware.NewAuthMiddleware(deps.JWTProvider))

	// Profile routes (protected)
	SetupUserRoutes(app, deps.UserHandler, middleware.NewAuthMiddleware(deps.JWTProvider))

	// Message routes (protected)
	SetupMessageRoutes(app, deps.MessageHandler, middleware.NewAuthMiddleware(deps.JWTProvider))

	// Media routes (protected)
	SetupMediaRoutes(app, deps.MediaHandler, middleware.NewAuthMiddleware(deps.JWTProvider))

}
