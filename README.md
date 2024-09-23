# Poteto

Simple Web Framework of GoLang

```sh
go get github.com/poteto0/poteto@v0.3.2
```

```go:main.go
package main

import (
	"net/http"

	"github.com/poteto0/poteto"
)

func main() {
	p := poteto.New()

	p.GET("/users", UserHandler)
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

```
