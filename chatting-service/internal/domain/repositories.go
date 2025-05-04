package domain

import (
	"context"
	"io"
	"time"
)

// MessageQuery is used for pagination and filtering
type MessageQuery struct {
	Limit       int       // Number of messages to return
	Offset      int       // Pagination offset
	Before      time.Time // Return messages before this time
	After       time.Time // Return messages after this time
	SortBy      string    // "asc" or "desc"
	MessageType string    // Filter by message type
	HasMedia    *bool     // Filter by media presence
	Status      string    // Filter by status
}

type UserRepository interface {
	Create(ctx context.Context, userName, email, passwordHash string) (*User, error)
	FindByID(ctx context.Context, userID uint) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindProfileByID(ctx context.Context, userID uint) (*User, error)
	Update(ctx context.Context, userID uint, username, email string) error
	UpdatePassword(ctx context.Context, userID uint, passwordHash string) error
	UpdateLastActiveAt(ctx context.Context, userID uint) error
	Exists(ctx context.Context, username string) (bool, error)
}

type MessageRepository interface {
	Create(ctx context.Context, senderID uint, content, mediaURL string, messageType MessageType) (*Message, error)
	FindByID(ctx context.Context, messageID uint) (*Message, error)
	FindConversation(ctx context.Context, user1ID, user2ID uint, query MessageQuery) ([]Message, error)
	FindUserMessages(ctx context.Context, userID uint, query MessageQuery) ([]Message, error)
	FindBroadcasts(ctx context.Context, broadcasterID uint, query MessageQuery) ([]Message, error)
	MarkAsDelivered(ctx context.Context, messageID uint) error
	MarkAsRead(ctx context.Context, messageID uint, recipientID uint) error
	Delete(ctx context.Context, messageID uint) error
}

type MessageRecipientRepository interface {
	Create(ctx context.Context, messageID uint, recipientID uint) error
	CreateBulk(ctx context.Context, messageID uint, recipientIDs []uint) error
}

type MessageService interface {
	SendText(ctx context.Context, senderID, recipientID uint, content string) (*Message, error)
	SendMedia(ctx context.Context, senderID, recipientID uint, content string, mediaURL string) (*Message, error)
	Broadcast(ctx context.Context, senderID uint, content string, recipientIDs []uint) (*Message, error)
	GetConversation(ctx context.Context, user1ID, user2ID uint, query MessageQuery) ([]Message, error)
	MarkAsRead(ctx context.Context, messageID uint, recipientID uint) error
}

//Auth Interfaces

type TokenProvider interface {
	GenerateToken(ctx context.Context, user *User) (string, error)
	GenerateRefreshToken(ctx context.Context, user *User) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error)
	ValidateRefreshToken(ctx context.Context, tokenString string) (*TokenClaims, error)
	GetAccessExpiry() time.Duration
	GetRefreshExpiry() time.Duration
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (interface{}, error)
	Refresh(ctx context.Context, refreshToken string) (interface{}, error)
	Logout(ctx context.Context, token string) error
}

//Media Interfaces

type MediaStorage interface {
	Upload(ctx context.Context, file io.Reader, filename string, size int64) (string, error)
	Delete(ctx context.Context, url string) error
	GetSignedURL(ctx context.Context, url string) (string, error)
}

type MediaService interface {
	Upload(ctx context.Context, userID uint, file io.Reader, filename string, size int64) (*MediaResponse, error)
	GetByUser(ctx context.Context, userID uint) ([]MediaResponse, error)
	Delete(ctx context.Context, userID uint, mediaID uint) error
}

//Real Time Interfaces

type MessageNotifier interface {
	Subscribe(ctx context.Context, userID uint) (<-chan *Message, error)
	Unsubscribe(ctx context.Context, userID uint) error
	Notify(ctx context.Context, message *Message) error
	Broadcast(ctx context.Context, message *Message, recipientIDs []uint) error
}

type PresenceTracker interface {
	Track(ctx context.Context, userID uint) error
	Untrack(ctx context.Context, userID uint) error
	GetOnline(ctx context.Context) ([]uint, error)
}

// Transaction Manager for repositories

type Repositories struct {
	Users             UserRepository
	Messages          MessageRepository
	MessageRecipients MessageRecipientRepository
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context, repos *Repositories) error) error
}
