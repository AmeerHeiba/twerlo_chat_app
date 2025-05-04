package database

import (
	"context"
	"errors"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
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

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Business validation happens in BeforeCreate hook
		if err := tx.Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return shared.ErrUserExists
			}
			return err
		}
		return nil
	})

	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Business validation happens in BeforeUpdate hook
		return tx.Save(user).Error
	})
}

func (r *userRepository) UpdateLastActiveAt(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Exec(
		"UPDATE users SET last_active_at = ? WHERE id = ?",
		time.Now().UTC(),
		userID,
	).Error
}

// QUERY OPERATIONS (Read)

func (r *userRepository) FindByID(ctx context.Context, userID uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("id", "username", "email", "last_active_at", "status").
		First(&user, userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("id", "username", "email", "password_hash", "last_active_at").
		Where("username = ?", username).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) Exists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("username = ?", username).
		Count(&count).Error

	return count > 0, err
}
