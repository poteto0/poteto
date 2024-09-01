package poteto

import "net/http"

type Context interface {
	JSON(code int, value any) error

	WriteHeader(code int)
	SetPath(path string)
	GetResponse() *response
}

type context struct {
	response       Response
	request        *http.Request
	path           string
	jsonSerializer JsonSerializer
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return &context{
		response: NewResponse(w),
		request:  r,
		path:     "",
	}
}

func (ctx *context) JSON(code int, value any) error {
	ctx.response.SetStatus(code)
	return ctx.jsonSerializer.Serialize(ctx, value)
}

func (ctx *context) SetPath(path string) {
	ctx.path = path
}

func (ctx *context) WriteHeader(code int) {
	ctx.response.WriteHeader(code)
}

func (ctx *context) GetResponse() *response {
	return ctx.response.(*response)
}
