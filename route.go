package poteto

import (
	"strings"
)

type Route interface {
	Search(path string) Route
	Insert(method, path string, handler func(ctx Context))

	GetHandler() func(ctx Context)
}

type route struct {
	key      string
	method   string
	children map[string]Route
	handler  func(ctx Context)
}

func NewRoute() Route {
	return &route{
		key:      "",
		method:   "",
		children: make(map[string]Route),
	}
}

func (r *route) Search(path string) Route {
	if len(r.key) == 0 && len(r.children) == 0 {
		return nil
	}

	currentRoute := r
	params := strings.Split(path, "/")

	for _, param := range params {
		if nextRoute, ok := currentRoute.children[param]; !ok {
			currentRoute = nextRoute.(*route)
		} else {
			return nil
		}
	}
	return currentRoute
}

func (r *route) Insert(method, path string, handler func(ctx Context)) {
	currentRoute := r
	params := strings.Split(path, "/")

	for _, param := range params {
		if nextRoute, ok := currentRoute.children[param]; !ok {
			currentRoute = nextRoute.(*route)
		} else {
			currentRoute.children[param] = &route{
				key:      param,
				method:   method,
				children: make(map[string]Route),
				handler:  handler,
			}
		}
	}
}

func (r *route) GetHandler() func(ctx Context) {
	return r.handler
}
