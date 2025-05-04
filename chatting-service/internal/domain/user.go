package domain

import (
	"strings"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"golang.org/x/crypto/bcrypt"

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

//Value Objects and Bussiness rules for user "behaviour of user object"

// BeforeCreate sets the LastActiveAt field to the current time
// before a new User record is created in the database.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Only validate new users
	if u.ID == 0 { // New user
		return u.Validate()
	}
	return nil
}

func (u *User) AfterFind(tx *gorm.DB) error {
	return nil // Skip all validation on find
}

// check user rules for user critical fields validty
func (u *User) Validate() error {
	switch {
	case len(u.Username) < 3:
		return shared.ErrUsernameTooShort
	case !strings.Contains(u.Email, "@") || !strings.Contains(u.Email, "."):
		return shared.ErrInvalidEmail
	case u.PasswordHash == "" || len(u.PasswordHash) < 8:
		return shared.ErrWeakPassword
	}
	return nil
}

// SetPassword securely hashes and stores password
func (u *User) SetPassword(plainText string) error {
	if len(plainText) < 8 {
		return shared.ErrWeakPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies against stored hash
func (u *User) CheckPassword(plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainText))
	return err == nil
}
func (u *User) UpdateLastActive() {
	u.LastActiveAt = time.Now().UTC()
}
