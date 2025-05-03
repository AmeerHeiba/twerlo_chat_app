package middleware

import (
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err != nil {
		// Convert to standardized HTTP error
		httpErr := shared.ToHTTPError(err)

		// Log structured error information
		logFields := []zap.Field{
			zap.Error(err),
			zap.String("path", ctx.Path()),
			zap.String("method", ctx.Method()),
			zap.Int("status", httpErr.Status),
			zap.String("error_code", httpErr.Code),
		}

		// Add request ID if available
		if requestID := ctx.GetRespHeader("X-Request-ID"); requestID != "" {
			logFields = append(logFields, zap.String("request_id", requestID))
		}

		// Log differently based on error type
		switch {
		case httpErr.Status >= 500:
			shared.Log.Error("Server error", logFields...)
		case httpErr.Status >= 400:
			shared.Log.Warn("Client error", logFields...)
		default:
			shared.Log.Info("Application error", logFields...)
		}

		// Return error response
		return ctx.Status(httpErr.Status).JSON(httpErr)
	}

	return nil
}
