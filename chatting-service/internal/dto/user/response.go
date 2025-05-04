package user

import "time"

type ProfileResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	LastActive time.Time `json:"last_active"`
	Status     string    `json:"status"`
}
