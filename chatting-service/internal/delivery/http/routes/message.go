package routes

import (
	"github.com/AmeerHeiba/chatting-service/internal/delivery/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupMessageRoutes(app *fiber.App, handler *handlers.MessageHandler, wsHandler *handlers.WebSocketHandler, authMiddleware fiber.Handler) {
	messageGroup := app.Group("/api/messages", authMiddleware)

	messageGroup.Post("/", handler.SendMessage)
	messageGroup.Post("/broadcast", handler.SendBroadcast)
	messageGroup.Get("/conversations", handler.GetLoggedInUserConversations)
	messageGroup.Get("/conversation/:userID", handler.GetConversation)
	messageGroup.Put("/:id/read", handler.MarkAsRead)
	messageGroup.Delete("/:id", handler.DeleteMessage)

	// Todo get all conversation for signed in user

}
