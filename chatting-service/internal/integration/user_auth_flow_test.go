package integration

import (
	"context"
	"testing"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/config"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/auth"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/database"
	"github.com/stretchr/testify/assert"
)

func TestCompleteAuthFlow(t *testing.T) {
	db := setupTestDB(t)
	userRepo := database.NewUserRepository(db)
	userService := application.NewUserService(userRepo)

	authCfg := config.LoadAuthConfig()
	jwtProvider := auth.NewJWTProvider(authCfg)
	authService := application.NewAuthService(userRepo, userService, jwtProvider)

	// Test registration
	registerResp, err := authService.Register(
		context.Background(),
		"testuser",
		"test@example.com",
		"password123",
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, registerResp.AccessToken)
	assert.NotEmpty(t, registerResp.RefreshToken)

	// Test login
	loginResp, err := authService.Login(
		context.Background(),
		"testuser",
		"password123",
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResp.AccessToken)

	// Test token validation
	claims, err := jwtProvider.ValidateToken(
		context.Background(),
		loginResp.AccessToken,
	)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)

	// Test refresh token
	refreshResp, err := authService.Refresh(
		context.Background(),
		loginResp.RefreshToken,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEqual(t, loginResp.AccessToken, refreshResp.AccessToken)
}
