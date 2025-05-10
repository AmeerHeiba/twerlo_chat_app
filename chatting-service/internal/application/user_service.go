package application

import (
	"context"
	"errors"
	"strings"

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
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		shared.Log.Error("username check failed", zap.Error(err), zap.String("username", username))
		return nil, err
	}
	if exists {
		shared.Log.Debug("username already exists", zap.String("username", username))
		return nil, shared.ErrUserExists.WithDetails("username already exists")
	}

	// Check if email exists
	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		shared.Log.Error("email check failed", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	if exists {
		shared.Log.Debug("email already exists", zap.String("email", email))
		return nil, shared.ErrUserExists.WithDetails("email already exists")
	}

	// Create domain user
	user := &domain.User{
		Username: username,
		Email:    email,
	}

	// Hash password
	if err := user.SetPassword(password); err != nil {
		shared.Log.Error("password must be at least 8 characters", zap.Error(err), zap.String("hashed password", password))
		return nil, shared.ErrValidation.WithDetails("Set Password failed").WithDetails(err.Error())
	}

	// Persist user
	createdUser, err := s.userRepo.Create(ctx, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		shared.Log.Error("create user failed", zap.Error(err), zap.String("username", username), zap.String("email", email), zap.String("password_hash", user.PasswordHash))
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		shared.Log.Error("find user by ID failed", zap.Error(err), zap.Uint("userID", id))
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, id uint) (*domain.User, error) {
	profile, err := s.userRepo.FindProfileByID(ctx, id)
	if err != nil {
		shared.Log.Error("find user profile by ID failed", zap.Error(err), zap.Uint("userID", id))
		return nil, err
	}
	return profile, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		shared.Log.Error("find user by username failed", zap.Error(err), zap.String("username", username))
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUserLastActive(ctx context.Context, userID uint) error {
	err := s.userRepo.UpdateLastActiveAt(ctx, userID)
	if err != nil {
		shared.Log.Error("update user last active failed", zap.Error(err), zap.Uint("userID", userID))
		return err
	}
	return nil
}

func (s *UserService) VerifyCredentials(ctx context.Context, username, password string) (*domain.User, error) {

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		shared.Log.Error("Failed to find user", zap.Error(err), zap.String("username", username))
		return nil, err
	}
	// Direct password check without validation
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			shared.Log.Error("Password comparison failed", zap.Error(err), zap.String("username", username))
			return nil, shared.ErrInvalidCredentials.WithDetails("password comparison failed")
		}
		shared.Log.Error("Password comparison failed", zap.Error(err), zap.String("username", username))
		return nil, shared.ErrInvalidCredentials.WithDetails("password comparison failed")
	}

	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, username, email string) (*domain.User, error) {
	// Input validation (e.g., email format, username length)
	if username != "" && len(username) < 3 {
		shared.Log.Debug("username must be at least 3 characters long", zap.String("username", username))
		return nil, shared.ErrValidation.WithDetails("username must be at least 3 characters long")
	}
	if email != "" && !strings.Contains(email, "@") {
		shared.Log.Debug("invalid email format", zap.String("email", email))
		return nil, shared.ErrValidation.WithDetails("invalid email format")
	}

	// Delegate to repository
	if err := s.userRepo.Update(ctx, userID, username, email); err != nil {
		shared.Log.Error("update user failed", zap.Error(err), zap.Uint("userID", userID), zap.String("username", username), zap.String("email", email))
		return nil, err
	}

	updatedProfile, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		shared.Log.Error("find user by ID failed", zap.Error(err), zap.Uint("userID", userID))
		return nil, err
	}

	return updatedProfile, nil
}

func (s *UserService) GettAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetMessageHistory(ctx context.Context, userID uint, limit, offset int) ([]domain.Message, int64, error) {
	// I'll implement this fully after I create the message logic
	// For now just the signature
	return nil, 0, nil
}
