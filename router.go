package poteto

import "errors"

type Router interface {
	add(method, path string, handler HandlerFunc) error
	GET(path string, handler HandlerFunc) error
	POST(path string, handler HandlerFunc) error
	PUT(path string, handler HandlerFunc) error
	DELETE(path string, handler HandlerFunc) error

	GetRoutesByMethod(method string) *route
}

type router struct {
	routes map[string]Route
}

func NewRouter(methods []string) Router {
	rs := make(map[string]Route)
	for _, method := range methods {
		rs[method] = NewRoute()
	}

	return &router{
		routes: rs,
	}
}

func (r *router) add(method, path string, handler HandlerFunc) error {
	routes := r.GetRoutesByMethod(method)

	if that_route := routes.Search(path); that_route != nil {
		return errors.New("[" + method + "] " + path + " is already used")
	}

	routes.Insert(method, path, handler)
	return nil
}

func (r *router) GET(path string, handler HandlerFunc) error {
	return r.add("GET", path, handler)
}

func (r *router) POST(path string, handler HandlerFunc) error {
	return r.add("POST", path, handler)
}

func (r *router) PUT(path string, handler HandlerFunc) error {
	return r.add("PUT", path, handler)
}

func (r *router) DELETE(path string, handler HandlerFunc) error {
	return r.add("DELETE", path, handler)
}

func (r *router) GetRoutesByMethod(method string) *route {
	return r.routes[method].(*route)
}
