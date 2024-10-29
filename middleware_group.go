package poteto

import (
	"strings"
)

/*
 * Middleware Group
 * Trie Tree of path -> middleware
 * If you want to apply into all path, Apply to "" path
 * Refer, if once found.
 */

type MiddlewareGroup interface {
	Search(pattern string) *middlewareGroup
	Insert(pattern string, middlewares ...MiddlewareFunc) *middlewareGroup
	ApplyMiddleware(handler HandlerFunc) HandlerFunc
	Register(middlewares ...MiddlewareFunc)
}

type middlewareGroup struct {
	children    map[string]MiddlewareGroup
	middlewares []MiddlewareFunc
	key         string
}

func NewMiddlewareGroup() MiddlewareGroup {
	return &middlewareGroup{
		children: make(map[string]MiddlewareGroup),
	}
}

func (mg *middlewareGroup) Search(pattern string) *middlewareGroup {
	currentNode := mg
	patterns := strings.Split(pattern, "/")

	for _, p := range patterns {
		if p == "" {
			continue
		}

		if nextNode, ok := currentNode.children[p]; ok {
			currentNode = nextNode.(*middlewareGroup)
		} else {
			// if found ever
			// You got Middleware Group
			break
		}
	}

	return currentNode
}

func (mg *middlewareGroup) Insert(pattern string, middlewares ...MiddlewareFunc) *middlewareGroup {
	currentNode := mg
	patterns := strings.Split(pattern, "/")

	for _, p := range patterns {
		if p == "" {
			continue
		}

		if _, ok := currentNode.children[p]; !ok {
			currentNode.children[p] = &middlewareGroup{
				children:    make(map[string]MiddlewareGroup),
				middlewares: []MiddlewareFunc{},
				key:         p,
			}
		}
		currentNode = currentNode.children[p].(*middlewareGroup)
	}
	currentNode.Register(middlewares...)
	return currentNode
}

func (mg *middlewareGroup) Register(middlewares ...MiddlewareFunc) {
	mg.middlewares = append(mg.middlewares, middlewares...)
}

func (mg *middlewareGroup) ApplyMiddleware(handler HandlerFunc) HandlerFunc {
	for _, middleware := range mg.middlewares {
		handler = middleware(handler)
	}
	return handler
}
