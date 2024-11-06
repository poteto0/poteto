package poteto

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/poteto0/poteto/constant"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	context := NewContext(w, req).(*context)

	tests := []struct {
		name     string
		code     int
		val      testVal
		expected string
	}{
		{"status ok & can serialize",
			http.StatusOK,
			testVal{Name: "test", Val: "val"},
			`{"name":"test","val":"val"}`,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			context.JSON(it.code, it.val)
			if body := w.Body.String(); body[0:27] != it.expected[0:27] {
				t.Errorf("FATAL: context json")
			}

			if header := w.Header(); header[constant.HEADER_CONTENT_TYPE][0] != constant.APPLICATION_JSON {
				t.Errorf("FATAL: wrong content-type")
			}
		})
	}
}

func TestQueryParam(t *testing.T) {
	url := "https://example.com?hello=world"

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", url, nil)
	ctx := NewContext(w, req).(*context)

	queryParams := req.URL.Query()
	ctx.SetQueryParam(queryParams)

	queryParam1 := ctx.QueryParam("hello")
	if queryParam1 != "world" {
		t.Errorf("Cannot Get Query Param")
	}

	queryParam2 := ctx.QueryParam("unknown")
	if queryParam2 != nil {
		t.Errorf("Cannot Get nil If Unknown key")
	}
}

func TestPathParam(t *testing.T) {
	url := "https://example.com/users/1"

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", url, nil)
	ctx := NewContext(w, req).(*context)

	ctx.SetParam(constant.PARAM_TYPE_PATH, ParamUnit{key: ":id", value: "1"})

	tests := []struct {
		name     string
		key      string
		expected any
	}{
		{"Can get PathParam", "id", "1"},
		{"If unexpected key", "unexpected", nil},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			param := ctx.PathParam(it.key)

			if param != it.expected {
				t.Errorf("Unmatched")
			}
		})

	}
}

func TestSetPath(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com", nil)
	ctx := NewContext(w, req).(*context)

	expected := "http://expected.com"
	ctx.SetPath(expected)
	if ctx.path != expected {
		t.Errorf("Not Matched")
	}
}

func BenchmarkJSON(b *testing.B) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://example.com", strings.NewReader(userJSON))
	ctx := NewContext(w, req).(*context)

	testUser := user{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx.JSON(http.StatusOK, testUser)
	}
}
