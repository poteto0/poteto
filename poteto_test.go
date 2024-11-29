package poteto

import (
	stdContext "context"
	"net/http"
	"testing"
	"time"
)

func TestAddRouteToPoteto(t *testing.T) {
	poteto := New()

	tests := []struct {
		name   string
		method string
		path   string
		want   bool
	}{
		{"success add new route", http.MethodGet, "/users/find", false},
		{"fail add already existed route", http.MethodGet, "/users/find", true},
		{"success add new method already existed route", http.MethodPost, "/users/find", false},
		{"success add new method already existed route", http.MethodPut, "/users/find", false},
		{"success add new method already existed route", http.MethodPatch, "/users/find", false},
		{"success add new method already existed route", http.MethodDelete, "/users/find", false},
		{"success add new method already existed route", http.MethodHead, "/users/find", false},
		{"success add new method already existed route", http.MethodOptions, "/users/find", false},
		{"success add new method already existed route", http.MethodTrace, "/users/find", false},
		{"success add new method already existed route", http.MethodConnect, "/users/find", false},
	}

	for _, it := range tests {
		t.Run(it.name, func(tt *testing.T) {
			var err error

			switch it.method {
			case http.MethodGet:
				err = poteto.GET(it.path, nil)
			case http.MethodPost:
				err = poteto.POST(it.path, nil)
			case http.MethodPut:
				err = poteto.PUT(it.path, nil)
			case http.MethodPatch:
				err = poteto.PATCH(it.path, nil)
			case http.MethodDelete:
				err = poteto.DELETE(it.path, nil)
			case http.MethodHead:
				err = poteto.HEAD(it.path, nil)
			case http.MethodOptions:
				err = poteto.OPTIONS(it.path, nil)
			case http.MethodTrace:
				err = poteto.TRACE(it.path, nil)
			case http.MethodConnect:
				err = poteto.CONNECT(it.path, nil)
			}
			if it.want {
				if err == nil {
					t.Errorf("FATAL: success already existed route")
				}
			} else {
				if err != nil {
					t.Errorf("FATAL: fail new route")
				}
			}
		})
	}
}

func TestRunAndStop(t *testing.T) {
	p := New()

	tests := []struct {
		name  string
		port1 string
		port2 string
	}{
		{"Test :8080", ":8080", ""},
		{"Test 8080", "8080", ""},
		{"Test collision panic", ":8080", ":8080"},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			errChan := make(chan error)
			go func() {
				errChan <- p.Run(it.port1)
			}()

			errChan2 := make(chan error)
			if it.port2 != "" {
				go func() {
					errChan2 <- p.Run(it.port2)
				}()
			}

			select {
			case <-time.After(500 * time.Millisecond):
				if err := p.Stop(stdContext.Background()); err != nil {
					t.Errorf("Unmatched")
				}
			case <-errChan:
				return
			case <-errChan2:
				return
			}
		})
	}
}

func TestSetLogger(t *testing.T) {
	p := New().(*poteto)
	logger := func(msg string) {
		return
	}

	if p.logger != nil {
		t.Error("Unmatched")
	}

	p.SetLogger(logger)
	if p.logger == nil {
		t.Errorf("Unmatched")
	}
}
