package application

import (
	"context"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	// Check if user exists
	exists, err := s.userRepo.Exists(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, shared.ErrUserExists
	}

	// Create domain user
	user := &domain.User{
		Username: username,
		Email:    email,
	}

	// Hash password
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// Persist user
	createdUser, err := s.userRepo.Create(ctx, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepo.FindByUsername(ctx, username)
}

func (s *UserService) UpdateUserLastActive(ctx context.Context, userID uint) error {
	return s.userRepo.UpdateLastActiveAt(ctx, userID)
}

func (s *UserService) VerifyCredentials(ctx context.Context, username, password string) (*domain.User, error) {
	shared.Log.Debug("VerifyCredentials input",
		zap.String("username", username),
		zap.Int("length", len(username)))

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	shared.Log.Debug("User retrieved",
		zap.String("db_username", user.Username),
		zap.Int("db_length", len(user.Username)))

	// Direct password check without validation
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}
