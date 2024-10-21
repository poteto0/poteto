package poteto

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/poteto0/poteto/constant"
)

type Poteto interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Run(addr string)
	GET(path string, handler HandlerFunc) error
	POST(path string, handler HandlerFunc) error
	PUT(path string, handler HandlerFunc) error
	DELETE(path string, handler HandlerFunc) error
	Register(middlewares ...MiddlewareFunc)
	applyMiddleware(handler HandlerFunc) HandlerFunc
}

type poteto struct {
	router      Router
	middlewares []MiddlewareFunc
}

func New() Poteto {
	return &poteto{
		router:      NewRouter([]string{"GET", "POST", "PUT", "DELETE"}),
		middlewares: []MiddlewareFunc{},
	}
}

func (p *poteto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	routes := p.router.GetRoutesByMethod(r.Method)

	queryParams := r.URL.Query()
	ctx.SetQueryParam(queryParams)

	targetRoute, httpParam := routes.Search(r.URL.Path)
	handler := targetRoute.GetHandler()

	if targetRoute == nil || handler == nil {
		ctx.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.SetPath(r.URL.Path)
	ctx.SetParam(constant.PARAM_TYPE_PATH, httpParam)
	handler = p.applyMiddleware(handler)
	handler(ctx)
}

func (p *poteto) applyMiddleware(handler HandlerFunc) HandlerFunc {
	for _, middleware := range p.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (p *poteto) Run(addr string) {
	// Print Banner
	coloredBanner := color.HiGreenString(Banner)
	fmt.Println(coloredBanner)

	if err := http.ListenAndServe(addr, p); err != nil {
		panic(err)
	}
}

func (p *poteto) GET(path string, handler HandlerFunc) error {
	return p.router.GET(path, handler)
}

func (p *poteto) POST(path string, handler HandlerFunc) error {
	return p.router.POST(path, handler)
}

func (p *poteto) PUT(path string, handler HandlerFunc) error {
	return p.router.PUT(path, handler)
}

func (p *poteto) DELETE(path string, handler HandlerFunc) error {
	return p.router.DELETE(path, handler)
}

func (p *poteto) Register(middlewares ...MiddlewareFunc) {
	p.middlewares = append(p.middlewares, middlewares...)
}
