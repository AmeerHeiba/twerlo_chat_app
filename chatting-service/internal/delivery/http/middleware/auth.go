package middleware

import (
	"strings"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(provider domain.TokenProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check Authorization header first
		authHeader := c.Get("Authorization")
		var token string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Fallback to query parameter s
			token = c.Query("token")
		}

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization token missing",
			})
		}

		claims, err := provider.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userClaims", claims)

		return c.Next()
	}
}
