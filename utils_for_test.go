package poteto

import "net/http"

func sampleMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx Context) error {
		res := ctx.GetResponse()

		res.Header().Set("Hello", "world")

		return next(ctx)
	}
}

func sampleMiddleware2(next HandlerFunc) HandlerFunc {
	return func(ctx Context) error {
		res := ctx.GetResponse()

		res.Header().Set("Hello2", "world2")

		return next(ctx)
	}
}

type user struct {
	Name string `json:"string"`
}

func getAllUserForTest(ctx Context) error {
	user := user{
		Name: "user",
	}
	return ctx.JSON(http.StatusOK, user)
}

func getAllUserForTestById(ctx Context) error {
	user := user{
		Name: "user1",
	}
	return ctx.JSON(http.StatusOK, user)
}

const (
	userJSON = `{"name":"poteto"}`
)
