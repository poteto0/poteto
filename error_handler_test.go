package poteto

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHttpError(t *testing.T) {
	handler := httpErrorHandler{}

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			"Test Not Handled Error -> Server Error",
			errors.New("not httpError"),
			`{"message":"Internal Server Error"}`,
		},
		{
			"Test Handled Error",
			NewHttpError(http.StatusBadRequest),
			`{"message":"Bad Request"}`,
		},
		{
			"Test wrapped Error",
			&httpError{
				Code:          http.StatusBadRequest,
				Message:       "",
				InternalError: NewHttpError(http.StatusBadRequest),
			},
			`{"message":"Bad Request"}`,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			url := "/example.com"
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", url, nil)
			ctx := NewContext(w, req).(*context)

			handler.HandleHttpError(it.err, ctx)

			if res := w.Body.String()[0:20]; res != it.expected[0:20] {
				t.Errorf(res)
				t.Errorf(it.expected)
				t.Errorf("Unmatched")
			}
		})
	}
}

func TestIsCommitedHandlerHttpErrorCase(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users", nil)
	ctx := NewContext(w, req).(*context)
	ctx.GetResponse().IsCommitted = true

	handler := httpErrorHandler{}
	handler.HandleHttpError(errors.New("error"), ctx)
}
