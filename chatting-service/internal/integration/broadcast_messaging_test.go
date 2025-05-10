package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/database"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/realtime"
	"github.com/stretchr/testify/assert"
)

func TestBroadcastMessaging(t *testing.T) {
	db := setupTestDB(t)
	userRepo := database.NewUserRepository(db)
	messageRepo := database.NewMessageRepository(db)
	messageRecipientRepo := database.NewMessageRecipientRepository(db)
	notifier := realtime.NewWebSocketNotifier()

	messageService := application.NewMessageService(
		messageRepo,
		messageRecipientRepo,
		userRepo,
		notifier,
		nil,
	)

	// Create test users
	broadcaster, err := userRepo.Create(context.Background(), "broadcaster", "broadcaster@test.com", "password")
	assert.NoError(t, err)

	var recipients []uint
	for i := 0; i < 3; i++ {
		user, err := userRepo.Create(context.Background(),
			fmt.Sprintf("recipient%d", i),
			fmt.Sprintf("recipient%d@test.com", i),
			"password")
		assert.NoError(t, err)
		recipients = append(recipients, user.ID)
	}

	// Send broadcast
	broadcastMsg, err := messageService.SendBroadcast(
		context.Background(),
		broadcaster.ID,
		"Important announcement!",
		"",
		recipients,
	)
	assert.NoError(t, err)
	assert.Equal(t, domain.MessageBroadcast, broadcastMsg.MessageType)

	// Verify recipients
	msgFromDB, err := messageRepo.FindByID(context.Background(), broadcastMsg.ID)
	assert.NoError(t, err)
	assert.Len(t, msgFromDB.Recipients, 3)
}
