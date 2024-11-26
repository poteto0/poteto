package poteto

import (
	"net/http/httptest"
	"testing"
)

func TestLeaf(t *testing.T) {
	p := New()

	leaf := NewLeaf(p, "/users")
	leaf.Register(sampleMiddleware)
	leaf.GET("/", getAllUserForTest)
	leaf.POST("/create", getAllUserForTest)
	leaf.PUT("/change", getAllUserForTest)
	leaf.DELETE("/delete", getAllUserForTest)

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
