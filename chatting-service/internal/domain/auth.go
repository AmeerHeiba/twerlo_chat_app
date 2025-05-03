package domain

import jwt "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	jwt.RegisteredClaims
	Id       uint
	Username string
	Email    string
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
