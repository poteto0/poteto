package poteto

import (
	"strings"
)

type Route interface {
	Search(path string) *route
	Insert(method, path string, handler HandlerFunc)

	GetHandler() HandlerFunc
}

type route struct {
	key      string
	method   string
	children map[string]Route
	handler  HandlerFunc
}

func NewRoute() Route {
	return &route{
		key:      "",
		method:   "",
		children: make(map[string]Route),
	}
}

func (r *route) Search(path string) *route {
	currentRoute := r
	params := strings.Split(path, "/")

	for _, param := range params {
		if param == "" {
			continue
		}

		if nextRoute, ok := currentRoute.children[param]; ok {
			currentRoute = nextRoute.(*route)
		} else {
			return nil
		}
	}
	return currentRoute
}

func (r *route) Insert(method, path string, handler HandlerFunc) {
	currentRoute := r
	params := strings.Split(path, "/")

	for _, param := range params {
		if param == "" {
			continue
		}

		if nextRoute := currentRoute.children[param]; nextRoute == nil {
			currentRoute.children[param] = &route{
				key:      param,
				method:   method,
				children: make(map[string]Route),
			}
		}
		currentRoute = currentRoute.children[param].(*route)
	}
	currentRoute.handler = handler
}

func (r *route) GetHandler() HandlerFunc {
	return r.handler
}
