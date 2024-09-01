package poteto

import (
	"testing"
)

func ArrangeRoute() Route {
	url := "https://example.com/v1/users/find/poteto"

	route := NewRoute().(*route)

	route.Insert("GET", url, nil)

	return route
}

func TestNewRoute(t *testing.T) {
	// Arrange
	want := &route{
		key:      "",
		method:   "",
		children: make(map[string]Route),
	}

	got := NewRoute().(*route)

	if got.key != want.key {
		t.Errorf("Cannot initialize Route: key")
	}

	if got.method != want.method {
		t.Errorf("Cannot initialize Route: method")
	}

	if len(got.children) != 0 {
		t.Errorf("Cannot initialize Route: method")
	}
}

func TestSearch(t *testing.T) {
	route := ArrangeRoute().(*route)

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"FIND", "https://example.com", "example.com"},
		{"NOT FOUND", "https://fuck.com", ""},
	}

	for _, it := range tests {
		t.Run(it.name, func(tt *testing.T) {
			got := route.Search(it.arg)

			key := ""
			if got != nil {
				key = got.key
			}

			if key != it.want {
				tt.Errorf("Cannot search route")
			}
		})
	}

}
