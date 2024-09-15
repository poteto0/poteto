package poteto

import (
	"net/http/httptest"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w).(*response)

	resp.WriteHeader(200)

	if resp.Status != 200 {
		t.Errorf("Cannot write status of header")
	}
}

func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()
	resp := NewResponse(w).(*response)

	tests := []struct {
		name        string
		isCommitted bool
		b           []byte
		expected    int
	}{
		{"write not committed response", false, []byte("Hello"), 5},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			resp.isCommitted = it.isCommitted
			n, _ := resp.Write(it.b)
			if n != it.expected {
				t.Errorf("FATAL write response")
			}
		})
	}
}
