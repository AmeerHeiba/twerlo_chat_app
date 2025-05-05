package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MediaHandler struct {
	mediaService *application.MediaService
}

func NewMediaHandler(mediaService *application.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

func (h *MediaHandler) Upload(c *fiber.Ctx) error {
	ctx := c.UserContext()
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		shared.Log.Error("failed to get file from form",
			zap.Error(err))
		return shared.ErrBadRequest.WithDetails("file is required")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		shared.Log.Error("failed to open uploaded file",
			zap.Error(err))
		return shared.ErrBadRequest.WithDetails("invalid file")
	}
	defer src.Close()

	// Upload
	response, err := h.mediaService.Upload(
		ctx,
		claims.UserID,
		src,
		file.Filename,
		file.Header.Get("Content-Type"),
		file.Size,
	)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
