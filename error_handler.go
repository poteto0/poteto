package poteto

import (
	"net/http"
)

type HttpErrorHandler interface {
	HandleHttpError(err error, ctx Context)
}

type httpErrorHandler struct{}

func (heh *httpErrorHandler) HandleHttpError(err error, ctx Context) {
	if ctx.GetResponse().isCommitted {
		return
	}

	httpErr, ok := err.(*httpError)
	if !ok { // Not Handled
		httpErr = NewHttpError(http.StatusInternalServerError).(*httpError)
	}
	// Unwrap wrapped error
	if httpErr.InternalError != nil {
		if warpedErr, ok := httpErr.InternalError.(*httpError); ok {
			httpErr = warpedErr
		}
	}

	message := httpErr.Message

	switch m := httpErr.Message.(type) {
	case string:
		message = map[string]string{"message": m}
	}

	// Send response
	err = ctx.JSON(httpErr.Code, message)
}
