package middleware

import (
	"strings"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(provider domain.TokenProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := provider.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store userID in context for WebSocket handler
		c.Locals("userID", claims.UserID)
		c.Locals("userClaims", claims)

		return c.Next()
	}
}
