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
	recipients := make([]domain.MessageRecipient, len(recipientIDs))
	for i, id := range recipientIDs {
		recipients[i] = domain.MessageRecipient{
			MessageID:  messageID,
			UserID:     id,
			ReceivedAt: time.Now().UTC(),
		}
	}

	err := r.db.WithContext(ctx).Create(&recipients).Error
	if err != nil {
		shared.Log.Error("create bulk message recipients failed",
			zap.String("operation", "CreateBulk"),
			zap.Uint("messageID", messageID),
			zap.Int("recipientCount", len(recipientIDs)),
			zap.Error(err))
		return shared.ErrDatabaseOperation.WithDetails("create bulk message recipients failed").WithDetails(err.Error())
	}

	shared.Log.Debug("bulk message recipients created",
		zap.Uint("messageID", messageID),
		zap.Int("recipientCount", len(recipientIDs)))
	return nil
}
