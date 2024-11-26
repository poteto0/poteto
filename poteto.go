package poteto

import (
	"net"
	"net/http"
	"strings"
	"sync"

	stdContext "context"

	"github.com/fatih/color"
	"github.com/poteto0/poteto/constant"
	"github.com/poteto0/poteto/utils"
)

type Poteto interface {
	// If requested, call this
	// you can make WithRequestIdOption false: you can faster request
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Run(addr string) error
	Stop(ctx stdContext.Context) error
	setupServer() error
	GET(path string, handler HandlerFunc) error
	POST(path string, handler HandlerFunc) error
	PUT(path string, handler HandlerFunc) error
	DELETE(path string, handler HandlerFunc) error
	Register(middlewares ...MiddlewareFunc)
	Combine(pattern string, middlewares ...MiddlewareFunc) *middlewareTree
	SetLogger(logger any)
}

type poteto struct {
	router         Router
	errorHandler   HttpErrorHandler
	middlewareTree MiddlewareTree
	logger         any
	cache          sync.Pool
	option         PotetoOption
	startupMutex   sync.RWMutex
	Server         http.Server
	Listener       net.Listener
}

func New() Poteto {
	return &poteto{
		router:         NewRouter([]string{"GET", "POST", "PUT", "DELETE"}),
		errorHandler:   &httpErrorHandler{},
		middlewareTree: NewMiddlewareTree(),
		option:         DefaultPotetoOption,
	}
}

func NewWithOption(option PotetoOption) Poteto {
	return &poteto{
		router:         NewRouter([]string{"GET", "POST", "PUT", "DELETE"}),
		errorHandler:   &httpErrorHandler{},
		middlewareTree: NewMiddlewareTree(),
		option:         option,
	}
}

func (p *poteto) initializeContext(w http.ResponseWriter, r *http.Request) *context {
	if ctx, ok := p.cache.Get().(*context); ok {
		ctx.Reset(w, r)
		return ctx
	}

	newCtx := NewContext(w, r).(*context)
	if p.logger != nil {
		newCtx.SetLogger(p.logger)
	}
	return newCtx
}

func (p *poteto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get from cache & reset context
	ctx := p.initializeContext(w, r)

	// default get & set RequestId
	if p.option.WithRequestId {
		reqId := ctx.RequestId()
		ctx.Set(constant.STORE_REQUEST_ID, reqId)
		if id := ctx.GetRequest().Header.Get(constant.HEADER_X_REQUEST_ID); id == "" {
			ctx.GetResponse().Header().Set(constant.HEADER_X_REQUEST_ID, reqId)
		}
	}

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

func (p *poteto) Run(addr string) error {
	p.startupMutex.Lock()

	if !strings.Contains(addr, constant.PARAM_PREFIX) {
		addr = constant.PARAM_PREFIX + addr
	}

	if err := p.setupServer(); err != nil {
		p.startupMutex.Unlock()
		return err
	}

	utils.PotetoPrint("server is available at http://localhost" + addr + "\n")

	p.startupMutex.Unlock()
	return p.Server.Serve(p.Listener)
}

func (p *poteto) setupServer() error {
	// Print Banner
	coloredBanner := color.HiGreenString(Banner)
	utils.PotetoPrint(coloredBanner)

	// setup Server
	p.Server.Handler = p
	if p.Listener == nil {
		ln, err := net.Listen(p.option.ListenerNetwork, p.Server.Addr)
		if err != nil {
			return err
		}

		p.Listener = ln
	}

	return nil
}

// Shutdown stops the server gracefully.
// It internally calls `http.Server#Shutdown()`.
func (p *poteto) Stop(ctx stdContext.Context) error {
	p.startupMutex.Lock()

	if err := p.Server.Shutdown(ctx); err != nil {
		p.startupMutex.Unlock()
		return err
	}

	p.startupMutex.Unlock()
	return nil
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

func (p *poteto) SetLogger(logger any) {
	p.logger = logger
}
