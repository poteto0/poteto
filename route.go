package poteto

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/poteto0/poteto/constant"
	"github.com/poteto0/poteto/utils"
)

type Route interface {
	Search(path string) (*route, ParamUnit)
	Insert(path string, handler HandlerFunc)

	GetHandler() HandlerFunc
}

type route struct {
	key           string
	children      map[string]Route
	childParamKey string
	handler       HandlerFunc
}

func NewRoute() Route {
	return &route{
		key:           "",
		children:      make(map[string]Route),
		childParamKey: "",
	}
}

func (r *route) Search(path string) (*route, ParamUnit) {
	currentRoute := r
	rightPath := path[1:]
	param := ""
	var httpParam ParamUnit

	if rightPath == "" {
		return currentRoute, ParamUnit{}
	}

	// optimized router insert
	// https://github.com/poteto0/poteto/issues/113
	for {
		id := strings.Index(rightPath, "/")
		if id < 0 {
			param = rightPath
		} else {
			param = rightPath[:id]
			rightPath = rightPath[(id + 1):]
		}

		if nextRoute, ok := currentRoute.children[param]; ok {
			currentRoute = nextRoute.(*route)
		} else {
			// last path includes url param ex: /users/:id
			if chParam := currentRoute.childParamKey; (id < 0) && chParam != "" {
				if nextRoute, ok = currentRoute.children[chParam]; ok {
					currentRoute = nextRoute.(*route)
					httpParam = ParamUnit{key: chParam, value: param}
				}
			} else {
				return nil, ParamUnit{}
			}
		}

		if id < 0 {
			break
		}
	}

	return currentRoute, httpParam
}

func (r *route) Insert(path string, handler HandlerFunc) {
	currentRoute := r
	rightPath := path[1:]
	param := ""

	// optimized router insert
	// https://github.com/poteto0/poteto/issues/113
	for {
		id := strings.Index(rightPath, "/")
		if id < 0 { // means last
			param = rightPath
		} else {
			param = rightPath[:id]
			rightPath = rightPath[(id + 1):]
		}

		if nextRoute := currentRoute.children[param]; nextRoute == nil {
			// last path includes url param ex: /users/:id
			if (id < 0) && hasParamPrefix(param) {
				currentRoute.childParamKey = param
			}

			currentRoute.children[param] = &route{
				key:      param,
				children: make(map[string]Route),
			}
		}
		currentRoute = currentRoute.children[param].(*route)

		if id < 0 {
			break
		}
	}

	if currentRoute.handler != nil {
		coloredWarn := color.HiRedString(fmt.Sprintf("Handler Collision on %s \n", path))
		utils.PotetoPrint(coloredWarn)
		return
	}

	currentRoute.handler = handler
}

func hasParamPrefix(param string) bool {
	return strings.HasPrefix(param, constant.PARAM_PREFIX)
}

func (r *route) GetHandler() HandlerFunc {
	return r.handler
}
