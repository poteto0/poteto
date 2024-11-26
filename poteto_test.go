package poteto

import (
	stdContext "context"
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
		{"success add new route", "GET", "/users/find", false},
		{"fail add already existed route", "GET", "/users/find", true},
		{"success add new method already existed route", "POST", "/users/find", false},
		{"success add new method already existed route", "PUT", "/users/find", false},
		{"success add new method already existed route", "DELETE", "/users/find", false},
	}

	for _, it := range tests {
		t.Run(it.name, func(tt *testing.T) {
			var err error

			switch it.method {
			case "GET":
				err = poteto.GET(it.path, nil)
			case "POST":
				err = poteto.POST(it.path, nil)
			case "PUT":
				err = poteto.PUT(it.path, nil)
			case "DELETE":
				err = poteto.DELETE(it.path, nil)
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

			select {
			case <-time.After(500 * time.Millisecond):
				if err := p.Stop(stdContext.Background()); err != nil {
					t.Errorf("Unmatched")
				}
			case <-errChan:
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
