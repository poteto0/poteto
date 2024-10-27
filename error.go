package poteto

import (
	"fmt"
	"net/http"
)

type HttpError interface {
	Error() string
	SetInternalError(err error)
	Unwrap() error
}

type httpError struct {
	InternalError error `json:"internal_error"`
	Message       any   `json:"message"`
	Code          int   `json:"code"`
}

func NewHTTPError(code int, messages ...any) HttpError {
	httpErr := &httpError{Code: code, Message: http.StatusText(code)}
	if len(messages) > 0 {
		httpErr.Message = messages[0]
	}
	return httpErr
}

func (he *httpError) Error() string {
	if he.InternalError == nil {
		return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
	}

	return fmt.Sprintf(
		"code=%d, message=%v, internal=%v",
		he.Code, he.Message, he.InternalError,
	)
}

func (he *httpError) SetInternalError(err error) {
	he.InternalError = err
}

// ↓ For Satisfy Error Interface ↓ //
func (he *httpError) Unwrap() error {
	return he.InternalError
}
