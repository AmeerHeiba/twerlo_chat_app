package user

import "time"

type ProfileResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	LastActive time.Time `json:"last_active"`
	Status     string    `json:"status"`
}

type MessageHistoryResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int64             `json:"total"`
}

type MessageResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
}
