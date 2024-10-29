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
	Combine(pattern string, middlewares ...MiddlewareFunc) *middlewareGroup
	applyMiddleware(mg MiddlewareGroup, handler HandlerFunc) HandlerFunc
}

type poteto struct {
	router          Router
	middlewares     []MiddlewareFunc
	errorHandler    HttpErrorHandler
	middlewareGroup MiddlewareGroup
}

func New() Poteto {
	return &poteto{
		router:          NewRouter([]string{"GET", "POST", "PUT", "DELETE"}),
		middlewares:     []MiddlewareFunc{},
		errorHandler:    &httpErrorHandler{},
		middlewareGroup: NewMiddlewareGroup(),
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

	// Search middleware
	mg := p.middlewareGroup.Search(r.URL.Path)
	handler = p.applyMiddleware(mg, handler)
	if err := handler(ctx); err != nil {
		p.errorHandler.HandleHttpError(err, ctx)
	}
}

func (p *poteto) applyMiddleware(mg MiddlewareGroup, handler HandlerFunc) HandlerFunc {
	return mg.ApplyMiddleware(handler)
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
	p.middlewareGroup.Insert("", middlewares...)
}

func (p *poteto) Combine(pattern string, middlewares ...MiddlewareFunc) *middlewareGroup {
	return p.middlewareGroup.Insert(pattern, middlewares...)
}
