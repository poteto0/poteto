package poteto

import (
	"testing"
)

var rtr *router

func TestMain(m *testing.M) {
	beforeEach()

	m.Run()

	//os.Exit(code)
}

func beforeEach() {
	rtr = NewRouter([]string{"GET", "POST", "PUT", "DELETE"}).(*router)
}

func TestAdd(t *testing.T) {
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
			err := rtr.add(it.method, it.path, nil)
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

func TestGetRoutesByMethod(t *testing.T) {
	rtr.GET("users/get", nil)

	routes := rtr.GetRoutesByMethod("GET")
	child, ok := routes.children["users"].(*route)
	if !ok || child.key != "users" {
		t.Errorf("FATAL add top param")
	}

	cchild, ok := child.children["get"].(*route)
	if !ok || cchild.key != "get" {
		t.Errorf("FATAL add bottom param")
	}
}
