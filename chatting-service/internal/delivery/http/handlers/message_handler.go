package handlers

import (
	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/dto/message"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
		shared.Log.Error("Invalid request body", zap.Error(err), zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body failed to parse request body").WithDetails(err.Error())
	}

	// Additional validation
	if body.Type == "direct" && body.RecipientID == 0 {
		shared.Log.Debug("Missing recipient ID ",
			zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body check recipient ID")
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
		shared.Log.Warn("Wrong broadcast endpoint", zap.String("url", c.OriginalURL()))
		return shared.ErrBadRequest.WithDetails("Invalid request body use broadcast endpoint for sending broadcast messages")
	default:
		shared.Log.Warn("Invalid message type", zap.String("type", body.Type))
		return shared.ErrBadRequest.WithDetails("Invalid request body check message type")
	}

	if err != nil {
		shared.Log.Warn("Failed to send message", zap.Error(err), zap.ByteString("body", c.Body()))
		return err
	}

	return c.JSON(toMessageResponse(msg))
}

func (h *MessageHandler) SendBroadcast(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	var body message.BroadcastRequest
	if err := c.BodyParser(&body); err != nil {
		shared.Log.Error("Invalid request body", zap.Error(err), zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body failed to parse request body").WithDetails(err.Error())
	}

	// validation for recipient IDs
	if len(body.RecipientIDs) == 0 {
		shared.Log.Debug("No recipients", zap.ByteString("body", c.Body()))
		return shared.ErrBadRequest.WithDetails("Invalid request body check recipient IDs")
	}

	msg, err := h.messageService.SendBroadcast(
		c.Context(),
		claims.UserID,
		body.Content,
		body.MediaURL,
		body.RecipientIDs,
	)
	if err != nil {
		shared.Log.Error("Failed to send broadcast message", zap.Error(err), zap.ByteString("body", c.Body()))
		return err
	}

	return c.JSON(toMessageResponse(msg))
}

func (h *MessageHandler) GetConversation(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)
	otherUserID, err := c.ParamsInt("userID")
	if err != nil {
		shared.Log.Error("Invalid user ID", zap.Error(err), zap.String("userID", c.Params("userID")))
		return shared.ErrBadRequest.WithDetails("Invalid or missing user ID").WithDetails(err.Error())
	}

	var query message.QueryRequest
	if err := c.QueryParser(&query); err != nil {
		shared.Log.Error("Invalid conversation request query", zap.Error(err), zap.String("path", c.Path()),
			zap.Any("query", map[string]interface{}{
				"limit":  c.Query("limit"),
				"offset": c.Query("offset"),
				"before": c.Query("before"),
			}))
		return shared.ErrBadRequest.WithDetails("Invalid or missing query params").WithDetails(err.Error())
	}

	messages, err := h.messageService.GetConversation(
		c.Context(),
		claims.UserID,
		uint(otherUserID),
		toDomainQuery(query),
	)
	if err != nil {
		shared.Log.Error("Failed to get conversation", zap.Error(err))
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
		shared.Log.Error("Invalid message ID", zap.Error(err))
		return shared.ErrBadRequest.WithDetails("Invalid or missing message ID").WithDetails(err.Error())
	}

	if err := h.messageService.MarkAsRead(
		c.Context(),
		uint(messageID),
		claims.UserID,
	); err != nil {
		shared.Log.Error("Failed to mark message as read", zap.Error(err))
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
		shared.Log.Error("Invalid message ID", zap.Error(err), zap.String("id", c.Params("id")))
		return shared.ErrBadRequest.WithDetails("Invalid or missing message ID").WithDetails(err.Error())
	}

	if err := h.messageService.DeleteMessage(
		c.Context(),
		uint(messageID),
		claims.UserID,
	); err != nil {
		shared.Log.Error("Failed to delete message", zap.Error(err), zap.String("id", c.Params("id")))
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

func (h *MessageHandler) GetLoggedInUserConversations(c *fiber.Ctx) error {
	claims := c.Locals("userClaims").(*domain.TokenClaims)

	messages, err := h.messageService.GetMessageHistory(
		c.Context(),
		claims.UserID,
		domain.MessageQuery{},
	)
	if err != nil {
		shared.Log.Error("Failed to get conversation", zap.Error(err))
		return err
	}

	return c.JSON(messages)
}
