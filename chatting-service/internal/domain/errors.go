package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWeakPassword       = errors.New("weak password")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidToken       = errors.New("invalid token")
)
