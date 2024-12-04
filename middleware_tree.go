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
	currentNode := mt
	// faster
	middlewares := mt.middlewares
	if pattern == "/" {
		return middlewares
	}

	rightPattern := pattern[1:]
	param := ""

	for {
		id := strings.Index(rightPattern, "/")
		if id < 0 {
			param = rightPattern
		} else {
			param = rightPattern[:id]
			rightPattern = rightPattern[(id + 1):]
		}
		if nextNode, ok := currentNode.children[param]; ok {
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
	if pattern == "/" || pattern == "" {
		currentNode.Register(middlewares...)
		return currentNode
	}
	rightPattern := pattern[1:]
	param := ""

	for {
		id := strings.Index(rightPattern, "/")
		if id < 0 {
			param = rightPattern
		} else {
			param = rightPattern[:id]
			rightPattern = rightPattern[(id + 1):]
		}
		if _, ok := currentNode.children[param]; !ok {
			currentNode.children[param] = &middlewareTree{
				children:    make(map[string]MiddlewareTree),
				middlewares: []MiddlewareFunc{},
				key:         param,
			}
		}
		currentNode = currentNode.children[param].(*middlewareTree)

		if id < 0 {
			break
		}
	}
	currentNode.Register(middlewares...)
	return currentNode
}

func (mt *middlewareTree) Register(middlewares ...MiddlewareFunc) {
	mt.middlewares = append(mt.middlewares, middlewares...)
}
