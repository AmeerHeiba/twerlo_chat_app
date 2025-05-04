package config

import (
	"os"
	"time"
)

type AuthConfig struct {
	JWTSecret          string
	AccessTokenExpiry  time.Duration // e.g., 15 minutes
	RefreshTokenExpiry time.Duration // e.g., 7 days
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AccessTokenExpiry:  time.Minute * 15,
		RefreshTokenExpiry: time.Hour * 24 * 7,
	}
}
