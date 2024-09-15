package poteto

import (
	"encoding/json"
	"net/http"
)

type Context interface {
	JSON(code int, value any) error

	WriteHeader(code int)
	SetPath(path string)
	GetResponse() *response
	JsonSerialize(value any) error
}

type context struct {
	response Response
	request  *http.Request
	path     string
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
	return ctx.JsonSerialize(value)
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

func (ctx *context) JsonSerialize(value any) error {
	encoder := json.NewEncoder(ctx.GetResponse())
	return encoder.Encode(value)
}
