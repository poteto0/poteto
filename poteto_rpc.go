package poteto

import (
	"net/http"
	"reflect"
)

// TODO: error status code
var (
	rpcVersion               = "2.0"
	rpcErrorStatus           = -32700
	rpcErrorStatusBadRequest = -32600
	rpcErrorStatusNotFound   = -32601
)

// inspired by
// https://github.com/kanocz/goginjsonrpc/blob/master/jsonrpc.go
// * Just Method Post
func PotetoJsonRPCAdapter[T any](ctx Context, api T) error {
	if ctx.GetRequest().Method != http.MethodPost {
		return ctx.JSONRPCError(
			rpcErrorStatus,
			"parse error",
			"POST method excepted",
			"",
		)
	}

	if ctx.GetRequest().Body == nil {
		return ctx.JSONRPCError(
			rpcErrorStatus,
			"parse error",
			"No Post data",
			"",
		)
	}

	data := make(map[string]any)
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSONRPCError(
			rpcErrorStatus,
			"parse error",
			"error during decode json",
			"",
		)
	}

	id, ok := data["id"].(string)
	if !ok {
		return ctx.JSONRPCError(
			rpcErrorStatusBadRequest,
			"BadRequest",
			"invalid id",
			"",
		)
	}

	version, ok := data["jsonrpc"]
	if !ok || version != rpcVersion {
		return ctx.JSONRPCError(
			rpcErrorStatusBadRequest,
			"BadRequest",
			"version of jsonrpc is not 2.0",
			id,
		)
	}

	method, ok := data["method"].(string)
	if !ok {
		return ctx.JSONRPCError(
			rpcErrorStatusBadRequest,
			"BadRequest",
			"invalid method",
			id,
		)
	}

	params, ok := data["params"].([]T)
	if !ok {
		return ctx.JSONRPCError(
			rpcErrorStatusBadRequest,
			"BadRequest",
			"invalid params",
			id,
		)
	}

	call := reflect.ValueOf(api).MethodByName(method)
	if !call.IsValid() {
		return ctx.JSONRPCError(
			rpcErrorStatusNotFound,
			"NotFound",
			"Method is not found",
			id,
		)
	}

	if call.Type().NumIn() != len(params) {
		return ctx.JSONRPCError(
			rpcErrorStatusBadRequest,
			"BadRequest",
			"invalid params length",
			id,
		)
	}

	args := make([]reflect.Value, len(params))
	result := call.Call(args)

	if len(result) > 0 {
		return ctx.JSON(http.StatusOK, map[string]any{
			"result":  result[0].Interface(),
			"jsonrpc": "2.0",
			"id":      id,
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"result":  nil,
		"jsonrpc": "2.0",
		"id":      id,
	})
}
