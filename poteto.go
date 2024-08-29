package poteto

import (
	"net/http"

	"github.com/poteto0/poteto/router"
)

type Poteto struct {
	Router *router.Router
}

func NewPoteto() *Poteto {
	return &Poteto{
		Router: &router.Router{
			RoutingTable: make(map[string]func(w http.ResponseWriter, r *http.Request)),
		},
	}
}

func (p *Poteto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handler := p.Router.RoutingTable[r.URL.Path]
		if handler == nil {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (p *Poteto) Run(addr string) {
	http.ListenAndServe(addr, p)
}
