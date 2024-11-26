package poteto

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	p := New()
	p.GET("/users", getAllUserForTest)
	p.GET("/users/:id", getAllUserForTestById)
	logger := func(msg string) {
		return
	}
	p.SetLogger(logger)

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

func BenchmarkServeHTTP(b *testing.B) {
	p := New()
	p.GET("/users/:id", getAllUserForTestById)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1", strings.NewReader(userJSON))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ServeHTTP(w, req)
	}
}

func BenchmarkServeHTTPWORequestId(b *testing.B) {
	option := PotetoOption{
		WithRequestId: false,
	}
	p := NewWithOption(option)
	p.GET("/users/:id", getAllUserForTestById)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1", strings.NewReader(userJSON))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ServeHTTP(w, req)
	}
}

func TestServeHTTPWithMiddleware(t *testing.T) {
	p := New()
	p.Register(sampleMiddleware)
	p.GET("/users", getAllUserForTest)

	group := p.Combine("/tests")
	group.Register(sampleMiddleware2)
	p.GET("/tests", getAllUserForTest)
	p.GET("/tests/:id", getAllUserForTestById)

	tests := []struct {
		name          string
		reqMethod     string
		reqUrl        string
		worked        bool
		expectedKey   string
		expectedValue string
	}{
		{
			"Test Middleware",
			"GET",
			"/users",
			true,
			"Hello",
			"world",
		},
		{
			"Test MiddlewareTree",
			"GET",
			"/tests",
			true,
			"Hello2",
			"world2",
		},
		{
			"Test if MiddlewareTree includes all middleware",
			"GET",
			"/tests",
			true,
			"Hello",
			"world",
		},
		{
			"Test parent Middleware",
			"GET",
			"/tests/1",
			true,
			"Hello2",
			"world2",
		},
		{
			"Test not apply middleware without group",
			"GET",
			"/users",
			false,
			"Hello2",
			"world2",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(it.reqMethod, it.reqUrl, nil)

			p.ServeHTTP(w, req)
			header := w.Header()
			if it.worked {
				if header[it.expectedKey][0] != it.expectedValue {
					t.Errorf("Unmatched")
				}
			} else {
				if header[it.expectedKey] != nil {
					t.Errorf("Unmatched")
				}
			}
		})
	}
}

func TestLeafHandler(t *testing.T) {
	p := New()

	p.Leaf("/users", func(leaf Leaf) {
		leaf.Register(sampleMiddleware)
		leaf.GET("/", getAllUserForTest)
		leaf.POST("/create", getAllUserForTest)
		leaf.PUT("/change", getAllUserForTest)
		leaf.DELETE("/delete", getAllUserForTest)
	})

	tests := []struct {
		name          string
		reqMethod     string
		reqUrl        string
		expectedKey   string
		expectedValue string
		expectedRes   string
	}{
		{
			"Test Get",
			"GET",
			"/users",
			"Hello",
			"world",
			`{"string":"user"}`,
		},
		{
			"Test Post",
			"POST",
			"/users/create",
			"Hello",
			"world",
			`{"string":"user"}`,
		},
		{
			"Test Put",
			"PUT",
			"/users/change",
			"Hello",
			"world",
			`{"string":"user"}`,
		},
		{
			"Test Delete",
			"DELETE",
			"/users/delete",
			"Hello",
			"world",
			`{"string":"user"}`,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(it.reqMethod, it.reqUrl, nil)

			p.ServeHTTP(w, req)
			if w.Header()[it.expectedKey][0] != it.expectedValue {
				t.Errorf("Unmatched")
			}
			if w.Body.String()[0:10] != it.expectedRes[0:10] {
				t.Errorf("Unmatched")
			}
		})
	}
}
