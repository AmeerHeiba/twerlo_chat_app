package application

import (
	"context"
	"errors"
	"strings"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
)

type MessageService struct {
	messageRepo          domain.MessageRepository
	messageRecipientRepo domain.MessageRecipientRepository
	userRepo             domain.UserRepository
	notifier             domain.MessageNotifier // Optional for real-time
}

func NewMessageService(
	messageRepo domain.MessageRepository,
	messageRecipientRepo domain.MessageRecipientRepository,
	userRepo domain.UserRepository,
	notifier domain.MessageNotifier,
) *MessageService {
	return &MessageService{
		messageRepo:          messageRepo,
		messageRecipientRepo: messageRecipientRepo,
		userRepo:             userRepo,
		notifier:             notifier,
	}
}

func (s *MessageService) SendDirectMessage(ctx context.Context, senderID, recipientID uint, content string, mediaURL string) (*domain.Message, error) {
	// Validate recipient exists
	if _, err := s.userRepo.FindByID(ctx, recipientID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, shared.ErrInvalidRecipient
		}
		return nil, err
	}

	// Validate content
	if strings.TrimSpace(content) == "" && mediaURL == "" {
		return nil, shared.ErrEmptyMessage
	}

	msg, err := s.messageRepo.Create(ctx, senderID, content, mediaURL, domain.MessageDirect)
	if err != nil {
		shared.Log.Error("Failed to create message",
			zap.Uint("sender", senderID),
			zap.Uint("recipient", recipientID),
			zap.Error(err))
		return nil, err
	}

	// Update recipient association
	if err := s.messageRepo.Update(ctx, msg.ID, &recipientID, nil); err != nil {
		return nil, err
	}

	// Notify recipient (if notifier is configured)
	if s.notifier != nil {
		if err := s.notifier.Notify(ctx, msg); err != nil {
			shared.Log.Warn("Failed to notify recipient",
				zap.Uint("message", msg.ID),
				zap.Error(err))
		}
	}

	return msg, nil
}

func (s *MessageService) SendBroadcast(ctx context.Context, broadcasterID uint, content string, mediaURL string, recipientIDs []uint) (*domain.Message, error) {
	if len(recipientIDs) == 0 {
		return nil, shared.ErrNoRecipients
	}

	// Validate all recipients exist first
	for _, id := range recipientIDs {
		if _, err := s.userRepo.FindByID(ctx, id); err != nil {
			return nil, shared.ErrUserNotFound
		}
	}

	// Create message object (validation happens in CreateWithRecipients)
	msg := &domain.Message{
		SenderID:    broadcasterID,
		Content:     content,
		MediaURL:    mediaURL,
		MessageType: domain.MessageBroadcast,
		Status:      domain.StatusSent,
	}

	// Create message with recipients in one transaction
	createdMsg, err := s.messageRepo.CreateWithRecipients(ctx, msg, recipientIDs)
	if err != nil {
		shared.Log.Error("Failed to create broadcast",
			zap.Uint("broadcaster", broadcasterID),
			zap.Error(err))
		return nil, err
	}

	// Reload the message with all relationships
	fullMessage, err := s.messageRepo.FindByID(ctx, createdMsg.ID) // Use createdMsg.ID here
	if err != nil {
		return nil, err
	}

	// Notify recipients
	if s.notifier != nil {
		if err := s.notifier.Broadcast(ctx, fullMessage, recipientIDs); err != nil {
			shared.Log.Warn("Failed to notify broadcast recipients",
				zap.Uint("message", fullMessage.ID),
				zap.Error(err))
		}
	}

	return fullMessage, nil
}
func (s *MessageService) GetConversation(ctx context.Context, user1ID, user2ID uint, query domain.MessageQuery) ([]domain.Message, error) {
	// Validate both users exist
	if _, err := s.userRepo.FindByID(ctx, user1ID); err != nil {
		return nil, shared.ErrUserNotFound
	}
	if _, err := s.userRepo.FindByID(ctx, user2ID); err != nil {
		return nil, shared.ErrUserNotFound
	}

	messages, err := s.messageRepo.FindConversation(ctx, user1ID, user2ID, query)
	if err != nil {
		shared.Log.Error("Failed to get conversation",
			zap.Uint("user1", user1ID),
			zap.Uint("user2", user2ID),
			zap.Error(err))
		return nil, err
	}

	return messages, nil
}

func (s *MessageService) GetMessageHistory(ctx context.Context, userID uint, query domain.MessageQuery) ([]domain.Message, error) {
	messages, err := s.messageRepo.FindUserMessages(ctx, userID, query)
	if err != nil {
		shared.Log.Error("Failed to get message history",
			zap.Uint("user", userID),
			zap.Error(err))
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) MarkAsDelivered(ctx context.Context, messageID uint) error {
	if err := s.messageRepo.MarkAsDelivered(ctx, messageID); err != nil {
		shared.Log.Error("Failed to mark message as delivered",
			zap.Uint("message", messageID),
			zap.Error(err))
		return err
	}
	return nil
}

func (s *MessageService) MarkAsRead(ctx context.Context, messageID uint, recipientID uint) error {
	if err := s.messageRepo.MarkAsRead(ctx, messageID, recipientID); err != nil {
		shared.Log.Error("Failed to mark message as read",
			zap.Uint("message", messageID),
			zap.Uint("recipient", recipientID),
			zap.Error(err))
		return err
	}
	return nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID uint, userID uint) error {
	// Verify user has permission to delete (either sender or recipient)
	msg, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	if msg.SenderID != userID && (msg.RecipientID != nil && *msg.RecipientID != userID) {
		return shared.ErrInvalidCredentials
	}

	if err := s.messageRepo.Delete(ctx, messageID); err != nil {
		shared.Log.Error("Failed to delete message",
			zap.Uint("message", messageID),
			zap.Uint("user", userID),
			zap.Error(err))
		return err
	}

	return nil
}
