// pkg/errors/errors.go
package errors

import (
	"fmt"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

var (
	ErrInvalidEmail   = NewAppError(400, "Invalid email address", "")
	ErrUserNotFound   = NewAppError(404, "User not found", "")
	ErrInvalidToken   = NewAppError(401, "Invalid or expired token", "")
	ErrEmailExists    = NewAppError(409, "Email address already exists", "")
	ErrInternalServer = NewAppError(500, "Internal server error", "")
)
