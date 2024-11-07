package poteto

import (
	"fmt"
	"net/http"
	"sync"

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
	Combine(pattern string, middlewares ...MiddlewareFunc) *middlewareTree
}

type poteto struct {
	router         Router
	errorHandler   HttpErrorHandler
	middlewareTree MiddlewareTree
	cache          sync.Pool
}

func New() Poteto {
	return &poteto{
		router:         NewRouter([]string{"GET", "POST", "PUT", "DELETE"}),
		errorHandler:   &httpErrorHandler{},
		middlewareTree: NewMiddlewareTree(),
	}
}

func (p *poteto) initializeContext(w http.ResponseWriter, r *http.Request) *context {
	if ctx, ok := p.cache.Get().(*context); ok {
		ctx.Reset(w, r)
		return ctx
	}
	return NewContext(w, r).(*context)
}

func (p *poteto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get from cache & reset context
	ctx := p.initializeContext(w, r)

	routes := p.router.GetRoutesByMethod(r.Method)

	targetRoute, httpParam := routes.Search(r.URL.Path)
	handler := targetRoute.GetHandler()

	if targetRoute == nil || handler == nil {
		ctx.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.SetQueryParam(r.URL.Query())
	ctx.SetPath(r.URL.Path)
	ctx.SetParam(constant.PARAM_TYPE_PATH, httpParam)

	// Search middleware
	middlewares := p.middlewareTree.SearchMiddlewares(r.URL.Path)
	handler = p.applyMiddleware(middlewares, handler)
	if err := handler(ctx); err != nil {
		p.errorHandler.HandleHttpError(err, ctx)
	}

	// cached context
	p.cache.Put(ctx)
}

func (p *poteto) applyMiddleware(middlewares []MiddlewareFunc, handler HandlerFunc) HandlerFunc {
	for _, middleware := range middlewares {
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
	p.middlewareTree.Insert("", middlewares...)
}

func (p *poteto) Combine(pattern string, middlewares ...MiddlewareFunc) *middlewareTree {
	return p.middlewareTree.Insert(pattern, middlewares...)
}
