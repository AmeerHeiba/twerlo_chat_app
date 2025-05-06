package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupWebSocketRoutes(app *fiber.App, wsHandler *handlers.WebSocketHandler, authMiddleware fiber.Handler) {
	app.Get("/ws",
		authMiddleware,
		wsHandler.Upgrade,
		websocket.New(wsHandler.HandleConnection),
	)
}
