package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/application"
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
