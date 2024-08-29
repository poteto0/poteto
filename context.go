package poteto

import "net/http"

type Context struct {
	w    http.ResponseWriter
	r    *http.Request
	path string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:    w,
		r:    r,
		path: "",
	}
}

func (ctx *Context) SetPath(path string) {
	ctx.path = path
}
