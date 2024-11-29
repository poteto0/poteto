package poteto

import (
	"errors"
	"net/http"
)

type Router interface {
	add(method, path string, handler HandlerFunc) error
	GET(path string, handler HandlerFunc) error
	POST(path string, handler HandlerFunc) error
	PUT(path string, handler HandlerFunc) error
	DELETE(path string, handler HandlerFunc) error

	GetRoutesByMethod(method string) *route
}

type router struct {
	routesGET    Route
	routesPOST   Route
	routesPUT    Route
	routesDELETE Route
}

func NewRouter() Router {
	return &router{
		routesGET:    NewRoute(),
		routesPOST:   NewRoute(),
		routesPUT:    NewRoute(),
		routesDELETE: NewRoute(),
	}
}

func (r *router) add(method, path string, handler HandlerFunc) error {
	routes := r.GetRoutesByMethod(method)
	if routes == nil {
		return errors.New("Unexpected method error: [GET, POST, PUT, DELETE]")
	}

	if that_route, _ := routes.Search(path); that_route != nil {
		return errors.New("[" + method + "] " + path + " is already used")
	}

	routes.Insert(path, handler)
	return nil
}

func (r *router) GET(path string, handler HandlerFunc) error {
	return r.add(http.MethodGet, path, handler)
}

func (r *router) POST(path string, handler HandlerFunc) error {
	return r.add(http.MethodPost, path, handler)
}

func (r *router) PUT(path string, handler HandlerFunc) error {
	return r.add(http.MethodPut, path, handler)
}

func (r *router) DELETE(path string, handler HandlerFunc) error {
	return r.add(http.MethodDelete, path, handler)
}

func (r *router) GetRoutesByMethod(method string) *route {
	switch method {
	case http.MethodGet:
		return r.routesGET.(*route)
	case http.MethodPost:
		return r.routesPOST.(*route)
	case http.MethodPut:
		return r.routesPUT.(*route)
	case http.MethodDelete:
		return r.routesDELETE.(*route)
	default:
		return nil
	}
}
