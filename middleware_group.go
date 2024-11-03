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
	SearchMiddlewares(pattern string) []MiddlewareFunc
	Insert(pattern string, middlewares ...MiddlewareFunc) *middlewareGroup
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

func (mg *middlewareGroup) SearchMiddlewares(pattern string) []MiddlewareFunc {
	middlewares := []MiddlewareFunc{}
	currentNode := mg
	middlewares = append(middlewares, mg.middlewares...)
	patterns := strings.Split(pattern, "/")

	for _, p := range patterns {
		if p == "" {
			continue
		}

		if nextNode, ok := currentNode.children[p]; ok {
			currentNode = nextNode.(*middlewareGroup)
			middlewares = append(middlewares, currentNode.middlewares...)
		} else {
			// if found ever
			// You got Middleware Group
			break
		}
	}

	return middlewares
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
