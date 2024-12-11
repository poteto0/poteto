package poteto

import (
	"errors"
	"net/http"
	"strings"
)

type Router interface {
	add(method, path string, handler HandlerFunc) error
	GET(path string, handler HandlerFunc) error
	POST(path string, handler HandlerFunc) error
	PUT(path string, handler HandlerFunc) error
	PATCH(path string, handler HandlerFunc) error
	DELETE(path string, handler HandlerFunc) error
	HEAD(path string, handler HandlerFunc) error
	OPTIONS(path string, handler HandlerFunc) error
	TRACE(path string, handler HandlerFunc) error
	CONNECT(path string, handler HandlerFunc) error

	GetRoutesByMethod(method string) *route
}

type router struct {
	routesGET     Route
	routesPOST    Route
	routesPUT     Route
	routesPATCH   Route
	routesDELETE  Route
	routesHEAD    Route
	routesOPTIONS Route
	routesTRACE   Route
	routesCONNECT Route
}

func NewRouter() Router {
	return &router{
		routesGET:     NewRoute(),
		routesPOST:    NewRoute(),
		routesPUT:     NewRoute(),
		routesPATCH:   NewRoute(),
		routesDELETE:  NewRoute(),
		routesHEAD:    NewRoute(),
		routesOPTIONS: NewRoute(),
		routesTRACE:   NewRoute(),
		routesCONNECT: NewRoute(),
	}
}

func (r *router) add(method, path string, handler HandlerFunc) error {
	routes := r.GetRoutesByMethod(method)
	if routes == nil {
		return errors.New("Unexpected method error: " + method)
	}

	if that_route, _ := routes.Search(path); that_route != nil {
		if path == "/" {
			that_route.handler = handler
			return nil
		}
		return errors.New("[" + method + "] " + path + " is already used")
	}

	// "/users/" -> "/users"
	// if just "/" -> handler set by above
	path = strings.TrimSuffix(path, "/")

	routes.Insert(path, handler)
	return nil
}

// These are router Method
// Seems redundant, but you can register your own router with poteto
// And call it with `Poteto.GET()` etc.
func (r *router) GET(path string, handler HandlerFunc) error {
	return r.add(http.MethodGet, path, handler)
}

func (r *router) POST(path string, handler HandlerFunc) error {
	return r.add(http.MethodPost, path, handler)
}

func (r *router) PUT(path string, handler HandlerFunc) error {
	return r.add(http.MethodPut, path, handler)
}

func (r *router) PATCH(path string, handler HandlerFunc) error {
	return r.add(http.MethodPatch, path, handler)
}

func (r *router) DELETE(path string, handler HandlerFunc) error {
	return r.add(http.MethodDelete, path, handler)
}

func (r *router) HEAD(path string, handler HandlerFunc) error {
	return r.add(http.MethodHead, path, handler)
}

func (r *router) OPTIONS(path string, handler HandlerFunc) error {
	return r.add(http.MethodOptions, path, handler)
}

func (r *router) TRACE(path string, handler HandlerFunc) error {
	return r.add(http.MethodTrace, path, handler)
}

func (r *router) CONNECT(path string, handler HandlerFunc) error {
	return r.add(http.MethodConnect, path, handler)
}

func (r *router) GetRoutesByMethod(method string) *route {
	switch method {
	case http.MethodGet:
		return r.routesGET.(*route)
	case http.MethodPost:
		return r.routesPOST.(*route)
	case http.MethodPut:
		return r.routesPUT.(*route)
	case http.MethodPatch:
		return r.routesPATCH.(*route)
	case http.MethodDelete:
		return r.routesDELETE.(*route)
	case http.MethodHead:
		return r.routesHEAD.(*route)
	case http.MethodOptions:
		return r.routesOPTIONS.(*route)
	case http.MethodTrace:
		return r.routesTRACE.(*route)
	case http.MethodConnect:
		return r.routesCONNECT.(*route)
	default:
		return nil
	}
}
