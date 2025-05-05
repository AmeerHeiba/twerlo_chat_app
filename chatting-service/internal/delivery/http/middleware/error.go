package middleware

import (
	"errors"

	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ErrorHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err == nil {
		return nil
	}

	logger := ctx.Locals("logger").(*zap.Logger)
	var appErr shared.Error

	if !errors.As(err, &appErr) {
		appErr = shared.ErrInternalServer.WithDetails(err.Error())
		logger.Error("Unhandled error",
			zap.Error(err),
			zap.String("path", ctx.Path()),
			zap.ByteString("body", ctx.Body()))
	} else {
		logLevel := zapcore.ErrorLevel
		if appErr.Status < 500 {
			logLevel = zapcore.DebugLevel
		}

		logger.Log(logLevel, "Request error",
			zap.Error(err),
			zap.String("path", ctx.Path()),
			zap.Int("status", appErr.Status),
			zap.Any("details", appErr.Details))
	}

	return ctx.Status(appErr.Status).JSON(appErr)
}
