package integration

import (
	"context"
	"testing"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/database"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/realtime"
	"github.com/stretchr/testify/assert"
)

func TestMessageFlow(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	messageRepo := database.NewMessageRepository(db)
	messageRecipientRepo := database.NewMessageRecipientRepository(db)

	// Initialize services
	notifier := realtime.NewWebSocketNotifier()
	messageService := application.NewMessageService(
		messageRepo,
		messageRecipientRepo,
		userRepo,
		notifier,
		nil, // no media uploader for this test
	)

	// Create test users
	sender, err := userRepo.Create(context.Background(), "sender", "sender@test.com", "password")
	assert.NoError(t, err)

	recipient, err := userRepo.Create(context.Background(), "recipient", "recipient@test.com", "password")
	assert.NoError(t, err)

	// Test sending a message
	msg, err := messageService.SendDirectMessage(
		context.Background(),
		sender.ID,
		recipient.ID,
		"Hello",
		"",
	)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", msg.Content)

	// Test retrieving conversation
	messages, err := messageService.GetConversation(
		context.Background(),
		sender.ID,
		recipient.ID,
		domain.MessageQuery{Limit: 10},
	)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "Hello", messages[0].Content)

	// Test marking as read
	err = messageService.MarkAsRead(context.Background(), msg.ID, recipient.ID)
	assert.NoError(t, err)

	// Verify status updated
	updatedMsg, err := messageRepo.FindByID(context.Background(), msg.ID)
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusRead, updatedMsg.Status)
	assert.NotNil(t, updatedMsg.ReadAt)
}
