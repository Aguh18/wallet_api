package errors

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrNotFound       = New(404, "Resource not found", nil)
	ErrUnauthorized   = New(401, "Unauthorized", nil)
	ErrForbidden      = New(403, "Forbidden", nil)
	ErrBadRequest     = New(400, "Bad request", nil)
	ErrInternalServer = New(500, "Internal server error", nil)
	ErrConflict       = New(409, "Resource already exists", nil)
)
