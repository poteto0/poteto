package poteto

import (
	"fmt"
	"net/http"
)

type Poteto struct {
	Router Router
}

func New() *Poteto {
	return &Poteto{
		Router: NewRouter([]string{"GET", "POST", "PUT", "DELETE"}),
	}
}

func (p *Poteto) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	routes := p.Router.GetRoutesByMethod(r.Method)

	targetRoute := routes.Search(r.URL.Path)
	handler := targetRoute.GetHandler()

	if targetRoute == nil || handler == nil {
		ctx.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.SetPath(r.URL.Path)
	handler(ctx)
}

func (p *Poteto) Run(addr string) {
	if err := http.ListenAndServe(addr, p); err != nil {
		fmt.Println("Error occurred")
	}
}
