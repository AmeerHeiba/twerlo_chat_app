package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"

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

	profile, err := h.userService.GetUserByID(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.JSON(profile)
}
