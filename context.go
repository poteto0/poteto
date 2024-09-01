package poteto

import "net/http"

type Context interface {
	SetPath(path string)
}

type context struct {
	response http.ResponseWriter
	request  *http.Request
	path     string
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return &context{
		response: w,
		request:  r,
		path:     "",
	}
}

func (ctx *context) SetPath(path string) {
	ctx.path = path
}
