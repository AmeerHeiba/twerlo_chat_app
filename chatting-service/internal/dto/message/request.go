package message

import "time"

type SendRequest struct {
	Content     string `json:"content" validate:"required_without=MediaURL"`
	MediaURL    string `json:"media_url" validate:"omitempty,url"`
	RecipientID uint   `json:"recipient_id" validate:"required_if=Type direct"`
	Type        string `json:"type" validate:"required,oneof=direct broadcast"`
}

type BroadcastRequest struct {
	Content      string `json:"content" validate:"required_without=MediaURL"`
	MediaURL     string `json:"media_url" validate:"omitempty,url"`
	RecipientIDs []uint `json:"recipient_ids" validate:"required,min=1"`
}

type QueryRequest struct {
	Limit       int       `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset      int       `json:"offset" validate:"omitempty,min=0"`
	Before      time.Time `json:"before"`
	After       time.Time `json:"after"`
	MessageType string    `json:"message_type" validate:"omitempty,oneof=direct broadcast"`
	HasMedia    *bool     `json:"has_media"`
	Status      string    `json:"status" validate:"omitempty,oneof=sent delivered read"`
}
