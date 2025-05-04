package shared

import (
	"errors"
	"net/http"
)

// Error Codes (Machine-readable)
const (
	CodeUserNotFound         = "USER_NOT_FOUND"
	CodeUserExists           = "USER_EXISTS"
	CodeInvalidCreds         = "INVALID_CREDENTIALS"
	CodeWeakPassword         = "WEAK_PASSWORD"
	CodeDBConnection         = "DB_CONNECTION_FAILED"
	CodeInternal             = "INTERNAL_ERROR"
	CodeInvalidUsername      = "INVALID_USERNAME"
	CodeMessageNotFound      = "MESSAGE_NOT_FOUND"
	CodeInvalidMessageType   = "INVALID_MESSAGE_TYPE"
	CodeInvalidMessageStatus = "INVALID_MESSAGE_STATUS"
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
	ErrTest                 = errors.New("place - holder for now ")
	ErrUserExists           = errors.New("User already exists")
	ErrDatabaseConnection   = errors.New("database connection failed")
	ErrFileUploadFailed     = errors.New("file upload failed")
	ErrEmailExists          = errors.New("email already exists")
	ErrMessageNotFound      = errors.New("message not found")
	ErrInvalidMessageType   = errors.New("invalid message type")
	ErrInvalidMessageStatus = errors.New("invalid message status")
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
	case errors.Is(err, ErrUserExists):
		return HTTPError{
			Status:  http.StatusConflict,
			Code:    CodeUserExists,
			Message: "User already exists",
		}

	case errors.Is(err, ErrUsernameTooShort):
		return HTTPError{
			Status:  http.StatusBadRequest,
			Code:    CodeInvalidUsername,
			Message: "Username must be at least 3 characters",
		}
	case errors.Is(err, ErrMessageNotFound):
		return HTTPError{
			Status:  http.StatusNotFound,
			Code:    CodeMessageNotFound,
			Message: "Message not found",
		}
	case errors.Is(err, ErrInvalidMessageType):
		return HTTPError{
			Status:  http.StatusBadRequest,
			Code:    CodeInvalidMessageType,
			Message: "Invalid message type",
		}
	case errors.Is(err, ErrInvalidMessageStatus):
		return HTTPError{
			Status:  http.StatusBadRequest,
			Code:    CodeInvalidMessageStatus,
			Message: "Invalid message status",
		}
	case errors.Is(err, ErrInvalidRecipient):
		return HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_RECIPIENT",
			Message: "Direct messages require exactly one recipient",
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
