package domain

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content     string        `gorm:"type:text"`
	MediaURL    string        `gorm:"type:varchar(255)"`
	MessageType MessageType   `gorm:"type:message_type;default:'direct'"`
	Status      MessageStatus `gorm:"type:message_status;default:'sent'"`

	// Relationships
	//Using foreign keys and gorm models to allow eager/lazy loading
	SenderID uint `gorm:"index"` // Foreign key to User
	Sender   User `gorm:"foreignKey:SenderID"`

	RecipientID *uint `gorm:"index;null"` // Null for broadcasts
	Recipient   *User `gorm:"foreignKey:RecipientID"`

	BroadcasterID *uint `gorm:"index;null"` // For broadcast origin
	Broadcaster   *User `gorm:"foreignKey:BroadcasterID"`

	// For broadcast recipients (many-to-many)
	Recipients []User `gorm:"many2many:message_recipients;"`

	// Metadata
	SentAt      time.Time `gorm:"index;default:CURRENT_TIMESTAMP"`
	DeliveredAt *time.Time
	ReadAt      *time.Time
}

// MessageRecipient join table for broadcasts
type MessageRecipient struct {
	MessageID  uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"primaryKey"`
	ReceivedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ReadAt     *time.Time
}
