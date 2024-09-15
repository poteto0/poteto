package poteto

import (
	"testing"
)

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

func TestInsertAndSearch(t *testing.T) {
	url := "https://example.com/v1/users/find/poteto"

	route := NewRoute().(*route)

	route.Insert("GET", url, nil)

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

func BenchmarkInsertAndSearch(b *testing.B) {
	urls := []string{
		"https://example.com/v1/users/find/poteto",
		"https://example.com/v1/users/find/potato",
		"https://example.com/v1/users/find/jagaimo",
		"https://example.com/v1/users/create/poteto",
		"https://example.com/v1/users/create/potato",
		"https://example.com/v1/users/create/jagaimo",
		"https://example.com/v1/members/find/poteto",
		"https://example.com/v1/members/find/potato",
		"https://example.com/v1/members/find/jagaimo",
		"https://example.com/v1/members/create/poteto",
		"https://example.com/v1/members/create/potato",
		"https://example.com/v1/members/create/jagaimo",
		"https://example.com/v2/users/find/poteto",
		"https://example.com/v2/users/find/potato",
		"https://example.com/v2/users/find/jagaimo",
		"https://example.com/v2/users/create/poteto",
		"https://example.com/v2/users/create/potato",
		"https://example.com/v2/users/create/jagaimo",
		"https://example.com/v2/members/find/poteto",
		"https://example.com/v2/members/find/potato",
		"https://example.com/v2/members/find/jagaimo",
		"https://example.com/v2/members/create/poteto",
		"https://example.com/v2/members/create/potato",
		"https://example.com/v2/members/create/jagaimo",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Insert
		route := NewRoute().(*route)
		for _, url := range urls {
			route.Insert("GET", url, nil)
		}

		// Search
		for _, url := range urls {
			route.Search(url)
		}
	}
}
