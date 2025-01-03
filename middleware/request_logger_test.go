package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poteto-go/poteto"
)

func TestRequestLogger(t *testing.T) {
	p := poteto.New()

	tests := []struct {
		name     string
		config   RequestLoggerConfig
		expected bool
	}{
		{"Test default config", DefaultRequestLoggerConfig, true},
		{"Test empty config", RequestLoggerConfig{}, false},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			logConfig := it.config
			var result RequestLoggerValues
			logConfig.LogHandleFunc = func(ctx poteto.Context, values RequestLoggerValues) error {
				result = values
				return nil
			}
			p.Register(RequestLoggerWithConfig(logConfig))

			p.GET("/test", func(ctx poteto.Context) error {
				return ctx.JSON(http.StatusTeapot, nil)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			p.ServeHTTP(rec, req)

			if http.StatusTeapot != rec.Code {
				t.Errorf("Not go through logger")
			}

			if (result.Status != http.StatusTeapot) == it.expected {
				t.Errorf("Not matched")
			}

			if (result.Method != "GET") == it.expected {
				t.Errorf("Not matched")
			}
		})
	}
}
