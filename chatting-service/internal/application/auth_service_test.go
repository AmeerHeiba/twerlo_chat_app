package application

import (
	"context"
	"testing"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/auth"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, userName, email, passwordHash string) (*domain.User, error) {
	args := m.Called(ctx, userName, email, passwordHash)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID uint) (*domain.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindProfileByID(ctx context.Context, userID uint) (*domain.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, userID uint, username, email string) error {
	args := m.Called(ctx, userID, username, email)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID uint, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastActiveAt(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, userID uint) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

type MockTokenProvider struct {
	mock.Mock
}

func (m *MockTokenProvider) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenProvider) GenerateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenProvider) ValidateToken(ctx context.Context, tokenString string) (*domain.TokenClaims, error) {
	args := m.Called(ctx, tokenString)
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockTokenProvider) ValidateRefreshToken(ctx context.Context, tokenString string) (*domain.TokenClaims, error) {
	args := m.Called(ctx, tokenString)
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockTokenProvider) GetAccessExpiry() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockTokenProvider) GetRefreshExpiry() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}
func TestAuthService_Register(t *testing.T) {
	shared.InitLogger("test")

	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		mockSetup   func(*MockUserRepository, *MockTokenProvider)
		expected    *auth.AuthResponse
		expectedErr error
	}{
		{
			name:     "Success",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(userRepo *MockUserRepository, tokenProvider *MockTokenProvider) {
				// Setup all required mock expectations
				userRepo.On("ExistsByUsername", mock.Anything, "testuser").Return(false, nil)
				userRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
				userRepo.On("Create", mock.Anything, "testuser", "test@example.com", mock.Anything).
					Return(&domain.User{

						Username: "testuser",
						Email:    "test@example.com",
					}, nil)

				tokenProvider.On("GetAccessExpiry").Return(time.Hour)
				tokenProvider.On("GenerateToken", mock.Anything, mock.Anything).Return("access_token", nil)
				tokenProvider.On("GenerateRefreshToken", mock.Anything, mock.Anything).Return("refresh_token", nil)
			},
			expected: &auth.AuthResponse{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				ExpiresIn:    3600, // 1 hour in seconds
				TokenType:    "Bearer",
				Username:     "testuser",
				Email:        "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &MockUserRepository{}
			tokenProvider := &MockTokenProvider{}
			userService := NewUserService(userRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(userRepo, tokenProvider)
			}

			authService := NewAuthService(userRepo, userService, tokenProvider)
			res, err := authService.Register(context.Background(), tt.username, tt.email, tt.password)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, res)
			}

			userRepo.AssertExpectations(t)
			tokenProvider.AssertExpectations(t)
		})
	}
}
