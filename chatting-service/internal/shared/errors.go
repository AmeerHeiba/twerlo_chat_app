package shared

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

const (
	EmailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type Error struct {
	Code    string      `json:"code"`    // Machine-readable code
	Message string      `json:"message"` // Human-readable message
	Status  int         `json:"-"`       // HTTP status code
	Details interface{} `json:"details,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) WithDetails(details interface{}) Error {
	return Error{
		Code:    e.Code,
		Message: e.Message,
		Status:  e.Status,
		Details: details,
	}
}

// Predefined common errors
var (
	// 4xx Errors
	ErrBadRequest = Error{
		Code:    "BAD_REQUEST",
		Message: "Invalid request",
		Status:  http.StatusBadRequest,
	}

	ErrUnauthorized = Error{
		Code:    "UNAUTHORIZED",
		Message: "Not authorized",
		Status:  http.StatusUnauthorized,
	}

	ErrForbidden = Error{
		Code:    "FORBIDDEN",
		Message: "Access denied",
		Status:  http.StatusForbidden,
	}

	ErrNotFound = Error{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
		Status:  http.StatusNotFound,
	}
	ErrConflict = Error{
		Code:    "CONFLICT",
		Message: "Resource already exists",
		Status:  http.StatusConflict,
	}

	// Validation Errors
	ErrValidation = Error{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Status:  http.StatusBadRequest,
	}

	ErrUsernameTooShort = Error{
		Code:    "USERNAME_TOO_SHORT",
		Message: "Username must be at least 3 characters",
		Status:  http.StatusBadRequest,
	}
	ErrInvalidEmailFormat = Error{
		Code:    "INVALID_EMAIL_FORMAT",
		Message: "Email must be a valid format",
		Status:  http.StatusBadRequest,
	}
	ErrPasswordTooWeak = Error{
		Code:    "PASSWORD_TOO_WEAK",
		Message: "Password must be at least 8 characters",
		Status:  http.StatusBadRequest,
	}

	// 5xx Errors
	ErrInternalServer = Error{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}

	// Domain-specific errors
	ErrUserNotFound = Error{
		Code:    "USER_NOT_FOUND",
		Message: "User not found",
		Status:  http.StatusNotFound,
	}

	ErrUserExists = Error{
		Code:    "USER_EXISTS",
		Message: "User already exists",
		Status:  http.StatusConflict,
	}

	ErrInvalidCredentials = Error{
		Code:    "INVALID_CREDENTIALS",
		Message: "Invalid username or password",
		Status:  http.StatusUnauthorized,
	}

	ErrWeakPassword = Error{
		Code:    "WEAK_PASSWORD",
		Message: "Password must be at least 8 characters",
		Status:  http.StatusBadRequest,
	}

	// Database errors
	ErrDatabaseOperation = Error{
		Code:    "DATABASE_ERROR",
		Message: "Database operation failed",
		Status:  http.StatusInternalServerError,
	}

	ErrRecordNotFound = Error{
		Code:    "RECORD_NOT_FOUND",
		Message: "Record not found",
		Status:  http.StatusNotFound,
	}

	ErrDuplicateEntry = Error{
		Code:    "DUPLICATE_ENTRY",
		Message: "Duplicate entry",
		Status:  http.StatusConflict,
	}

	//Generic
	ErrRateLimited = Error{
		Code:    "RATE_LIMITED",
		Message: "Too many requests",
		Status:  http.StatusTooManyRequests,
	}
	ErrServiceUnavailable = Error{
		Code:    "SERVICE_UNAVAILABLE",
		Message: "Service temporarily unavailable",
		Status:  http.StatusServiceUnavailable,
	}
)

// Helper to convert third-party errors to the local error type
func NormalizeError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}

	return err
}
