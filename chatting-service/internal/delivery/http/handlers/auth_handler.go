package handlers

import (
	"errors"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/auth"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := h.authService.Login(c.Context(), body.Username, body.Password)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := h.authService.Register(c.Context(), body.Username, body.Email, body.Password)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	claims, ok := c.Locals("userClaims").(*domain.TokenClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var body auth.ChangePasswordRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.authService.ChangePassword(c.Context(), claims.UserID, body.CurrentPassword, body.NewPassword); err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Current password is incorrect",
			})
		}
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}
