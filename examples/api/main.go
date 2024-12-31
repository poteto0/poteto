package main

import (
	"fmt"
	"net/http"

	"github.com/poteto0/poteto"
	"github.com/poteto0/poteto/middleware"
)

func main() {
	p := poteto.New()

	// Security Header
	p.Register(middleware.CamaraWithConfig(middleware.DefaultCamaraConfig))

	// CORS
	p.Register(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	p.GET("/", func(ctx poteto.Context) error {
		fmt.Println("hello worlda")
		return ctx.JSON(http.StatusOK, "Potato Simple Web framework")
	})

	p.Leaf("/users", func(userApi poteto.Leaf) {
		userApi.GET("/", func(ctx poteto.Context) error {
			return ctx.JSON(http.StatusOK, "get users")
		})
		userApi.GET("/:id", func(ctx poteto.Context) error {
			id, _ := ctx.PathParam("id")
			return ctx.JSON(http.StatusOK, id)
		})
	})

	p.Run("127.0.0.1:8000")
}
