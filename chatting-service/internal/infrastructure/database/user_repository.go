package database

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// COMMAND OPERATIONS (Write)

func (r *userRepository) Create(ctx context.Context, username, email, passwordHash string) (*domain.User, error) {
	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if validationErr := user.ValidateRegistration(); validationErr != nil {
		return nil, validationErr
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				if strings.Contains(err.Error(), "users_username_key") {
					shared.Log.Debug("username already exists", zap.String("username", username))
					return shared.ErrValidation.WithDetails("username already exists")
				}
				if strings.Contains(err.Error(), "users_email_key") {
					shared.Log.Debug("email already exists", zap.String("email", email))
					return shared.ErrValidation.WithDetails("email already exists")
				}
			}
			shared.Log.Error("create user failed",
				zap.String("operation", "Create"),
				zap.String("username", username),
				zap.String("email", email),
				zap.String("password_hash", passwordHash),
				zap.Error(err))
			return shared.ErrDatabaseOperation.WithDetails("create user failed").WithDetails(err.Error())
		}
		return nil
	})

	return user, err
}

func (r *userRepository) Update(ctx context.Context, userID uint, username, email string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if username != "" {
			exists, err := r.usernameExists(tx, username, userID)
			if err != nil {
				shared.Log.Error("username check failed",
					zap.String("operation", "Update"),
					zap.Uint("userID", userID),
					zap.String("username", username),
					zap.Error(err))
				return shared.ErrDatabaseOperation.WithDetails("username check failed").WithDetails(err.Error())
			}
			if exists {
				return shared.ErrDuplicateEntry.WithDetails("username already exists")
			}
		}

		if email != "" {
			exists, err := r.emailExists(tx, email, userID)
			if err != nil {
				shared.Log.Error("email check failed",
					zap.String("operation", "Update"),
					zap.Uint("userID", userID),
					zap.String("email", email),
					zap.Error(err))
				return shared.ErrDatabaseOperation.WithDetails("email check failed").WithDetails(err.Error())
			}
			if exists {
				return shared.ErrDuplicateEntry.WithDetails("email already exists")
			}
		}

		updates := map[string]interface{}{}
		if username != "" {
			updates["username"] = username
		}
		if email != "" {
			updates["email"] = email
		}

		if len(updates) > 0 {
			if err := tx.Model(&domain.User{}).
				Where("id = ?", userID).
				Updates(updates).Error; err != nil {
				shared.Log.Error("update user failed",
					zap.String("operation", "Update"),
					zap.Uint("userID", userID),
					zap.Any("updates", updates),
					zap.Error(err))
				return shared.ErrDatabaseOperation.WithDetails("update user failed").WithDetails(err.Error())
			}
		}

		return nil
	})
}

func (r *userRepository) UpdateLastActiveAt(ctx context.Context, userID uint) error {
	err := r.db.WithContext(ctx).Exec(
		"UPDATE users SET last_active_at = ? WHERE id = ?",
		time.Now().UTC(),
		userID,
	).Error
	if err != nil {
		shared.Log.Error("update last active at failed",
			zap.String("operation", "UpdateLastActiveAt"),
			zap.Uint("userID", userID),
			zap.Error(err))
		return shared.ErrDatabaseOperation.WithDetails("update last active at failed").WithDetails(err.Error())
	}
	return nil
}

// QUERY OPERATIONS (Read)

func (r *userRepository) FindByID(ctx context.Context, userID uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		shared.Log.Debug("user not found",
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, shared.ErrRecordNotFound.WithDetails("user not found").WithDetails(err.Error())
	}
	if err != nil {
		shared.Log.Error("find user by ID failed",
			zap.String("operation", "FindByID"),
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, shared.ErrDatabaseOperation.WithDetails("find user by ID failed").WithDetails(err.Error())
	}
	return &user, nil
}

func (r *userRepository) FindProfileByID(ctx context.Context, userID uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("id", "username", "email", "last_active_at", "status").
		First(&user, userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		shared.Log.Debug("user profile not found",
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, shared.ErrRecordNotFound.WithDetails("user profile not found").WithDetails(err.Error())
	}
	if err != nil {
		shared.Log.Error("find user profile by ID failed",
			zap.String("operation", "FindProfileByID"),
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, shared.ErrDatabaseOperation.WithDetails("find user profile by ID failed").WithDetails(err.Error())
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("id", "username", "email", "password_hash", "last_active_at").
		Where("username = ?", username).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		shared.Log.Debug("user not found", zap.String("username", username))
		return nil, shared.ErrRecordNotFound.WithDetails("user not found").WithDetails(err.Error())
	}
	if err != nil {
		shared.Log.Error("find user by username failed",
			zap.String("operation", "FindByUsername"),
			zap.String("username", username),
			zap.Error(err))
		return nil, shared.ErrRecordNotFound.WithDetails("find user by username failed").WithDetails(err.Error())
	}
	return &user, nil
}

func (r *userRepository) Exists(ctx context.Context, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Count(&count).Error

	if err != nil {
		shared.Log.Error("user exists check failed",
			zap.String("operation", "Exists"),
			zap.Uint("userID", userID),
			zap.Error(err))
		return false, shared.ErrDatabaseOperation.WithDetails("user exists check failed").WithDetails(err.Error())
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		shared.Log.Error("user exists check failed",
			zap.String("operation", "ExistsByUsername"),
			zap.String("username", username),
			zap.Error(err))
		return false, shared.ErrDatabaseOperation.WithDetails("user exists check failed").WithDetails(err.Error())
	}
	return count > 0, nil
}

// Helper methods (private to repository)
func (r *userRepository) usernameExists(tx *gorm.DB, username string, excludeID uint) (bool, error) {
	var count int64
	if err := tx.Model(&domain.User{}).
		Where("username = ? AND id <> ?", username, excludeID).
		Count(&count).Error; err != nil {
		shared.Log.Error("username check failed",
			zap.String("operation", "usernameExists"),
			zap.String("username", username),
			zap.Error(err))
		return false, shared.ErrDatabaseOperation.WithDetails("username check failed").WithDetails(err.Error())
	}
	return count > 0, nil
}

func (r *userRepository) emailExists(tx *gorm.DB, email string, excludeID uint) (bool, error) {
	var count int64
	if err := tx.Model(&domain.User{}).
		Where("email = ? AND id <> ?", email, excludeID).
		Count(&count).Error; err != nil {
		shared.Log.Error("email check failed",
			zap.String("operation", "emailExists"),
			zap.String("email", email),
			zap.Error(err))
		return false, shared.ErrDatabaseOperation.WithDetails("email check failed").WithDetails(err.Error())
	}
	return count > 0, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uint, passwordHash string) error {
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
	if err != nil {
		shared.Log.Error("update password failed",
			zap.String("operation", "UpdatePassword"),
			zap.Uint("userID", userID),
			zap.Error(err))
		return shared.ErrDatabaseOperation.WithDetails("update password failed").WithDetails(err.Error())
	}
	return nil
}
