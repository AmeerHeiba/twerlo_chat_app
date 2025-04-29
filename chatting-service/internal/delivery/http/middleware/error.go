package middleware

import (
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err != nil {
		httpErr := shared.ToHTTPError(err)

		// Log unexpected errors (500s)
		if httpErr.Code >= 500 {
			shared.Log.Error("Server error",
				zap.Error(err),
				zap.String("path", ctx.Path()),
				zap.String("method", ctx.Method()),
			)
		}

		return ctx.Status(httpErr.Code).JSON(httpErr)
	}

	return nil
}
