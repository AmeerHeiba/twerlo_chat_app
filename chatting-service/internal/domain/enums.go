package domain

type MessageStatus string

const (
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
)

type UserStatus string

const (
	UserOnline  UserStatus = "online"
	UserOffline UserStatus = "offline"
	UserAway    UserStatus = "away"
)

type MessageType string

const (
	MessageDirect    MessageType = "direct"
	MessageBroadcast MessageType = "broadcast"
)
