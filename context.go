package poteto

import (
	"encoding/json"
	"net/http"

	"github.com/poteto0/poteto/constant"
)

type Context interface {
	JSON(code int, value any) error

	WriteHeader(code int)
	writeContentType(value string)
	SetPath(path string)
	GetResponse() *response
	GetRequest() *http.Request
	GetRequestHeaderValue(key string) string
	JsonSerialize(value any) error
	NoContent() error
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
	ctx.writeContentType(constant.APPLICATION_JSON)
	ctx.response.SetStatus(code)
	return ctx.JsonSerialize(value)
}

func (ctx *context) SetPath(path string) {
	ctx.path = path
}

func (ctx *context) WriteHeader(code int) {
	ctx.response.WriteHeader(code)
}

func (ctx *context) writeContentType(value string) {
	header := ctx.response.Header()

	if header.Get(constant.HEADER_CONTENT_TYPE) == "" {
		header.Set(constant.HEADER_CONTENT_TYPE, value)
	}
}

func (ctx *context) GetResponse() *response {
	return ctx.response.(*response)
}

func (ctx *context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *context) GetRequestHeaderValue(key string) string {
	return ctx.request.Header.Get(key)
}

func (ctx *context) JsonSerialize(value any) error {
	encoder := json.NewEncoder(ctx.GetResponse())
	return encoder.Encode(value)
}

func (c *context) NoContent() error {
	c.response.WriteHeader(http.StatusNoContent)
	return nil
}
