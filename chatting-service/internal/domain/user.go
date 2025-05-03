package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string     `gorm:"uniqueIndex;size:50;not null"`
	Email        string     `gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string     `gorm:"type:text;not null"`
	LastActiveAt time.Time  `gorm:"index"`
	Status       UserStatus `gorm:"type:user_status;default:'offline'"`

	// Relationships
	SentMessages     []Message `gorm:"foreignKey:SenderID"`
	ReceivedMessages []Message `gorm:"foreignKey:RecipientID"`
	Broadcasts       []Message `gorm:"foreignKey:BroadcasterID"`
}

//Value Objects

// BeforeCreate sets the LastActiveAt field to the current time
// before a new User record is created in the database.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.LastActiveAt = time.Now()
	return nil
}
