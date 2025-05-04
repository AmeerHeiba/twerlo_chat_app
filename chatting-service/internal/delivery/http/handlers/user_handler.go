package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/user"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updatedUser, err := h.userService.UpdateProfile(c.Context(), claims.UserID, body.Username, body.Email)
	if err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Default values
	if query.Limit == 0 {
		query.Limit = 20
	}

	messages, total, err := h.userService.GetMessageHistory(c.Context(), claims.UserID, query.Limit, query.Offset)
	if err != nil {
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
