package poteto

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w).(*response)

	tests := []struct {
		name        string
		IsCommitted bool
		expected    int
	}{
		{"Test not committed case", false, http.StatusOK},
		{"Test committed case", true, 0},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			resp.IsCommitted = it.IsCommitted

			resp.WriteHeader(http.StatusOK)

			if resp.Status != http.StatusOK {
				t.Errorf(
					"Unmatched actual(%d) -> expected(%d)",
					resp.Status,
					it.expected,
				)
			}
		})
	}

}

func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w).(*response)

	tests := []struct {
		name        string
		IsCommitted bool
		b           []byte
		expected    int
	}{
		{"write not committed response", false, []byte("Hello"), 5},
		{"don't write committed reponse", true, []byte(""), 0},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			resp.IsCommitted = it.IsCommitted
			n, _ := resp.Write(it.b)
			if n != it.expected {
				t.Errorf(
					"Unmatched actual(%d) -> expected(%d)",
					n,
					it.expected,
				)
			}
		})
	}
}
