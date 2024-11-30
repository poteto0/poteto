package poteto

import (
	"reflect"
	"runtime"
	"testing"
)

func TestNewRoute(t *testing.T) {
	// Arrange
	want := &route{
		key:      "",
		children: make(map[string]Route),
	}

	got := NewRoute().(*route)

	if got.key != want.key {
		t.Errorf("Cannot initialize Route: key")
	}

	if len(got.children) != 0 {
		t.Errorf("Cannot initialize Route: method")
	}
}

func TestInsertAndSearch(t *testing.T) {
	url := "/example.com/v1/users/find/poteto"

	route := NewRoute().(*route)

	route.InsertNew("/", nil)
	route.InsertNew(url, nil)
	route.InsertNew("/users/:id", nil)

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"FIND empty", "/", ""},
		{"FIND", "/example.com", "example.com"},
		{"NOT FOUND", "/test.com", ""},
		{"PARAM ROUTING", "/users/1", ":id"},
	}

	for _, it := range tests {
		t.Run(it.name, func(tt *testing.T) {
			got, _ := route.Search(it.arg)

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
		"/example.com/v1/users/find/poteto",
		"/example.com/v1/users/find/potato",
		"/example.com/v1/users/find/jagaimo",
		"/example.com/v1/users/create/poteto",
		"/example.com/v1/users/create/potato",
		"/example.com/v1/users/create/jagaimo",
		"/example.com/v1/members/find/poteto",
		"/example.com/v1/members/find/potato",
		"/example.com/v1/members/find/jagaimo",
		"/example.com/v1/members/create/poteto",
		"/example.com/v1/members/create/potato",
		"/example.com/v1/members/create/jagaimo",
		"/example.com/v2/users/find/poteto",
		"/example.com/v2/users/find/potato",
		"/example.com/v2/users/find/jagaimo",
		"/example.com/v2/users/create/poteto",
		"/example.com/v2/users/create/potato",
		"/example.com/v2/users/create/jagaimo",
		"/example.com/v2/members/find/poteto",
		"/example.com/v2/members/find/potato",
		"/example.com/v2/members/find/jagaimo",
		"/example.com/v2/members/create/poteto",
		"/example.com/v2/members/create/potato",
		"/example.com/v2/members/create/jagaimo",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Insert
		route := NewRoute().(*route)
		for _, url := range urls {
			route.Insert(url, nil)
		}

		// Search
		for _, url := range urls {
			route.Search(url)
		}
	}
}

func BenchmarkInsertAndSearchNewRoute(b *testing.B) {
	urls := []string{
		"/example.com/v1/users/find/poteto",
		"/example.com/v1/users/find/potato",
		"/example.com/v1/users/find/jagaimo",
		"/example.com/v1/users/create/poteto",
		"/example.com/v1/users/create/potato",
		"/example.com/v1/users/create/jagaimo",
		"/example.com/v1/members/find/poteto",
		"/example.com/v1/members/find/potato",
		"/example.com/v1/members/find/jagaimo",
		"/example.com/v1/members/create/poteto",
		"/example.com/v1/members/create/potato",
		"/example.com/v1/members/create/jagaimo",
		"/example.com/v2/users/find/poteto",
		"/example.com/v2/users/find/potato",
		"/example.com/v2/users/find/jagaimo",
		"/example.com/v2/users/create/poteto",
		"/example.com/v2/users/create/potato",
		"/example.com/v2/users/create/jagaimo",
		"/example.com/v2/members/find/poteto",
		"/example.com/v2/members/find/potato",
		"/example.com/v2/members/find/jagaimo",
		"/example.com/v2/members/create/poteto",
		"/example.com/v2/members/create/potato",
		"/example.com/v2/members/create/jagaimo",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Insert
		route := NewRoute().(*route)
		for _, url := range urls {
			route.InsertNew(url, nil)
		}

		// Search
		for _, url := range urls {
			route.Search(url)
		}
	}
}

func TestCollisionDetection(t *testing.T) {
	route := NewRoute().(*route)

	route.Insert("/users/test", getAllUserForTest)
	route.Insert("/users/test", getAllUserForTestById)

	r, _ := route.Search("/users/test")
	if getFunctionName(r.handler) != getFunctionName(getAllUserForTest) {
		t.Errorf("Unmatched")
	}
}

func getFunctionName(handler any) string {
	return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
}
