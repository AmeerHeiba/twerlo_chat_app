package domain

import (
	"io"
	"time"
)

type MediaUpload struct {
	File        io.Reader `validate:"required"`
	Filename    string    `validate:"required,min=1,max=255"`
	ContentType string    `validate:"required,oneof=image/jpeg image/png application/pdf"`
	Size        int64     `validate:"required,min=1,max=10485760"` // 10MB max
	UserID      uint      `validate:"required"`
}
type MediaResponse struct {
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"` // For temporary URLs "future enhancment"
}
