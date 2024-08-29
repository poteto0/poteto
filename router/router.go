package router

import (
	"errors"
	"net/http"
)

type Router struct {
	RoutingTable map[string]func(w http.ResponseWriter, h *http.Request)
}

func (r *Router) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) error {
	if r.RoutingTable[path] != nil {
		return errors.New("the route " + path + " is already existed!!")
	}

	r.RoutingTable[path] = handler
	return nil
}
