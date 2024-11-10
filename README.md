# Poteto

![](https://github.com/user-attachments/assets/7e503083-0af0-4b95-8277-46dfb8166cb9)

## Simple Web Framework of GoLang

```sh
go get github.com/poteto0/poteto@v0.15.3
go mod tidy
```

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

	p.GET("/example", UserHandler)

	// Group of Middleware
	userGroup := p.Combine("/users")
	userGroup.Register(middleware.CamaraWithConfig(middleware.DefaultCamaraConfig))

	p.GET("/users", UserHandler)
	p.Get("/users/:id", UserIdHandler)
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

func UserIdHandler(ctx poteto.Context) error {
	user := User{
		Name: ctx.PathParam("id")
	}
	return ctx.JSON(http.StatusOK, user)
}

```
