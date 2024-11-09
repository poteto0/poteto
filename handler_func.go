package poteto

type HandlerFunc func(ctx Context) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc
