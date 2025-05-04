package message

import "time"

type MessageResponse struct {
	ID          uint      `json:"id"`
	Content     string    `json:"content"`
	MediaURL    string    `json:"media_url,omitempty"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	SenderID    uint      `json:"sender_id"`
	RecipientID uint      `json:"recipient_id,omitempty"`
	SentAt      time.Time `json:"sent_at"`
	DeliveredAt time.Time `json:"delivered_at,omitempty"`
	ReadAt      time.Time `json:"read_at,omitempty"`
}

type ConversationResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int64             `json:"total"`
}
