package poteto

import (
	"fmt"
	"net/http"
)

type Response interface {
	WriteHeader(code int)
	Write(b []byte) (int, error)

	SetStatus(code int)
	Header() http.Header
}

type response struct {
	writer      http.ResponseWriter
	Status      int
	Size        int64
	isCommitted bool
}

func NewResponse(w http.ResponseWriter) Response {
	return &response{writer: w}
}

func (r *response) WriteHeader(code int) {
	if r.isCommitted {
		fmt.Println("response has already committed")
		return
	}

	r.Status = code
	r.writer.WriteHeader(r.Status)
	r.isCommitted = true
}

func (r *response) Write(b []byte) (int, error) {
	if !r.isCommitted {
		if r.Status == 0 {
			r.SetStatus(http.StatusOK)
		}
		r.WriteHeader(r.Status)
		r.isCommitted = true
	}

	n, err := r.writer.Write(b)
	r.Size += int64(n)

	return n, err
}

func (r *response) SetStatus(code int) {
	r.Status = code
}

func (r *response) Header() http.Header {
	return r.writer.Header()
}
