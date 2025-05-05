package middleware

import (
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)

		// Create contextual logger
		logger := shared.Log.With(
			zap.String("request_id", requestID),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("ip", c.IP()),
		)
		c.Locals("logger", logger)

		return c.Next()
	}
}
