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

// Upload handles media file upload for authenticated users.
//
// @Summary Upload a media file
// @Description Upload a media file (JPEG, PNG, PDF) with a maximum size of 10MB.
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "Media File (JPEG, PNG, PDF)"
// @Success 200 {object} domain.MediaResponse "Successfully uploaded"
// @Failure 400 {object} shared.Error "Bad request (missing or invalid file)"
// @Failure 401 {object} shared.Error "Unauthorized"
// @Failure 500 {object} shared.Error "Internal server error"
// @Router /api/media/upload [post]
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
		claims.UserID,
	)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
