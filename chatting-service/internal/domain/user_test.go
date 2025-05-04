package domain

import (
	"testing"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	baseUser := User{
		Username:     "validuser",
		Email:        "valid@example.com",
		PasswordHash: "$2a$10$fakehashedpassword",
	}

	t.Run("Validation", func(t *testing.T) {
		tests := []struct {
			name        string
			modifyFn    func(*User)
			expectedErr error
		}{
			{
				name:        "ValidUser",
				modifyFn:    func(u *User) {},
				expectedErr: nil,
			},
			{
				name: "UsernameTooShort",
				modifyFn: func(u *User) {
					u.Username = "ab"
				},
				expectedErr: shared.ErrUsernameTooShort,
			},
			{
				name: "InvalidEmail",
				modifyFn: func(u *User) {
					u.Email = "invalid-email"
				},
				expectedErr: shared.ErrInvalidEmail,
			},
			{
				name: "EmptyPasswordHash",
				modifyFn: func(u *User) {
					u.PasswordHash = ""
				},
				expectedErr: shared.ErrWeakPassword,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				user := baseUser
				tt.modifyFn(&user)
				err := user.Validate()
				assert.ErrorIs(t, err, tt.expectedErr)
			})
		}
	})

	t.Run("PasswordHashing", func(t *testing.T) {
		t.Run("SetValidPassword", func(t *testing.T) {
			u := baseUser
			err := u.SetPassword("validpass123")
			assert.NoError(t, err)
			assert.NotEqual(t, "validpass123", u.PasswordHash)
			assert.True(t, len(u.PasswordHash) > 0)
		})

		t.Run("RejectWeakPassword", func(t *testing.T) {
			u := baseUser
			err := u.SetPassword("short")
			assert.ErrorIs(t, err, shared.ErrWeakPassword)
			assert.Equal(t, baseUser.PasswordHash, u.PasswordHash)
		})

		t.Run("PasswordVerification", func(t *testing.T) {
			u := baseUser
			plainPass := "testpassword"
			assert.NoError(t, u.SetPassword(plainPass))

			assert.True(t, u.CheckPassword(plainPass))
			assert.False(t, u.CheckPassword("wrongpassword"))
			assert.False(t, u.CheckPassword(""))
		})
	})

	t.Run("ActivityTracking", func(t *testing.T) {
		u := baseUser
		initialTime := u.LastActiveAt

		time.Sleep(10 * time.Millisecond)
		u.UpdateLastActive()

		assert.True(t, u.LastActiveAt.After(initialTime))
		assert.False(t, u.LastActiveAt.IsZero())
	})

	t.Run("GORMHooks", func(t *testing.T) {
		t.Run("BeforeCreateSetsActivity", func(t *testing.T) {
			u := User{
				Username:     "newuser",
				Email:        "new@example.com",
				PasswordHash: "$2a$10$fakehash",
			}
			err := u.BeforeCreate(nil)
			assert.NoError(t, err)
			assert.False(t, u.LastActiveAt.IsZero())
		})

		t.Run("BeforeCreateValidates", func(t *testing.T) {
			u := User{
				Username: "ab", // Too short
				Email:    "new@example.com",
			}
			err := u.BeforeCreate(nil)
			assert.ErrorIs(t, err, shared.ErrUsernameTooShort)
		})

		// t.Run("BeforeUpdateValidates", func(t *testing.T) {
		// 	u := baseUser
		// 	u.Email = "invalid-email"
		// 	err := u.BeforeUpdate(nil)
		// 	assert.ErrorIs(t, err, shared.ErrInvalidEmail)
		// })
	})
}
