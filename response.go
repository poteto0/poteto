package poteto

import (
	"bytes"
	"net/http"
	"os"
)

type Response interface {
	WriteHeader(code int)
	Write(b []byte) (int, error)

	SetStatus(code int)
	Header() http.Header
}

type response struct {
	Writer      http.ResponseWriter
	Status      int
	Size        int64
	IsCommitted bool
}

func NewResponse(w http.ResponseWriter) Response {
	return &response{Writer: w}
}

func (r *response) WriteHeader(code int) {
	if r.IsCommitted {
		buf := bytes.NewBuffer([]byte("response has already committed\n"))
		buf.WriteTo(os.Stdout)
		return
	}

	r.Status = code
	r.Writer.WriteHeader(r.Status)
	r.IsCommitted = true
}

func (r *response) Write(b []byte) (int, error) {
	if !r.IsCommitted {
		if r.Status == 0 {
			r.SetStatus(http.StatusOK)
		}
		r.WriteHeader(r.Status)
		r.IsCommitted = true
	}

	n, err := r.Writer.Write(b)
	r.Size += int64(n)

	return n, err
}

func (r *response) SetStatus(code int) {
	r.Status = code
}

func (r *response) Header() http.Header {
	return r.Writer.Header()
}
