package domain

import (
	"net/url"
	"strings"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content     string        `gorm:"type:text" json:"content"`
	MediaURL    string        `gorm:"type:varchar(255)" json:"media_url,omitempty"`
	MessageType MessageType   `gorm:"type:message_type;default:'direct'" json:"message_type"`
	Status      MessageStatus `gorm:"type:message_status;default:'sent'" json:"status"`

	// Relationships
	//Using foreign keys and gorm models to allow eager/lazy loading
	SenderID uint `gorm:"index" json:"sender_id"` // Foreign key to User
	Sender   User `gorm:"foreignKey:SenderID" json:"-"`

	RecipientID *uint `gorm:"index;null" json:"recipient_id"`
	Recipient   *User `gorm:"foreignKey:RecipientID" json:"-"`

	BroadcasterID *uint `gorm:"index;null"` // For broadcast origin
	Broadcaster   *User `gorm:"foreignKey:BroadcasterID"`

	// For broadcast recipients (many-to-many)
	Recipients []User `gorm:"many2many:message_recipients;joinForeignKey:MessageID;joinReferences:UserID"`

	// Metadata
	SentAt      time.Time  `gorm:"index;default:CURRENT_TIMESTAMP" json:"sent_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
}

// MessageRecipient join table for broadcasts
type MessageRecipient struct {
	MessageID  uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"primaryKey"`
	ReceivedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ReadAt     *time.Time
}

//Key Business Rules
//1 - Content Validation
//		-Either text or media must be present
//		-Text length limits (1000 chars)
//		-Media URL format validation
//1 - Recipient Rules
//		-Broadcasts can't have single RecipientID
//		-Direct messages must have RecipientID
//		-Broadcasts must have RecipientIDs

func (m *Message) Validate() error {
	shared.Log.Debug("Validating message",
		zap.Any("message", m),
		zap.Bool("isBroadcast", m.IsBroadcast()),
		zap.Int("numRecipients", len(m.Recipients)),
		zap.Any("recipients", m.Recipients))

	if m.RequiresRecipientsList() {
		if len(m.Recipients) == 0 {
			shared.Log.Error("Validation failed - no recipients for broadcast")
			return ErrNoRecipients
		}
	}
	// Content Validation
	if strings.TrimSpace(m.Content) == "" && m.MediaURL == "" {
		return ErrEmptyMessage
	}
	if len(m.Content) > 1000 {
		return ErrMessageTooLong
	}
	if m.MediaURL != "" {
		if _, err := url.ParseRequestURI(m.MediaURL); err != nil {
			return ErrInvalidMediaURL
		}
	}

	// Recipient Rules
	if m.RequiresRecipientsList() {
		if m.RecipientID != nil {
			return ErrInvalidBroadcast
		}
		if len(m.Recipients) == 0 {
			return ErrNoRecipients
		}
	} else {
		if m.RecipientID == nil {
			return ErrInvalidRecipient
		}
		if len(m.Recipients) > 0 {
			return ErrDirectMessageNoList
		}
	}

	return nil
}

// Message Reciption val
func (mr *MessageRecipient) Validate() error {
	if mr.MessageID == 0 || mr.UserID == 0 {
		return ErrMissingRecOrSenderID
	}
	return nil
}

// GORM Hooks
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	// Auto-set SentAt if not specified
	if m.SentAt.IsZero() {
		m.SentAt = time.Now().UTC()
	}

	// Enforce default status
	if m.Status == "" {
		m.Status = StatusSent
	}

	return nil // REMOVED THE VALIDATION CALL HERE
}

// Helper methods
func (m *Message) IsDirect() bool {
	return m.MessageType == MessageDirect
}

func (m *Message) IsBroadcast() bool {
	return m.MessageType == MessageBroadcast
}

// State management
func (m *Message) MarkDelivered() {
	now := time.Now().UTC()
	m.DeliveredAt = &now
	m.Status = StatusDelivered
}

func (m *Message) MarkRead() {
	if m.DeliveredAt == nil {
		m.MarkDelivered()
	}
	now := time.Now().UTC()
	m.ReadAt = &now
	m.Status = StatusRead
}

func (m *Message) MarkFailed() {
	m.Status = StatusFailed
}

func (m *Message) RequiresRecipientsList() bool {
	return m.IsBroadcast()
	// Later: return m.IsBroadcast() || m.IsGroup()
}
