package poteto

import (
	"fmt"
	"net/http"
)

type Poteto struct {
	Router *Router
}

func New() *Poteto {
	return &Poteto{
		Router: &Router{
			Routes: map[string]*Route{
				"GET":    NewRoute(),
				"POST":   NewRoute(),
				"PUT":    NewRoute(),
				"DELETE": NewRoute(),
			},
		},
	}
}

func (p *Poteto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	routes := p.Router.Routes[r.Method]

	targetRoute := routes.Search(r.URL.Path)

	if targetRoute == nil || targetRoute.handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.SetPath(r.URL.Path)
	targetRoute.handler(ctx)
}

func (p *Poteto) Run(addr string) {
	if err := http.ListenAndServe(addr, p); err != nil {
		fmt.Println("Error occured")
	}
}
