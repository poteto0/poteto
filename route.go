package poteto

import (
	"strings"

	"github.com/poteto0/poteto/constant"
)

type Route interface {
	Search(path string) (*route, map[string]any)
	Insert(method, path string, handler HandlerFunc)

	GetHandler() HandlerFunc
}

type route struct {
	key           string
	method        string
	children      map[string]Route
	childParamKey string
	handler       HandlerFunc
}

func NewRoute() Route {
	return &route{
		key:           "",
		method:        "",
		children:      make(map[string]Route),
		childParamKey: "",
	}
}

func (r *route) Search(path string) (*route, map[string]any) {
	currentRoute := r
	params := strings.Split(path, "/")
	httpParam := map[string]any{}

	for i, param := range params {
		if param == "" {
			continue
		}

		if nextRoute, ok := currentRoute.children[param]; ok {
			currentRoute = nextRoute.(*route)
		} else {
			// last path includes url param ex: /users/:id
			if chParam := currentRoute.childParamKey; i == len(params)-1 && chParam != "" {
				if nextRoute, ok = currentRoute.children[chParam]; ok {
					currentRoute = nextRoute.(*route)
					httpParam[chParam] = param
				}
			} else {
				return nil, map[string]any{}
			}
		}
	}
	return currentRoute, httpParam
}

func (r *route) Insert(method, path string, handler HandlerFunc) {
	currentRoute := r
	params := strings.Split(path, "/")

	for i, param := range params {
		if param == "" {
			continue
		}

		if nextRoute := currentRoute.children[param]; nextRoute == nil {

			// last path includes url param ex: /users/:id
			if i == len(params)-1 && hasParamPrefix(param) {
				currentRoute.childParamKey = param
			}

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

func hasParamPrefix(param string) bool {
	return strings.HasPrefix(param, constant.PARAM_PREFIX)
}

func (r *route) GetHandler() HandlerFunc {
	return r.handler
}
