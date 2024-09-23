package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poteto0/poteto"
)

type TestVal struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

func TestCORSWithConfigByDefault(t *testing.T) {
	config := CORSConfig{
		AllowOrigins: []string{},
		AllowMethods: []string{},
	}

	t.Run("allow all origins", func(t *testing.T) {
		cors := CORSWithConfig(config)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://example.com/test", nil)
		context := poteto.NewContext(w, req)

		handler := func(ctx poteto.Context) error {
			return ctx.JSON(http.StatusOK, TestVal{Name: "test", Val: "val"})
		}

		cors_handler := cors(handler)
		cors_handler(context)
		result := w.Body.String()
		expected := `{"name":"test","val":"val"}`
		if result[0:27] != expected[0:27] {
			t.Errorf("Wrong result")
			t.Errorf(fmt.Sprintf("expected: %s", expected))
			t.Errorf(fmt.Sprintf("actual: %s", result))
		}
	})
}

func TestGetAllowSubDomain(t *testing.T) {
	tests := []struct {
		name         string
		origin       string
		allowOrigins []string
		expected     string
	}{
		{"test wildcard return true", "https://example.com", []string{"*"}, "https://example.com"},
		{"test match same domain", "https://example.com", []string{"https://example.com"}, "https://example.com"},
		{"test matched subdomain", "https://exmaple.com.test", []string{"https://example.com.*"}, "https://exmaple.com.test"},
		{"test not matched", "https://hello.world.com", []string{"https://exmaple.com"}, ""},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			result := getAllowSubDomain(it.origin, it.allowOrigins)
			if result != it.expected {
				t.Errorf("Not matched")
				t.Errorf(fmt.Sprintf("expected: %s", it.expected))
				t.Errorf(fmt.Sprintf("actual: %s", result))
			}
		})
	}
}

func TestGetAllowOrigin(t *testing.T) {
	tests := []struct {
		name                string
		origin              string
		allowOriginPatterns []string
		expected            string
	}{
		{"test match case", "https://example.com", []string{wrapRegExp("https://example.*")}, "https://example.com"},
		{"test not match case", "https://example.com", []string{wrapRegExp("https://hello.world.com")}, ""},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			result := getAllowOrigin(it.origin, it.allowOriginPatterns)
			if result != it.expected {
				t.Errorf("Not matched")
				t.Errorf(fmt.Sprintf("expected: %s", it.expected))
				t.Errorf(fmt.Sprintf("actual: %s", result))
			}
		})
	}
}
