package database

import (
	"context"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"gorm.io/gorm"
)

type messageRecipientRepository struct {
	db *gorm.DB
}

func NewMessageRecipientRepository(db *gorm.DB) domain.MessageRecipientRepository {
	return &messageRecipientRepository{db: db}
}

func (r *messageRecipientRepository) Create(ctx context.Context, messageID uint, recipientID uint) error {
	return r.db.WithContext(ctx).Create(&domain.MessageRecipient{
		MessageID:  messageID,
		UserID:     recipientID,
		ReceivedAt: time.Now().UTC(),
	}).Error
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

	return r.db.WithContext(ctx).Create(&recipients).Error
}
