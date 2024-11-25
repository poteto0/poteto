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
	params := strings.Split(path, "/")
	last := len(params) - 1
	var httpParam ParamUnit

	for i, param := range params {
		if param == "" {
			continue
		}

		if nextRoute, ok := currentRoute.children[param]; ok {
			currentRoute = nextRoute.(*route)
		} else {
			// last path includes url param ex: /users/:id
			if chParam := currentRoute.childParamKey; i == last && chParam != "" {
				if nextRoute, ok = currentRoute.children[chParam]; ok {
					currentRoute = nextRoute.(*route)
					httpParam = ParamUnit{key: chParam, value: param}
				}
			} else {
				return nil, ParamUnit{}
			}
		}
	}
	return currentRoute, httpParam
}

func (r *route) Insert(path string, handler HandlerFunc) {
	currentRoute := r
	params := strings.Split(path, "/")
	last := len(params) - 1

	for i, param := range params {
		if param == "" {
			continue
		}

		if nextRoute := currentRoute.children[param]; nextRoute == nil {

			// last path includes url param ex: /users/:id
			if i == last && hasParamPrefix(param) {
				currentRoute.childParamKey = param
			}

			currentRoute.children[param] = &route{
				key:      param,
				children: make(map[string]Route),
			}
		}
		currentRoute = currentRoute.children[param].(*route)
	}

	if currentRoute.handler != nil {
		coloredWarn := color.HiRedString(fmt.Sprintf("Handler Collision on %s \n", utils.StrArrayToStr(params)))
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
