package application

import (
	"context"
	"strings"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/auth"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
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
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.tokenProvider.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
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
		return nil, shared.ErrUsernameTooShort
	}

	user, err := s.userService.VerifyCredentials(ctx, username, password)
	if err != nil {
		return nil, err
	}

	// Update last active
	if err := s.userService.UpdateUserLastActive(ctx, user.ID); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.tokenProvider.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
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
		return nil, domain.ErrInvalidToken
	}

	// Verify user still exists
	user, err := s.userService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Generate new tokens
	accessToken, err := s.tokenProvider.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenProvider.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
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
		return err
	}
	if !user.CheckPassword(currentPassword) {
		return domain.ErrInvalidCredentials
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, user.PasswordHash)
}
