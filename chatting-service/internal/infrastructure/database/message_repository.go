package database

import (
	"context"
	"errors"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
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
	return msg, err
}

func (r *messageRepository) CreateWithRecipients(ctx context.Context, msg *domain.Message, recipientIDs []uint) (*domain.Message, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
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
			return tx.Create(&recipients).Error
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
		Preload("Sender").
		Preload("Recipients").
		First(&message, messageID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrMessageNotFound
	}
	return &message, err
}

func (r *messageRepository) FindConversation(ctx context.Context, user1ID, user2ID uint, query domain.MessageQuery) ([]domain.Message, error) {
	var messages []domain.Message

	q := r.db.WithContext(ctx).
		Preload("Sender").
		Where("((sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?))",
			user1ID, user2ID, user2ID, user1ID).
		Where("deleted_at IS NULL")

	q = applyMessageQuery(q, query)

	err := q.Find(&messages).Error
	return messages, err
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
	return messages, err
}

func (r *messageRepository) FindBroadcasts(ctx context.Context, broadcasterID uint, query domain.MessageQuery) ([]domain.Message, error) {
	var messages []domain.Message

	q := r.db.WithContext(ctx).
		Preload("Broadcaster").
		Where("broadcaster_id = ?", broadcasterID).
		Where("deleted_at IS NULL")

	q = applyMessageQuery(q, query)

	err := q.Find(&messages).Error
	return messages, err
}

func (r *messageRepository) MarkAsDelivered(ctx context.Context, messageID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var message domain.Message
		if err := tx.First(&message, messageID).Error; err != nil {
			return err
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
		// Update message status if sender is viewing
		if err := tx.Model(&domain.Message{}).
			Where("id = ?", messageID).
			Updates(map[string]interface{}{
				"status":  domain.StatusRead,
				"read_at": time.Now().UTC(),
			}).Error; err != nil {
			return err
		}

		// Update recipient status for broadcasts
		return tx.Model(&domain.MessageRecipient{}).
			Where("message_id = ? AND user_id = ?", messageID, recipientID).
			Update("read_at", time.Now().UTC()).Error
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
			return nil // No updates needed
		}

		return tx.Model(&domain.Message{}).
			Where("id = ?", messageID).
			Updates(updates).Error
	})
}

func (r *messageRepository) Delete(ctx context.Context, messageID uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Message{}, messageID).Error
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
