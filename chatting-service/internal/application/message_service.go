package application

import (
	"context"
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
	MediaService         domain.MediaService    //responsible for media operations
}

func NewMessageService(
	messageRepo domain.MessageRepository,
	messageRecipientRepo domain.MessageRecipientRepository,
	userRepo domain.UserRepository,
	notifier domain.MessageNotifier,
	mediaService domain.MediaService,
) *MessageService {
	return &MessageService{
		messageRepo:          messageRepo,
		messageRecipientRepo: messageRecipientRepo,
		userRepo:             userRepo,
		notifier:             notifier,
		MediaService:         mediaService,
	}
}

func (s *MessageService) SendDirectMessage(ctx context.Context, senderID, recipientID uint, content string, mediaURL string) (*domain.Message, error) {
	// Validate recipient exists
	if exists, err := s.userRepo.Exists(ctx, recipientID); err != nil {
		shared.Log.Error("user exists check failed",
			zap.String("operation", "SendDirectMessage"),
			zap.Uint("recipientID", recipientID),
			zap.Uint("senderID", senderID),
			zap.Error(err))
		return nil, err
	} else if !exists {
		shared.Log.Error("user not found",
			zap.String("operation", "SendDirectMessage"),
			zap.Uint("recipientID", recipientID),
			zap.Uint("senderID", senderID))
		return nil, err
	}
	// Validate content
	if strings.TrimSpace(content) == "" && mediaURL == "" {
		shared.Log.Debug("Invalid or empty message content",
			zap.String("operation", "SendDirectMessage"),
			zap.Uint("recipientID", recipientID),
			zap.Uint("senderID", senderID),
			zap.String("content", content))
		return nil, shared.ErrValidation.WithDetails("Invalid or empty message content for direct message")
	}

	msg, err := s.messageRepo.Create(ctx, senderID, content, mediaURL, domain.MessageDirect)
	if err != nil {
		shared.Log.Error("create message failed",
			zap.String("operation", "SendDirectMessage"),
			zap.Uint("recipientID", recipientID),
			zap.Error(err))
		return nil, err
	}

	// Update recipient association
	if err := s.messageRepo.Update(ctx, msg.ID, &recipientID, nil); err != nil {
		shared.Log.Error("update message failed",
			zap.String("operation", "SendDirectMessage"),
			zap.Uint("recipientID", recipientID),
			zap.Uint("senderID", senderID),
			zap.Error(err))
		return nil, err
	}

	// Notify recipient (if notifier is configured)
	if s.notifier != nil {
		if err := s.notifier.Notify(ctx, msg); err != nil {
		}
	}

	return msg, nil
}

func (s *MessageService) SendBroadcast(ctx context.Context, broadcasterID uint, content string, mediaURL string, recipientIDs []uint) (*domain.Message, error) {
	if len(recipientIDs) == 0 {
		shared.Log.Warn("Invalid or empty recipient IDs", zap.Uint("broadcasterID", broadcasterID), zap.String("content", content), zap.String("mediaURL", mediaURL))
		return nil, shared.ErrBadRequest.WithDetails("Invalid or empty recipient IDs")
	}

	// Validate all recipients exist first
	for _, id := range recipientIDs {
		if _, err := s.userRepo.FindByID(ctx, id); err != nil {
			shared.Log.Error("user not found", zap.Uint("userID", id), zap.Error(err))
			return nil, err
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
		shared.Log.Error("create message with recipients failed", zap.Error(err))
		return nil, err
	}

	// Reload the message with all relationships
	fullMessage, err := s.messageRepo.FindByID(ctx, createdMsg.ID)
	if err != nil {
		shared.Log.Error("find message by ID failed", zap.Error(err))
		return nil, err
	}

	// Notify recipients
	if s.notifier != nil {
		if err := s.notifier.Broadcast(ctx, fullMessage, recipientIDs); err != nil {
		}
	}

	return fullMessage, nil
}
func (s *MessageService) GetConversation(ctx context.Context, user1ID, user2ID uint, query domain.MessageQuery) ([]domain.Message, error) {
	// Validate both users exist
	if _, err := s.userRepo.FindByID(ctx, user1ID); err != nil {
		shared.Log.Error("user not found", zap.Uint("userID", user1ID), zap.Error(err))
		return nil, err
	}
	if _, err := s.userRepo.FindByID(ctx, user2ID); err != nil {
		shared.Log.Error("user not found", zap.Uint("userID", user2ID), zap.Error(err))
		return nil, err
	}

	messages, err := s.messageRepo.FindConversation(ctx, user1ID, user2ID, query)
	if err != nil {
		shared.Log.Error("find conversation failed",
			zap.String("operation", "GetConversation"),
			zap.Uint("user1ID", user1ID),
			zap.Uint("user2ID", user2ID),
			zap.Error(err))
		return nil, err
	}

	return messages, nil
}

func (s *MessageService) GetMessageHistory(ctx context.Context, userID uint, query domain.MessageQuery) ([]domain.Message, error) {
	messages, err := s.messageRepo.FindUserMessages(ctx, userID, query)
	if err != nil {
		shared.Log.Error("find user messages failed",
			zap.String("operation", "GetMessageHistory"),
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) MarkAsDelivered(ctx context.Context, messageID uint) error {
	if err := s.messageRepo.MarkAsDelivered(ctx, messageID); err != nil {
		shared.Log.Error("mark message as delivered failed",
			zap.String("operation", "MarkAsDelivered"),
			zap.Uint("messageID", messageID),
			zap.Error(err))
		return err
	}
	return nil
}

func (s *MessageService) MarkAsRead(ctx context.Context, messageID uint, recipientID uint) error {
	if err := s.messageRepo.MarkAsRead(ctx, messageID, recipientID); err != nil {
		shared.Log.Error("mark message as read failed",
			zap.String("operation", "MarkAsRead"),
			zap.Uint("messageID", messageID),
			zap.Uint("recipientID", recipientID),
			zap.Error(err))
		return err
	}
	return nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID uint, userID uint) error {
	// Verify user has permission to delete (either sender or recipient)
	msg, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		shared.Log.Error("find message by ID failed",
			zap.String("operation", "DeleteMessage"),
			zap.Uint("messageID", messageID),
			zap.Error(err))
		return err
	}

	if msg.SenderID != userID && (msg.RecipientID != nil && *msg.RecipientID != userID) {
		shared.Log.Debug("invalid credentials",
			zap.String("operation", "DeleteMessage"),
			zap.Uint("messageID", messageID),
			zap.Uint("userID", userID))
		return shared.ErrInvalidCredentials
	}

	if err := s.messageRepo.Delete(ctx, messageID); err != nil {
		shared.Log.Error("delete message failed",
			zap.String("operation", "DeleteMessage"),
			zap.Uint("messageID", messageID),
			zap.Error(err))
		return err
	}

	return nil
}
