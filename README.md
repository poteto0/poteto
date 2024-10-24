# Poteto

![](https://github.com/user-attachments/assets/7e503083-0af0-4b95-8277-46dfb8166cb9)

## Simple Web Framework of GoLang

```sh
go get github.com/poteto0/poteto@v0.9.5
```

```go:main.go
package main

import (
	"net/http"

	"github.com/poteto0/poteto"
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

	e.Register(middleware.CamaraWithConfig(middleware.DefaultCamaraConfig))

	p.GET("/users", UserHandler)
	p.Get("/users/:id", UserIdHandler)
	p.Run(":8000")
}

type User struct {
	Name string `json:"string"`
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
