package domain

import jwt "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	SessionID string `json:"sid,omitempty"`
	IsRefresh bool   `json:"is_refresh,omitempty"` // Distinguish refresh tokens
}
