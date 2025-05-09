package realtime

import (
	"context"
	"testing"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/contrib/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketNotifier(t *testing.T) {
	shared.InitLogger("test")

	t.Run("Register and Notify", func(t *testing.T) {
		notifier := NewWebSocketNotifier()

		// Mock WebSocket connection
		conn := &websocket.Conn{}
		userID := uint(1)

		// Register client
		notifier.RegisterClient(userID, conn)

		// Test notification
		msg := &domain.Message{

			Content:     "Test",
			MessageType: domain.MessageDirect,
			Status:      domain.StatusSent,
			SenderID:    2,
			RecipientID: &userID,
			SentAt:      time.Now(),
		}

		err := notifier.Notify(context.Background(), msg)
		assert.NoError(t, err)

		// Cleanup
		notifier.RemoveClient(userID)
	})

}
