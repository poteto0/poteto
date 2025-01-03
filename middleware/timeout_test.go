package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/poteto-go/poteto"
)

type ret struct {
	Name string `json:"name"`
}

func TestTimeout(t *testing.T) {
	timeoutConfig := DefaultTimeoutConfig
	timeoutConfig.Limit = time.Millisecond * 300

	tests := []struct {
		name     string
		config   TimeoutConfig
		handler  poteto.HandlerFunc
		expected string
	}{
		{
			"Test done case",
			timeoutConfig,
			func(ctx poteto.Context) error {
				return ctx.JSON(http.StatusOK, ret{Name: "test"})
			},
			`{"name":"test"}`,
		},
		{
			"Test done case if not config",
			TimeoutConfig{},
			func(ctx poteto.Context) error {
				return ctx.JSON(http.StatusOK, ret{Name: "test"})
			},
			`{"name":"test"}`,
		},
		{
			"Test timeout case",
			timeoutConfig,
			func(ctx poteto.Context) error {
				time.Sleep(500 * time.Millisecond)
				return ctx.JSON(http.StatusOK, ret{Name: "test"})
			},
			`{"message":"Gateway Time Out"}`,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			timeout := TimeoutWithConfig(it.config)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "https://example.com/test", nil)
			ctx := poteto.NewContext(w, req)

			timeout_handler := timeout(it.handler)
			timeout_handler(ctx)

			if w.Body.String()[0:10] != it.expected[0:10] {
				t.Errorf(w.Body.String())
				t.Errorf("Unmatched")
			}
		})
	}
}

func TestPanicTimeout(t *testing.T) {
	timeoutConfig := DefaultTimeoutConfig
	timeoutConfig.Limit = time.Millisecond * 300
	timeout := TimeoutWithConfig(timeoutConfig)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://example.com/test", nil)
	ctx := poteto.NewContext(w, req)

	handler := func(poteto.Context) error {
		panic("panic")
	}

	timeout_handler := timeout(handler)
	err := timeout_handler(ctx)

	if err.Error() != "panic" {
		t.Errorf("Not recovered")
	}
}
