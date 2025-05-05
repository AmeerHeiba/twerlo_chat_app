package middleware

import (
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger, ok := c.Locals("logger").(*zap.Logger)
				if !ok {
					logger = shared.Log
				}

				logger.Error("Recovered from panic",
					zap.Any("panic", r),
					zap.Stack("stack"))

				err = shared.ErrInternalServer.WithDetails("internal server error")
			}
		}()

		return c.Next()
	}
}
