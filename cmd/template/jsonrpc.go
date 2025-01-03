package template

var JSONRPCTemplate = `
package main

import (
	"net/http"

	"github.com/poteto-go/poteto"
)

type (
	Calculator   struct{}
	AdditionArgs struct {
		Add, Added int
	}
)

func (c *Calculator) Add(r *http.Request, args *AdditionArgs) int {
	return args.Add + args.Added
}

func main() {
	p := poteto.New()

	rpc := Calculator{}
	p.POST("/add", func(ctx poteto.Context) error {
		return poteto.PotetoJsonRPCAdapter[Calculator, AdditionArgs](ctx, &rpc)
	})

	p.Run("127.0.0.1:8000")
}
`
