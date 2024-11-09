package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poteto0/poteto"
)

func TestRequestLogger(t *testing.T) {
	p := poteto.New()

	logConfig := DefaultRequestLoggerConfig
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

	if result.Status != http.StatusTeapot {
		t.Errorf("Not matched")
	}

	if result.Method != "GET" {
		t.Errorf("Not matched")
	}
}
