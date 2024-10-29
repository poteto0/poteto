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
	Insert(pattern string, middlewares ...MiddlewareFunc)
}

type middlewareGroup struct {
	children    map[string]MiddlewareGroup
	middlewares []MiddlewareFunc
	key         string
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

func (mg *middlewareGroup) Insert(pattern string, middlewares ...MiddlewareFunc) {
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
	currentNode.register(middlewares...)
}

func (mg *middlewareGroup) register(middlewares ...MiddlewareFunc) {
	mg.middlewares = append(mg.middlewares, middlewares...)
}
