package poteto

import (
	"strings"
)

type Route struct {
	key      string
	method   string
	children map[string]*Route
	handler  func(ctx *Context) // TODO: Context
}

func NewRoute() *Route {
	return &Route{
		key:      "",
		method:   "",
		children: make(map[string]*Route),
	}
}

func (r *Route) Search(path string) *Route {
	if len(r.key) == 0 && len(r.children) == 0 {
		return nil
	}

	currentRoute := r
	params := strings.Split(path, "/")

	for _, param := range params {
		if nextRoute, ok := currentRoute.children[param]; !ok {
			currentRoute = nextRoute
		} else {
			return nil
		}
	}
	return currentRoute
}

func (r *Route) Insert(method, path string, handler func(ctx *Context)) {
	currentRoute := r
	params := strings.Split(path, "/")

	for _, param := range params {
		if nextRoute, ok := currentRoute.children[param]; !ok {
			currentRoute = nextRoute
		} else {
			currentRoute.children[param] = &Route{
				key:      param,
				method:   method,
				children: make(map[string]*Route),
				handler:  handler,
			}
		}
	}
}
