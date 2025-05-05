package middleware

import (
	"context"
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

		ctx := context.WithValue(c.Context(), "userID", claims.UserID)
		c.SetUserContext(ctx)

		// Store claims in context for downstream handlers
		c.Locals("userClaims", claims)
		return c.Next()
	}
}
