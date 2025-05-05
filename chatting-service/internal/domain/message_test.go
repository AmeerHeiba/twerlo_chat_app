package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMessage_Validate(t *testing.T) {

	tests := []struct {
		name    string
		message Message
		wantErr error
	}{
		// Content Validation Tests
		{
			name:    "EmptyMessage",
			message: Message{Content: "", MediaURL: ""},
			wantErr: ErrEmptyMessage,
		},
		{
			name: "ValidTextOnly",
			message: Message{
				Content:     "Hello",
				MessageType: MessageDirect,
				RecipientID: uintPtr(1),
				SenderID:    1,
			},
			wantErr: nil,
		},
		{
			name: "ValidMediaOnly",
			message: Message{
				MediaURL:    "https://example.com/image.jpg", //TODO will decide later based on storage implmentation
				MessageType: MessageDirect,
				RecipientID: uintPtr(1),
				SenderID:    1,
			},
			wantErr: nil,
		},
		{
			name:    "TextTooLong",
			message: Message{Content: strings.Repeat("a", 1001), MediaURL: ""},
			wantErr: ErrMessageTooLong,
		},
		{
			name:    "InvalidMediaURL",
			message: Message{Content: "", MediaURL: "invalid-url"},
			wantErr: ErrInvalidMediaURL,
		},

		// Direct Message Recipient Tests
		{
			name: "ValidDirectMessage",
			message: Message{
				Content:     "Hi",
				MessageType: MessageDirect,
				RecipientID: uintPtr(2),
				SenderID:    1,
			},
			wantErr: nil,
		},
		{
			name: "DirectMessageMissingRecipient",
			message: Message{
				Content:     "Hi",
				MessageType: MessageDirect,
				SenderID:    1,
			},
			wantErr: ErrInvalidRecipient,
		},
		{
			name: "DirectMessageWithRecipientsList",
			message: Message{
				Content:     "Hi",
				MessageType: MessageDirect,
				RecipientID: uintPtr(2),
				Recipients:  []User{{Model: gorm.Model{ID: 3}}},
				SenderID:    1,
			},
			wantErr: ErrDirectMessageNoList,
		},

		// Broadcast Message Tests
		{
			name: "ValidBroadcast",
			message: Message{
				Content:     "Hello all",
				MessageType: MessageBroadcast,
				Recipients:  []User{{Model: gorm.Model{ID: 2}}, {Model: gorm.Model{ID: 3}}},
				SenderID:    1,
			},
			wantErr: nil,
		},
		{
			name: "BroadcastWithRecipientID",
			message: Message{
				Content:     "Hello",
				MessageType: MessageBroadcast,
				RecipientID: uintPtr(2),
				SenderID:    1,
			},
			wantErr: ErrInvalidBroadcast,
		},
		{
			name: "BroadcastNoRecipients",
			message: Message{
				Content:     "Hello",
				MessageType: MessageBroadcast,
				SenderID:    1,
			},
			wantErr: ErrNoRecipients,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestMessageRecipient_Validate(t *testing.T) {
	tests := []struct {
		name string
		mr   MessageRecipient
		want error
	}{
		{
			name: "Valid",
			mr:   MessageRecipient{MessageID: 1, UserID: 1},
			want: nil,
		},
		{
			name: "MissingMessageID",
			mr:   MessageRecipient{UserID: 1},
			want: ErrMissingRecOrSenderID,
		},
		{
			name: "MissingUserID",
			mr:   MessageRecipient{MessageID: 1},
			want: ErrMissingRecOrSenderID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mr.Validate()
			assert.ErrorIs(t, err, tt.want)
		})
	}
}

func TestMessage_StateTransitions(t *testing.T) {
	t.Run("MarkDelivered", func(t *testing.T) {
		msg := Message{Status: StatusSent, SentAt: time.Now()}
		msg.MarkDelivered()
		assert.Equal(t, StatusDelivered, msg.Status)
		assert.NotNil(t, msg.DeliveredAt)
	})

	t.Run("MarkRead", func(t *testing.T) {
		msg := Message{Status: StatusSent, SentAt: time.Now()}
		msg.MarkRead()
		assert.Equal(t, StatusRead, msg.Status)
		assert.NotNil(t, msg.DeliveredAt)
		assert.NotNil(t, msg.ReadAt)
	})
}

func TestMessage_HelperMethods(t *testing.T) {
	t.Run("IsDirect", func(t *testing.T) {
		msg := Message{MessageType: MessageDirect}
		assert.True(t, msg.IsDirect())
		assert.False(t, msg.IsBroadcast())
	})

	t.Run("IsBroadcast", func(t *testing.T) {
		msg := Message{MessageType: MessageBroadcast}
		assert.True(t, msg.IsBroadcast())
		assert.False(t, msg.IsDirect())
	})

	t.Run("RequiresRecipientsList", func(t *testing.T) {
		t.Run("Direct", func(t *testing.T) {
			msg := Message{MessageType: MessageDirect}
			assert.False(t, msg.RequiresRecipientsList())
		})
		t.Run("Broadcast", func(t *testing.T) {
			msg := Message{MessageType: MessageBroadcast}
			assert.True(t, msg.RequiresRecipientsList())
		})
	})
}

func TestMessage_BeforeCreate(t *testing.T) {
	baseMsg := Message{
		Content:     "test",
		MessageType: MessageDirect,
		RecipientID: uintPtr(1),
	}

	t.Run("SetsDefaultSentAt", func(t *testing.T) {
		msg := baseMsg
		msg.SentAt = time.Time{} // Zero value
		err := msg.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.False(t, msg.SentAt.IsZero())
	})

	t.Run("SetsDefaultStatus", func(t *testing.T) {
		msg := baseMsg
		msg.Status = ""
		err := msg.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, StatusSent, msg.Status)
	})
}

// Helper function to create uint pointers
func uintPtr(i uint) *uint {
	return &i
}
