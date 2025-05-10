package integration

import (
	"context"
	"testing"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/database"
	"github.com/stretchr/testify/assert"
)

func TestErrorScenarios(t *testing.T) {
	db := setupTestDB(t)
	userRepo := database.NewUserRepository(db)
	messageRepo := database.NewMessageRepository(db)
	messageRecipientRepo := database.NewMessageRecipientRepository(db)

	messageService := application.NewMessageService(
		messageRepo,
		messageRecipientRepo,
		userRepo,
		nil,
		nil,
	)

	t.Run("SendToNonexistentUser", func(t *testing.T) {
		sender, err := userRepo.Create(context.Background(), "sender1", "sender1@test.com", "password")
		assert.NoError(t, err)

		_, err = messageService.SendDirectMessage(
			context.Background(),
			sender.ID,
			9999, // non-existent user
			"Hello",
			"",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("EmptyMessage", func(t *testing.T) {
		sender, err := userRepo.Create(context.Background(), "sender2", "sender2@test.com", "password")
		assert.NoError(t, err)
		recipient, err := userRepo.Create(context.Background(), "recipient2", "recipient2@test.com", "password")
		assert.NoError(t, err)

		_, err = messageService.SendDirectMessage(
			context.Background(),
			sender.ID,
			recipient.ID,
			"", // empty message
			"", // no media
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid or empty message content")
	})
}
