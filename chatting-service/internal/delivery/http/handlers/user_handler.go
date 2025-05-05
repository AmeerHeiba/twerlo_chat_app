package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/user"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	// Get user ID from JWT claims
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	profile, err := h.userService.GetUserProfile(c.Context(), claims.UserID)
	if err != nil {
		shared.Log.Error("Failed to get user profile", zap.Error(err))
		return err
	}

	return c.JSON(user.ProfileResponse{
		ID:         profile.ID,
		Username:   profile.Username,
		Email:      profile.Email,
		LastActive: profile.LastActiveAt,
		Status:     string(profile.Status),
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	var body user.UpdateProfileRequest
	if err := c.BodyParser(&body); err != nil {
		shared.Log.Error("Invalid update profile request body",
			zap.Error(err),
			zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body failed to parse request body").WithDetails(err.Error())
	}

	updatedUser, err := h.userService.UpdateProfile(c.Context(), claims.UserID, body.Username, body.Email)
	if err != nil {
		shared.Log.Error("Failed to update user profile", zap.Error(err))
		return err
	}

	return c.JSON(user.ProfileResponse{
		ID:         updatedUser.ID,
		Username:   updatedUser.Username,
		Email:      updatedUser.Email,
		LastActive: updatedUser.LastActiveAt,
		Status:     string(updatedUser.Status),
	})
}

func (h *UserHandler) GetMessageHistory(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	var query user.MessageHistoryRequest
	if err := c.QueryParser(&query); err != nil {
		shared.Log.Error("Invalid message history request query", zap.Error(err), zap.Any("query", query))
		return shared.ErrBadRequest.WithDetails("Invalid or missing query params").WithDetails(err.Error())
	}

	// Default values
	if query.Limit == 0 {
		query.Limit = 20
		shared.Log.Debug("Using default limit", zap.Int("limit", query.Limit))
	}

	messages, total, err := h.userService.GetMessageHistory(c.Context(), claims.UserID, query.Limit, query.Offset)
	if err != nil {
		shared.Log.Error("Failed to get message history", zap.Error(err))
		return err
	}

	response := user.MessageHistoryResponse{
		Total:    total,
		Messages: make([]user.MessageResponse, 0, len(messages)),
	}

	for _, msg := range messages {
		response.Messages = append(response.Messages, user.MessageResponse{
			ID:        msg.ID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			Status:    string(msg.Status),
			// Type:      string(msg.Type),
		})
	}

	return c.JSON(response)
}
