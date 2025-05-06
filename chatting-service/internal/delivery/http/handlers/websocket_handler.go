package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/realtime"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocketHandler struct {
	notifier *realtime.WebSocketNotifier
}

func NewWebSocketHandler(notifier *realtime.WebSocketNotifier) *WebSocketHandler {
	return &WebSocketHandler{notifier: notifier}
}

func (h *WebSocketHandler) Upgrade(c *fiber.Ctx) error {
	return h.notifier.Upgrade(c)
}

func (h *WebSocketHandler) HandleConnection(conn *websocket.Conn) {
	h.notifier.HandleConnection(conn)
}
