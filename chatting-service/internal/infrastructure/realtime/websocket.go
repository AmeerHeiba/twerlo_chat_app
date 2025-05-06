package realtime

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WebSocketNotifier struct {
	clients   map[uint]*ConnectionWrapper
	clientsMu sync.Mutex
	logger    *zap.Logger
}

func NewWebSocketNotifier() *WebSocketNotifier {
	return &WebSocketNotifier{
		clients: make(map[uint]*ConnectionWrapper),
		logger:  shared.Log,
	}
}

func (w *WebSocketNotifier) Notify(ctx context.Context, message *domain.Message) error {
	if message == nil {
		return errors.New("nil message")
	}

	if message.RecipientID == nil {
		return errors.New("message has no recipient")
	}

	w.clientsMu.Lock()
	conn, ok := w.clients[*message.RecipientID]
	w.clientsMu.Unlock()

	if !ok {
		w.logger.Debug("recipient not connected",
			zap.Uint("recipientID", *message.RecipientID))
		return nil
	}

	wsMessage := struct {
		ID          uint      `json:"id"`
		Content     string    `json:"content"`
		MessageType string    `json:"message_type"`
		Status      string    `json:"status"`
		SenderID    uint      `json:"sender_id"`
		RecipientID uint      `json:"recipient_id"`
		SentAt      time.Time `json:"sent_at"`
	}{
		ID:          message.ID,
		Content:     message.Content,
		MessageType: string(message.MessageType),
		Status:      string(message.Status),
		SenderID:    message.SenderID,
		RecipientID: *message.RecipientID,
		SentAt:      message.SentAt,
	}

	w.logger.Debug("sending websocket message",
		zap.Uint("messageID", message.ID),
		zap.Uint("recipientID", *message.RecipientID))

	return conn.WriteJSON(wsMessage)
}

func (w *WebSocketNotifier) Broadcast(ctx context.Context, message *domain.Message, recipientIDs []uint) error {
	w.clientsMu.Lock()
	defer w.clientsMu.Unlock()

	for _, id := range recipientIDs {
		if conn, ok := w.clients[id]; ok {
			if err := conn.WriteJSON(message); err != nil {
				zap.L().Error("websocket broadcast failed",
					zap.Uint("userID", id),
					zap.Error(err))
			}
		}
	}
	return nil
}

func (w *WebSocketNotifier) RegisterClient(userID uint, conn *websocket.Conn) {
	w.clientsMu.Lock()
	defer w.clientsMu.Unlock()
	w.clients[userID] = &ConnectionWrapper{conn: conn}
	shared.Log.Info("WebSocket client registered",
		zap.Uint("userID", userID),
		zap.Int("activeConnections", len(w.clients)))
}

func (w *WebSocketNotifier) RemoveClient(userID uint) {
	w.clientsMu.Lock()
	defer w.clientsMu.Unlock()
	delete(w.clients, userID)
	shared.Log.Info("WebSocket client removed",
		zap.Uint("userID", userID),
		zap.Int("activeConnections", len(w.clients)))
}

func (w *WebSocketNotifier) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		shared.Log.Debug("WebSocket upgrade requested",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()))
		return c.Next()
	}
	shared.Log.Warn("WebSocket upgrade failed - not a WebSocket request")
	return fiber.ErrUpgradeRequired
}

func (w *WebSocketNotifier) HandleConnection(conn *websocket.Conn) {
	userID := conn.Locals("userID").(uint)

	// Register connection
	wrapper := &ConnectionWrapper{conn: conn}
	w.clientsMu.Lock()
	w.clients[userID] = wrapper
	w.clientsMu.Unlock()

	w.logger.Info("websocket connection established",
		zap.Uint("userID", userID),
		zap.String("remoteAddr", conn.RemoteAddr().String()))

	defer func() {
		w.clientsMu.Lock()
		delete(w.clients, userID)
		w.clientsMu.Unlock()
		conn.Close()
		w.logger.Info("websocket connection closed",
			zap.Uint("userID", userID))
	}()

	// Configure connection
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Heartbeat
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
					return
				}
			}
		}
	}()

	// Keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (w *WebSocketNotifier) Subscribe(ctx context.Context, userID uint) (<-chan *domain.Message, error) {
	// Return a channel that will never receive messages
	return nil, nil
}

func (w *WebSocketNotifier) Unsubscribe(ctx context.Context, userID uint) error {
	// No-op since i am not actually maintaining subscriptions
	return nil
}
