package poteto

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poteto0/poteto/constant"
)

type TestVal struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

type TestExpected struct {
	Code int   `json:"code"`
	Val  error `json:"val"`
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	context := NewContext(w, req).(*context)

	tests := []struct {
		name     string
		code     int
		val      TestVal
		expected string
	}{
		{"status ok & can serialize",
			http.StatusOK,
			TestVal{Name: "test", Val: "val"},
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
