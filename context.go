package poteto

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/goccy/go-json"
	"github.com/poteto0/poteto/constant"
	"github.com/poteto0/poteto/utils"
)

type Context interface {
	JSON(code int, value any) error
	WriteHeader(code int)
	SetQueryParam(queryParams url.Values)
	SetParam(paramType string, paramUnit ParamUnit)
	PathParam(key string) (string, bool)
	QueryParam(key string) (string, bool)
	Bind(object any) error
	GetPath() string
	SetPath(path string)
	GetResponse() *response
	GetRequest() *http.Request
	JsonSerialize(value any) error
	JsonDeserialize(object any) error
	NoContent() error
	Set(key string, val any)
	GetRemoteIP() (string, error)
	RegisterTrustIPRange(ranges *net.IPNet)
	GetIPFromXFFHeader() (string, error)
	Reset(w http.ResponseWriter, r *http.Request)
}

type context struct {
	response   Response
	request    *http.Request
	ipHandler  IPHandler
	path       string
	httpParams HttpParam
	store      map[string]any
	lock       sync.RWMutex

	// Method
	binder Binder
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return &context{
		response:   NewResponse(w),
		request:    r,
		ipHandler:  &ipHandler{isTrustPrivateIp: true},
		path:       "",
		httpParams: NewHttpParam(),
		binder:     NewBinder(),
	}
}

func (ctx *context) JSON(code int, value any) error {
	ctx.writeContentType(constant.APPLICATION_JSON)
	ctx.response.SetStatus(code)
	return ctx.JsonSerialize(value)
}

func (ctx *context) GetPath() string {
	return ctx.path
}

func (ctx *context) SetPath(path string) {
	ctx.path = path
}

func (ctx *context) SetQueryParam(queryParams url.Values) {
	if len(queryParams) > constant.MAX_QUERY_PARAM_LENGTH {
		fmt.Println("too many query params")
		return
	}

	for key, value := range queryParams {
		var paramUnit ParamUnit

		paramUnit = ParamUnit{key, utils.StrArrayToStr(value)}

		ctx.SetParam(constant.PARAM_TYPE_QUERY, paramUnit)
	}
}

func (ctx *context) SetParam(paramType string, paramUnit ParamUnit) {
	ctx.httpParams.AddParam(paramType, paramUnit)
}

func (ctx *context) PathParam(key string) (string, bool) {
	key = constant.PARAM_PREFIX + key
	return ctx.httpParams.GetParam(constant.PARAM_TYPE_PATH, key)
}

func (ctx *context) QueryParam(key string) (string, bool) {
	return ctx.httpParams.GetParam(constant.PARAM_TYPE_QUERY, key)
}

func (ctx *context) Bind(object any) error {
	err := ctx.binder.Bind(ctx, object)
	return err
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

func (ctx *context) JsonSerialize(value any) error {
	encoder := json.NewEncoder(ctx.GetResponse())
	return encoder.Encode(value)
}

func (ctx *context) JsonDeserialize(object any) error {
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(object)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return errors.New("error")
	} else if _, ok := err.(*json.SyntaxError); ok {
		return errors.New("error")
	}
	return err
}

func (c *context) NoContent() error {
	c.response.WriteHeader(http.StatusNoContent)
	return nil
}

func (ctx *context) Set(key string, val any) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if ctx.store == nil {
		ctx.store = make(map[string]any)
	}
	ctx.store[key] = val
}

func (ctx *context) GetRemoteIP() (string, error) {
	return ctx.ipHandler.GetRemoteIP(ctx)
}

func (ctx *context) RegisterTrustIPRange(ranges *net.IPNet) {
	ctx.ipHandler.RegisterTrustIPRange(ranges)
}

func (ctx *context) GetIPFromXFFHeader() (string, error) {
	return ctx.ipHandler.GetIPFromXFFHeader(ctx)
}

// using same binder
func (ctx *context) Reset(w http.ResponseWriter, r *http.Request) {
	ctx.request = r
	ctx.response = NewResponse(w)
	ctx.httpParams = NewHttpParam()
	ctx.store = make(map[string]any)
}
