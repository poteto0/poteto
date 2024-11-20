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

type MiddlewareTree interface {
	SearchMiddlewares(pattern string) []MiddlewareFunc
	Insert(pattern string, middlewares ...MiddlewareFunc) *middlewareTree
	Register(middlewares ...MiddlewareFunc)
}

type middlewareTree struct {
	children    map[string]MiddlewareTree
	middlewares []MiddlewareFunc
	key         string
}

func NewMiddlewareTree() MiddlewareTree {
	return &middlewareTree{
		children: make(map[string]MiddlewareTree),
	}
}

func (mt *middlewareTree) SearchMiddlewares(pattern string) []MiddlewareFunc {
	middlewares := []MiddlewareFunc{}
	currentNode := mt
	middlewares = append(middlewares, mt.middlewares...)
	patterns := strings.Split(pattern, "/")

	for _, p := range patterns {
		if p == "" {
			continue
		}

		if nextNode, ok := currentNode.children[p]; ok {
			currentNode = nextNode.(*middlewareTree)
			middlewares = append(middlewares, currentNode.middlewares...)
		} else {
			// if found ever
			// You got Middleware Tree
			break
		}
	}

	return middlewares
}

func (mt *middlewareTree) Insert(pattern string, middlewares ...MiddlewareFunc) *middlewareTree {
	currentNode := mt
	patterns := strings.Split(pattern, "/")

	for _, p := range patterns {
		if p == "" {
			continue
		}

		if _, ok := currentNode.children[p]; !ok {
			currentNode.children[p] = &middlewareTree{
				children:    make(map[string]MiddlewareTree),
				middlewares: []MiddlewareFunc{},
				key:         p,
			}
		}
		currentNode = currentNode.children[p].(*middlewareTree)
	}
	currentNode.Register(middlewares...)
	return currentNode
}

func (mt *middlewareTree) Register(middlewares ...MiddlewareFunc) {
	mt.middlewares = append(mt.middlewares, middlewares...)
}
