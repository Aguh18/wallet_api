package response

import "net/http"

// Response represents standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error details
type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Success creates success response
func Success(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// Error creates error response
func Error(code int, message string) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}

// WithStatus returns HTTP status code based on success
func (r Response) WithStatus() int {
	if r.Success {
		return http.StatusOK
	}
	if r.Error != nil {
		return r.Error.Code
	}
	return http.StatusInternalServerError
}
