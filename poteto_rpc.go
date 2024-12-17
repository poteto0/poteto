package poteto

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/poteto0/poteto/utils"
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
// * Only Support "POST" method
func PotetoJsonRPCAdapter[T any](ctx Context, api T) error {
	if ctx.GetRequest().Method != http.MethodPost {
		return ctx.JSONRPCError(
			rpcErrorStatus,
			"parse error",
			"POST method excepted",
			0,
		)
	}

	if ctx.GetRequest().Body == nil {
		return ctx.JSONRPCError(
			rpcErrorStatus,
			"parse error",
			"No Post data",
			0,
		)
	}

	data := make(map[string]any)
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSONRPCError(
			rpcErrorStatus,
			"parse error",
			"error during decode json",
			0,
		)
	}

	id, ok := utils.AssertToInt(data["id"])
	if !ok {
		return ctx.JSONRPCError(
			rpcErrorStatusBadRequest,
			"BadRequest",
			"invalid id",
			0,
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

	fmt.Println(data["params"]) // TODO: map
	params, ok := data["params"].([]T)
	if !ok {
		return ctx.JSONRPCError(
			rpcErrorStatusBadRequest,
			"BadRequest",
			"invalid params",
			id,
		)
	}

	fmt.Println("はろー")
	fmt.Println(id)

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
