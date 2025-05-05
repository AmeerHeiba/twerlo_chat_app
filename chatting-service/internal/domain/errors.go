package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrWeakPassword         = errors.New("weak password")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUsernameTooShort     = errors.New("user name should be longer than 3 letters")
	ErrEmptyMessage         = errors.New("message must contain text or media")
	ErrMessageTooLong       = errors.New("message exceeds 1000 character limit")
	ErrInvalidMediaURL      = errors.New("invalid media URL format")
	ErrInvalidRecipient     = errors.New("direct messages require exactly one recipient")
	ErrInvalidBroadcast     = errors.New("broadcasts cannot have direct recipient")
	ErrNoRecipients         = errors.New("message requires at least one recipient")
	ErrMessageNotFound      = errors.New("message not found")
	ErrInvalidMessageType   = errors.New("invalid message type")
	ErrInvalidMessageStatus = errors.New("invalid message status")
	ErrUserExists           = errors.New("User already exists")
	ErrMissingRecOrSenderID = errors.New("both message and user IDs are required")
	ErrDirectMessageNoList  = errors.New("direct messages should not specify recipients list")
	ErrInvalidRecipientID   = errors.New("invalid recipient ID")
	ErrInvalidSenderID      = errors.New("invalid sender ID")
	ErrEmailExists          = errors.New("email already exists")
	ErrUserAlreadyExists    = errors.New("user already exists")
)
