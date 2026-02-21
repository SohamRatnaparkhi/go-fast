package handler

import "fmt"

// HTTPError is an error that carries an HTTP status code.
//
// When a handler returns an HTTPError, the adapter uses its Code for the
// response status instead of the default 500.
type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string { return e.Message }

// NewHTTPError creates an HTTPError with the given status code and message.
func NewHTTPError(code int, msg string) *HTTPError {
	return &HTTPError{Code: code, Message: msg}
}

// BadRequest returns a 400 error.
func BadRequest(msg string) *HTTPError { return &HTTPError{Code: 400, Message: msg} }

// Unauthorized returns a 401 error.
func Unauthorized(msg string) *HTTPError { return &HTTPError{Code: 401, Message: msg} }

// Forbidden returns a 403 error.
func Forbidden(msg string) *HTTPError { return &HTTPError{Code: 403, Message: msg} }

// NotFound returns a 404 error.
func NotFound(msg string) *HTTPError { return &HTTPError{Code: 404, Message: msg} }

// Conflict returns a 409 error.
func Conflict(msg string) *HTTPError { return &HTTPError{Code: 409, Message: msg} }

// UnprocessableEntity returns a 422 error.
func UnprocessableEntity(msg string) *HTTPError { return &HTTPError{Code: 422, Message: msg} }

// InternalServerError returns a 500 error.
func InternalServerError(msg string) *HTTPError { return &HTTPError{Code: 500, Message: msg} }

// Errorf creates an HTTPError with a formatted message.
func Errorf(code int, format string, args ...interface{}) *HTTPError {
	return &HTTPError{Code: code, Message: fmt.Sprintf(format, args...)}
}
