package poteto

import (
	"testing"
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
