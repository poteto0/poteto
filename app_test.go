package poteto

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: echo見る
func TestServeHTTP(t *testing.T) {
	p := New()
	p.GET("/users", getAllUserForTest)
	p.GET("/users/:id", getAllUserForTestById)

	tests := []struct {
		name         string
		reqMethod    string
		reqUrl       string
		expectedCode int
	}{
		{
			"Test Not registered URL",
			"GET",
			"/test",
			http.StatusNotFound,
		},
		{
			"Test static url",
			"GET",
			"/users",
			http.StatusOK,
		},
		{
			"Test param url",
			"GET",
			"/users/1",
			http.StatusOK,
		},
		{
			"Test not registered method",
			"POST",
			"/users",
			http.StatusNotFound,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(it.reqMethod, it.reqUrl, nil)

			p.ServeHTTP(w, req)
			if w.Code != it.expectedCode {
				t.Errorf("Unmatched")
			}
		})
	}
}

type User struct {
	Name string `json:"string"`
}

func getAllUserForTest(ctx Context) error {
	user := User{
		Name: "user",
	}
	return ctx.JSON(http.StatusOK, user)
}

func getAllUserForTestById(ctx Context) error {
	user := User{
		Name: "user1",
	}
	return ctx.JSON(http.StatusOK, user)
}
