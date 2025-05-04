package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/message"
	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	messageService *application.MessageService
}

func NewMessageHandler(messageService *application.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	var body message.SendRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Additional validation
	if body.Type == "direct" && body.RecipientID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Recipient ID is required for direct messages",
		})
	}

	var msg *domain.Message
	var err error

	switch body.Type {
	case "direct":
		msg, err = h.messageService.SendDirectMessage(
			c.Context(),
			claims.UserID,
			body.RecipientID,
			body.Content,
			body.MediaURL,
		)
	case "broadcast":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Use /broadcast endpoint for broadcasts",
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message type",
		})
	}

	if err != nil {
		return err
	}

	return c.JSON(toMessageResponse(msg))
}

func (h *MessageHandler) SendBroadcast(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	var body message.BroadcastRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// validation for recipient IDs
	if len(body.RecipientIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one recipient is required for broadcast",
		})
	}

	msg, err := h.messageService.SendBroadcast(
		c.Context(),
		claims.UserID,
		body.Content,
		body.MediaURL,
		body.RecipientIDs,
	)
	if err != nil {
		return err
	}

	return c.JSON(toMessageResponse(msg))
}

func (h *MessageHandler) GetConversation(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)
	otherUserID, err := c.ParamsInt("userID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var query message.QueryRequest
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	messages, err := h.messageService.GetConversation(
		c.Context(),
		claims.UserID,
		uint(otherUserID),
		toDomainQuery(query),
	)
	if err != nil {
		return err
	}

	response := message.ConversationResponse{
		Messages: make([]message.MessageResponse, len(messages)),
	}
	for i, msg := range messages {
		response.Messages[i] = toMessageResponse(&msg)
	}

	return c.JSON(response)
}

func (h *MessageHandler) MarkAsRead(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)
	messageID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	if err := h.messageService.MarkAsRead(
		c.Context(),
		uint(messageID),
		claims.UserID,
	); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Message marked as read",
	})
}

func (h *MessageHandler) DeleteMessage(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)
	messageID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	if err := h.messageService.DeleteMessage(
		c.Context(),
		uint(messageID),
		claims.UserID,
	); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Message deleted",
	})
}

// Helpers
func toDomainQuery(q message.QueryRequest) domain.MessageQuery {
	return domain.MessageQuery{
		Limit:       q.Limit,
		Offset:      q.Offset,
		Before:      q.Before,
		After:       q.After,
		MessageType: q.MessageType,
		HasMedia:    q.HasMedia,
		Status:      q.Status,
	}
}

func toMessageResponse(m *domain.Message) message.MessageResponse {
	resp := message.MessageResponse{
		ID:       m.ID,
		Content:  m.Content,
		MediaURL: m.MediaURL,
		Type:     string(m.MessageType),
		Status:   string(m.Status),
		SenderID: m.SenderID,
		SentAt:   m.SentAt,
	}

	if m.RecipientID != nil {
		resp.RecipientID = *m.RecipientID
	}
	if m.DeliveredAt != nil {
		resp.DeliveredAt = *m.DeliveredAt
	}
	if m.ReadAt != nil {
		resp.ReadAt = *m.ReadAt
	}

	return resp
}
