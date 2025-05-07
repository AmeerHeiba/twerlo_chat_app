package database

import (
	"context"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type messageRecipientRepository struct {
	db *gorm.DB
}

func NewMessageRecipientRepository(db *gorm.DB) domain.MessageRecipientRepository {
	return &messageRecipientRepository{db: db}
}

func (r *messageRecipientRepository) Create(ctx context.Context, messageID uint, recipientID uint) error {
	err := r.db.WithContext(ctx).Create(&domain.MessageRecipient{
		MessageID:  messageID,
		UserID:     recipientID,
		ReceivedAt: time.Now().UTC(),
	}).Error

	if err != nil {
		shared.Log.Error("create message recipient failed",
			zap.String("operation", "Create"),
			zap.Uint("messageID", messageID),
			zap.Uint("recipientID", recipientID),
			zap.Error(err))
		return shared.ErrDatabaseOperation.WithDetails("create message recipient failed").WithDetails(err.Error())
	}

	shared.Log.Debug("message recipient created",
		zap.Uint("messageID", messageID),
		zap.Uint("recipientID", recipientID))
	return nil
}

func (r *messageRecipientRepository) CreateBulk(ctx context.Context, messageID uint, recipientIDs []uint) error {
	if len(recipientIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate recipients exist
		var count int64
		if err := tx.Model(&domain.User{}).
			Where("id IN ?", recipientIDs).
			Count(&count).Error; err != nil {
			return shared.ErrDatabaseOperation.WithDetails("recipient validation failed").WithDetails(err.Error())
		}

		if count != int64(len(recipientIDs)) {
			return domain.ErrInvalidRecipientID
		}

		// Create recipients
		recipients := make([]domain.MessageRecipient, len(recipientIDs))
		now := time.Now().UTC()
		for i, id := range recipientIDs {
			recipients[i] = domain.MessageRecipient{
				MessageID:  messageID,
				UserID:     id,
				ReceivedAt: now,
			}
		}

		return tx.CreateInBatches(recipients, 100).Error // Batch insert for performance
	})
}
