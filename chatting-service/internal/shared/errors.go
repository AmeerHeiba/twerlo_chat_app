package shared

import (
	"errors"
	"net/http"
)

// Error Codes (Machine-readable)
const (
	CodeUserNotFound = "USER_NOT_FOUND"
	CodeUserExists   = "USER_EXISTS"
	CodeInvalidCreds = "INVALID_CREDENTIALS"
	CodeWeakPassword = "WEAK_PASSWORD"
	CodeDBConnection = "DB_CONNECTION_FAILED"
	CodeInternal     = "INTERNAL_ERROR"
)

// Domain Errors (Business Rules)
var (
	ErrUsernameTooShort     = errors.New("user name should be longer than 3 letters")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrWeakPassword         = errors.New("password must be at least 8 characters")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmptyMessage         = errors.New("message must contain text or media")
	ErrMessageTooLong       = errors.New("message exceeds 1000 character limit")
	ErrInvalidMediaURL      = errors.New("invalid media URL format")
	ErrInvalidRecipient     = errors.New("direct messages require exactly one recipient")
	ErrInvalidBroadcast     = errors.New("broadcasts cannot have direct recipient")
	ErrMissingRecOrSenderID = errors.New("both message and user IDs are required")
	ErrNoRecipients         = errors.New("message requires at least one recipient")
	ErrDirectMessageNoList  = errors.New("direct messages should not specify recipients list")
)

// Application Errors (Use Cases)
var (
	ErrTest = errors.New("place - holder for now ")
)

// Infrastructure Errors (DB, External Services)
var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrFileUploadFailed   = errors.New("file upload failed")
)

// HTTP Error Responses
type HTTPError struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func ToHTTPError(err error) HTTPError {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return HTTPError{
			Status:  http.StatusNotFound,
			Code:    CodeUserNotFound,
			Message: "User not found",
		}
	case errors.Is(err, ErrUserAlreadyExists):
		return HTTPError{
			Status:  http.StatusConflict,
			Code:    CodeUserExists,
			Message: "Username already taken",
		}
	case errors.Is(err, ErrWeakPassword):
		return HTTPError{
			Status:  http.StatusBadRequest,
			Code:    CodeWeakPassword,
			Message: "Password must be at least 8 characters",
		}
	default:
		return HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    CodeInternal,
			Message: "Internal server error",
		}
	}
}

// WithDetails adds additional error context
func (e HTTPError) WithDetails(details interface{}) HTTPError {
	e.Details = details
	return e
}
