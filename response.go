package poteto

import "net/http"

type Response interface {
	WriteHeader(code int)
	Write(b []byte) (int, error)

	SetStatus(code int)
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
		return
	}

	r.Status = code
}

func (r *response) Write(b []byte) (int, error) {
	if !r.isCommitted {
		if r.Status == 0 {
			r.SetStatus(http.StatusOK)
		}
		r.WriteHeader(r.Status)
	}

	n, err := r.writer.Write(b)
	r.Size += int64(n)

	return n, err
}

func (r *response) SetStatus(code int) {
	r.Status = code
}
