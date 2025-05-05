package application

import (
	"context"
	"strings"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/auth"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepo      domain.UserRepository
	userService   *UserService
	tokenProvider domain.TokenProvider
}

func NewAuthService(
	repo domain.UserRepository,
	userService *UserService,
	provider domain.TokenProvider,
) *AuthService {
	return &AuthService{
		userRepo:      repo,
		userService:   userService,
		tokenProvider: provider,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*auth.AuthResponse, error) {
	// Create user through service (includes validation)
	user, err := s.userService.CreateUser(ctx, username, email, password)
	if err != nil {
		shared.Log.Error("create user failed",
			zap.String("operation", "Create"),
			zap.Error(err),
			zap.String("username", username),
			zap.String("email", email))
		return nil, shared.ErrDatabaseOperation.WithDetails("create user failed").WithDetails(err.Error())
	}

	// Generate tokens
	accessToken, err := s.tokenProvider.GenerateToken(ctx, user)
	if err != nil {
		shared.Log.Error("generate token failed", zap.Error(err))
		return nil, shared.ErrInternalServer.WithDetails("generate token failed").WithDetails(err.Error())
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(ctx, user)
	if err != nil {
		shared.Log.Error("generate refresh token failed", zap.Error(err))
		return nil, shared.ErrInternalServer.WithDetails("generate refresh token failed").WithDetails(err.Error())
	}

	return &auth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.tokenProvider.GetAccessExpiry().Seconds()),
		TokenType:    "Bearer",
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*auth.AuthResponse, error) {

	// Trim input first
	username = strings.TrimSpace(username)

	if len(username) < 3 {
		shared.Log.Debug("username must be at least 3 characters long", zap.String("username", username))
		return nil, shared.ErrBadRequest.WithDetails("Invalid username")
	}

	user, err := s.userService.VerifyCredentials(ctx, username, password)
	if err != nil {
		shared.Log.Error("verify credentials failed", zap.Error(err), zap.String("username", username))
		return nil, shared.ErrValidation.WithDetails("invalid credentials").WithDetails(err.Error())
	}

	// Update last active
	if err := s.userService.UpdateUserLastActive(ctx, user.ID); err != nil {
		shared.Log.Error("update user last active failed", zap.Error(err), zap.Uint("userID", user.ID))
		return nil, shared.ErrDatabaseOperation.WithDetails("update user last active failed").WithDetails(err.Error())
	}

	// Generate tokens
	accessToken, err := s.tokenProvider.GenerateToken(ctx, user)
	if err != nil {
		shared.Log.Error("generate token failed", zap.Error(err))
		return nil, shared.ErrInternalServer.WithDetails("generate token failed").WithDetails(err.Error())
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(ctx, user)
	if err != nil {
		shared.Log.Error("generate refresh token failed", zap.Error(err))
		return nil, shared.ErrInternalServer.WithDetails("generate refresh token failed").WithDetails(err.Error())
	}

	return &auth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.tokenProvider.GetAccessExpiry().Seconds()),
		TokenType:    "Bearer",
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*auth.AuthResponse, error) {
	claims, err := s.tokenProvider.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		shared.Log.Error("validate refresh token failed", zap.Error(err))
		return nil, shared.ErrUnauthorized.WithDetails("invalid refresh token").WithDetails(err.Error())
	}

	// Verify user still exists
	user, err := s.userService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		shared.Log.Error("find user by ID failed", zap.Error(err), zap.Uint("userID", claims.UserID))
		return nil, shared.ErrDatabaseOperation.WithDetails("find user by ID failed").WithDetails(err.Error())
	}

	// Generate new tokens
	accessToken, err := s.tokenProvider.GenerateToken(ctx, user)
	if err != nil {
		shared.Log.Error("generate token failed", zap.Error(err))
		return nil, shared.ErrInternalServer.WithDetails("generate token failed").WithDetails(err.Error())
	}

	newRefreshToken, err := s.tokenProvider.GenerateRefreshToken(ctx, user)
	if err != nil {
		shared.Log.Error("generate refresh token failed", zap.Error(err))
		return nil, shared.ErrInternalServer.WithDetails("generate refresh token failed").WithDetails(err.Error())
	}

	return &auth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.tokenProvider.GetAccessExpiry().Seconds()),
		TokenType:    "Bearer",
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// In a JWT system, you would typically:
	// 1. Add token to blacklist
	// 2. Clear client-side tokens
	// 3. For immediate invalidation
	//    - Shorten token expiry
	//    - Implement token blacklisting
	//    - Use refresh token rotation

	// For now, I'll just validate the token to ensure it was valid
	//TILL I DECIDE LATER BASED ON TIM LINE
	_, err := s.tokenProvider.ValidateToken(ctx, token)
	return err
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	return s.tokenProvider.ValidateToken(ctx, token)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		shared.Log.Error("find user by ID failed", zap.Error(err), zap.Uint("userID", userID))
		return shared.ErrDatabaseOperation.WithDetails("find user by ID failed").WithDetails(err.Error())
	}
	if !user.CheckPassword(currentPassword) {
		shared.Log.Debug("invalid credentials", zap.String("username", user.Username))
		return shared.ErrInvalidCredentials
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		shared.Log.Error("password must be at least 8 characters", zap.Error(err), zap.String("hashed password", newPassword))
		return shared.ErrValidation.WithDetails("password must be at least 8 characters").WithDetails(err.Error())
	}

	return s.userRepo.UpdatePassword(ctx, userID, user.PasswordHash)
}
