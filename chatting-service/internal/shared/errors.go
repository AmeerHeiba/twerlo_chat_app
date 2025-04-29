package shared

import (
	"errors"
	"net/http"
)

// Domain Errors (Business Rules)
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Application Errors (Use Cases)
var (
	ErrWeakPassword   = errors.New("password must be at least 8 characters")
	ErrMessageTooLong = errors.New("message exceeds 1000 characters")
)

// Infrastructure Errors (DB, External Services)
var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrFileUploadFailed   = errors.New("file upload failed")
)

// HTTP Error Responses
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ToHTTPError(err error) HTTPError {
	switch err {
	case ErrUserNotFound:
		return HTTPError{http.StatusNotFound, "User not found"}
	case ErrUserAlreadyExists:
		return HTTPError{http.StatusConflict, "Username taken"}
	case ErrWeakPassword:
		return HTTPError{http.StatusBadRequest, "Password too weak"}
	default:
		return HTTPError{http.StatusInternalServerError, "Internal server error"}
	}
}
