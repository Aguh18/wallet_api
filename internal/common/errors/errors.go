package errors

import "fmt"

// AppError represents application error
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// New creates new AppError
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common errors
var (
	ErrNotFound       = New(404, "Resource not found", nil)
	ErrUnauthorized   = New(401, "Unauthorized", nil)
	ErrForbidden      = New(403, "Forbidden", nil)
	ErrBadRequest     = New(400, "Bad request", nil)
	ErrInternalServer = New(500, "Internal server error", nil)
	ErrConflict       = New(409, "Resource already exists", nil)
)
