package database

import (
	"context"
	"errors"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, senderID uint, content, mediaURL string, messageType domain.MessageType) (*domain.Message, error) {
	msg := &domain.Message{
		SenderID:    senderID,
		Content:     content,
		MediaURL:    mediaURL,
		MessageType: messageType,
		Status:      domain.StatusSent,
	}

	err := r.db.WithContext(ctx).Create(msg).Error
	if err != nil {
		shared.Log.Error("create message failed",
			zap.String("operation", "Create"),
			zap.Uint("senderID", senderID),
			zap.String("content", content),
			zap.String("mediaURL", mediaURL),
			zap.String("messageType", string(messageType)),
			zap.Error(err))
		return nil, shared.ErrDatabaseOperation.WithDetails("create message failed").WithDetails(err.Error())
	}
	return msg, nil
}

func (r *messageRepository) CreateWithRecipients(ctx context.Context, msg *domain.Message, recipientIDs []uint) (*domain.Message, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return shared.ErrDatabaseOperation.WithDetails("create message failed")
		}

		if len(recipientIDs) > 0 {
			recipients := make([]domain.MessageRecipient, len(recipientIDs))
			for i, id := range recipientIDs {
				recipients[i] = domain.MessageRecipient{
					MessageID:  msg.ID,
					UserID:     id,
					ReceivedAt: time.Now().UTC(),
				}
			}
			if err := tx.Create(&recipients).Error; err != nil {
				return shared.ErrDatabaseOperation.WithDetails("create recipients failed")
			}
		}
		return nil
	})
	return msg, err
}

func (r *messageRepository) CreateWithTransaction(ctx context.Context, fn func(ctx context.Context, txRepo domain.MessageRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new repository instance with the transaction DB
		txRepo := NewMessageRepository(tx)
		return fn(ctx, txRepo)
	})
}

func (r *messageRepository) FindByID(ctx context.Context, messageID uint) (*domain.Message, error) {
	var message domain.Message
	err := r.db.WithContext(ctx).
		Preload("Sender", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email", "status", "last_active_at")
		}).
		Preload("Recipients", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email", "status", "last_active_at")
		}).
		First(&message, messageID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		shared.Log.Debug("message not found",
			zap.Uint("messageID", messageID),
			zap.Error(err))
		return nil, shared.ErrRecordNotFound.WithDetails("message not found").WithDetails(err.Error())
	}
	if err != nil {
		shared.Log.Error("find message by ID failed",
			zap.String("operation", "FindByID"),
			zap.Uint("messageID", messageID),
			zap.Error(err))
		return nil, shared.ErrDatabaseOperation.WithDetails("find message by ID failed").WithDetails(err.Error())
	}
	return &message, nil
}

func (r *messageRepository) FindConversation(ctx context.Context, user1ID, user2ID uint, query domain.MessageQuery) ([]domain.Message, error) {
	var messages []domain.Message

	q := r.db.WithContext(ctx).
		Joins("Sender").
		Where("((sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?))",
			user1ID, user2ID, user2ID, user1ID)

	q = applyMessageQuery(q, query)

	if err := q.Find(&messages).Error; err != nil {
		return nil, shared.ErrDatabaseOperation.WithDetails("find conversation failed")
	}
	return messages, nil
}

func (r *messageRepository) FindUserMessages(ctx context.Context, userID uint, query domain.MessageQuery) ([]domain.Message, error) {
	var messages []domain.Message

	q := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Recipient").
		Where("sender_id = ? OR recipient_id = ?", userID, userID).
		Where("deleted_at IS NULL")

	q = applyMessageQuery(q, query)

	err := q.Find(&messages).Error
	if err != nil {
		shared.Log.Error("find user messages failed",
			zap.String("operation", "FindUserMessages"),
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, shared.ErrDatabaseOperation.WithDetails("find user messages failed").WithDetails(err.Error())
	}
	return messages, nil
}

func (r *messageRepository) FindBroadcasts(ctx context.Context, broadcasterID uint, query domain.MessageQuery) ([]domain.Message, error) {
	var messages []domain.Message

	q := r.db.WithContext(ctx).
		Preload("Broadcaster").
		Where("broadcaster_id = ?", broadcasterID).
		Where("deleted_at IS NULL")

	q = applyMessageQuery(q, query)

	err := q.Find(&messages).Error
	if err != nil {
		shared.Log.Error("find broadcasts failed",
			zap.String("operation", "FindBroadcasts"),
			zap.Uint("broadcasterID", broadcasterID),
			zap.Error(err))
		return nil, shared.ErrDatabaseOperation.WithDetails("find broadcasts failed").WithDetails(err.Error())
	}
	return messages, nil
}

func (r *messageRepository) MarkAsDelivered(ctx context.Context, messageID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var message domain.Message
		if err := tx.First(&message, messageID).Error; err != nil {
			shared.Log.Error("mark message as delivered failed",
				zap.String("operation", "MarkAsDelivered"),
				zap.Uint("messageID", messageID),
				zap.Error(err))
			return shared.ErrDatabaseOperation.WithDetails("mark message as delivered failed").WithDetails(err.Error())
		}

		now := time.Now().UTC()
		return tx.Model(&message).
			Updates(map[string]interface{}{
				"status":       domain.StatusDelivered,
				"delivered_at": now,
			}).Error
	})
}

func (r *messageRepository) MarkAsRead(ctx context.Context, messageID uint, recipientID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		// Update main message status
		if err := tx.Model(&domain.Message{}).
			Where("id = ?", messageID).
			Updates(map[string]interface{}{
				"status":  domain.StatusRead,
				"read_at": now,
			}).Error; err != nil {
			return shared.ErrBadRequest.WithDetails("failed to update msg status").WithDetails(err.Error())
		}

		// Update recipient status if exists
		return tx.Model(&domain.MessageRecipient{}).
			Where("message_id = ? AND user_id = ?", messageID, recipientID).
			Update("read_at", now).Error
	})
}

func (r *messageRepository) Update(ctx context.Context, messageID uint, recipientID *uint, broadcasterID *uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := make(map[string]interface{})

		if recipientID != nil {
			updates["recipient_id"] = recipientID
			// Clear broadcaster if setting recipient (direct message)
			updates["broadcaster_id"] = nil
		}

		if broadcasterID != nil {
			updates["broadcaster_id"] = broadcasterID
			// Clear recipient if setting broadcaster (broadcast)
			updates["recipient_id"] = nil
		}

		if len(updates) == 0 {
			shared.Log.Debug("no updates needed",
				zap.String("operation", "Update"),
				zap.Uint("messageID", messageID))
			return nil // No updates needed
		}

		err := tx.Model(&domain.Message{}).
			Where("id = ?", messageID).
			Updates(updates).Error
		if err != nil {
			shared.Log.Error("update message failed",
				zap.String("operation", "Update"),
				zap.Uint("messageID", messageID),
				zap.Error(err))
			return shared.ErrDatabaseOperation.WithDetails("update message failed").WithDetails(err.Error())
		}
		return nil
	})
}

func (r *messageRepository) Delete(ctx context.Context, messageID uint, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verify ownership
		var count int64
		if err := tx.Model(&domain.Message{}).
			Where("id = ? AND (sender_id = ? OR recipient_id = ?)",
				messageID, userID, userID).
			Count(&count).Error; err != nil {
			return shared.ErrBadRequest.WithDetails("message ownership check failed").WithDetails(err.Error())
		}

		if count == 0 {
			return shared.ErrForbidden.WithDetails("message not owned by user")
		}

		// Delete message
		return tx.Delete(&domain.Message{}, messageID).Error
	})
}

// Helper function to apply query filters
func applyMessageQuery(q *gorm.DB, query domain.MessageQuery) *gorm.DB {
	if query.Limit > 0 {
		q = q.Limit(query.Limit)
	}
	if query.Offset > 0 {
		q = q.Offset(query.Offset)
	}
	if !query.Before.IsZero() {
		q = q.Where("sent_at < ?", query.Before)
	}
	if !query.After.IsZero() {
		q = q.Where("sent_at > ?", query.After)
	}
	if query.MessageType != "" {
		q = q.Where("message_type = ?", query.MessageType)
	}
	if query.HasMedia != nil {
		if *query.HasMedia {
			q = q.Where("media_url IS NOT NULL AND media_url != ''")
		} else {
			q = q.Where("media_url IS NULL OR media_url = ''")
		}
	}
	if query.Status != "" {
		q = q.Where("status = ?", query.Status)
	}

	// Default sorting - newest first
	sortOrder := "DESC"
	if query.SortBy == "asc" {
		sortOrder = "ASC"
	}
	q = q.Order("sent_at " + sortOrder)

	return q
}
