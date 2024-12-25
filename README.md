# Poteto

![](https://github.com/user-attachments/assets/7e503083-0af0-4b95-8277-46dfb8166cb9)

## Simple Web Framework of GoLang

We have confirmed that it works with various versions: go@1.21.x, go@1.22.x, go@1.23.x

```sh
go get github.com/poteto0/poteto@v0.26.3
go mod tidy
```

## Example App For Poteto

https://github.com/poteto0/poteto-sample-api/tree/main

## Poteto-Cli

We support cli tool. But if you doesn't like it, you can create poteto-app w/o cli of course.

```sh
go install github.com/poteto0/poteto/cmd/poteto-cli@v0.26.3
```

Create file.

```sh
poteto-cli new
```

fast mode.

```sh
poteto-cli new --fast
```

### Demo

https://github.com/user-attachments/assets/4b739964-1b4f-4913-b643-5984bf1ceae1

## Feature

### JSONRPCAdapter (`>=0.26.0`)

KeyNote: You can serve JSONRPC server easily.

```go
type (
  Calculator struct{}
  AdditionArgs   struct {
    Add, Added int
  }
)

func (tc *TestCalculator) Add(r *http.Request, args *AdditionArgs) int {
 return args.Add + args.Added
}

func main() {
  p := New()

  rpc := TestCalculator{}
  // you can access "/add/Calculator.Add"
  p.POST("/add", func(ctx Context) error {
    return PotetoJsonRPCAdapter[Calculator, AdditionArgs](ctx, &rpc)
  })

  p.Run("8080")
}
```

### Leaf router & middlewareTree (`>=0.21.0`)

```go
func main() {
	p := poteto.New()

	// Leaf >= 0.21.0
	p.Leaf("/users", func(userApi poteto.Leaf) {
		userApi.Register(middleware.CamaraWithConfig(middleware.DefaultCamaraConfig))
		userApi.GET("/", controller.UserHandler)
		userApi.GET("/:name", controller.UserIdHandler)
	})

	p.Run(":8000")
}
```

### Get RequestId Easily

```go
func handler(ctx poteto.Context) error {
	requestId := ctx.RequestId()
}
```

## How to use

```go:main.go
package main

import (
	"net/http"

	"github.com/poteto0/poteto"
	"github.com/poteto0/poteto/middleware"
)

func main() {
	p := poteto.New()

	// CORS
	p.Register(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		},
	))

	// Leaf >= 0.21.0
	p.Leaf("/users", func(userApi poteto.Leaf) {
		userApi.Register(middleware.CamaraWithConfig(middleware.DefaultCamaraConfig))
		userApi.GET("/", controller.UserHandler)
		userApi.GET("/:name", controller.UserNameHandler)
	})

	p.Run(":8000")
}

type User struct {
	Name any `json:"name"`
}

func UserHandler(ctx poteto.Context) error {
	user := User{
		Name: "user",
	}
	return ctx.JSON(http.StatusOK, user)
}

func UserNameHandler(ctx poteto.Context) error {
	name, _ := ctx.PathParam("name")
	user := User{
		Name: name,
	}
	return ctx.JSON(http.StatusOK, user)
}

```
