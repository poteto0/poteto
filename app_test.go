package poteto

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
