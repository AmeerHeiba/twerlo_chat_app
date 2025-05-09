package handlers

import (
	"regexp"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/auth"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary      Login to the system
// @Description  Authenticate a user and get a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      auth.LoginRequest  true  "Login credentials"
// @Success      200   {object}  auth.AuthResponse
// @Failure      400   {object}  shared.Error
// @Failure      401   {object}  shared.Error
// @Failure      500   {object}  shared.Error
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		shared.Log.Error("Invalid request body", zap.Error(err), zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body")
	}

	if body.Username == "" || body.Password == "" {
		shared.Log.Debug("Invalid request body must have all fields", zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body must have all fields")
	}

	res, err := h.authService.Login(c.Context(), body.Username, body.Password)
	if err != nil {
		shared.Log.Error("Login failed", zap.Error(err))
		return err
	}

	return c.JSON(res)
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with username, email, and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      auth.RegisterRequest  true  "Regesteration data"
// @Success      200   {object}  auth.AuthResponse
// @Failure      400   {object}  shared.Error
// @Failure      422   {object}  shared.Error
// @Failure      500   {object}  shared.Error
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		shared.Log.Error("Invalid request body", zap.Error(err), zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body").WithDetails(err.Error())
	}

	if body.Username == "" || body.Email == "" || body.Password == "" {
		shared.Log.Debug("Invalid request body must have all fields", zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body must have all fields")
	}

	if len(body.Password) < 8 {
		shared.Log.Debug("Password must be at least 8 characters long", zap.ByteString("body", c.Body()))
		return shared.ErrValidation.WithDetails("Password must be at least 8 characters long")
	}

	emailRegex := regexp.MustCompile(shared.EmailRegexPattern)
	if !emailRegex.MatchString(body.Email) {
		shared.Log.Debug("Please provide a valid email address", zap.ByteString("body", c.Body()))
		return shared.ErrInvalidEmailFormat.WithDetails("Please provide a valid email address")
	}

	res, err := h.authService.Register(c.Context(), body.Username, body.Email, body.Password)
	if err != nil {
		shared.Log.Error("Register failed", zap.Error(err))
		return err
	}

	return c.JSON(res)
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Authenticated user changes their password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body  auth.ChangePasswordRequest  true  "Change Password Request"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  shared.Error
// @Failure      401   {object}  shared.Error
// @Failure      500   {object}  shared.Error
// @Router       /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	claims, ok := c.Locals("userClaims").(*domain.TokenClaims)
	if !ok || claims == nil {
		shared.Log.Debug("Invalid user claims", zap.ByteString("body", c.Body()))
		return shared.ErrUnauthorized.WithDetails("Invalid user claims")
	}

	var body auth.ChangePasswordRequest
	if err := c.BodyParser(&body); err != nil {
		shared.Log.Debug("Invalid request body", zap.Error(err), zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body").WithDetails(err.Error())
	}

	if err := h.authService.ChangePassword(c.Context(), claims.UserID, body.CurrentPassword, body.NewPassword); err != nil {
		shared.Log.Error("Change password failed", zap.Error(err))
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}
