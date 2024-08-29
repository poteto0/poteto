package poteto

type Router struct {
	Routes map[string]*Route
}

func (r *Router) add(method, path string, handler func(ctx *Context)) error {
	router := r.Routes[method]

	if that_route := router.Search(path); that_route != nil {
		panic("[" + method + "] " + path + " is already used")
	}

	router.Insert(method, path, handler)
	return nil
}

func (r *Router) GET(path string, handler func(ctx *Context)) error {
	return r.add("GET", path, handler)
}

func (r *Router) POST(path string, handler func(ctx *Context)) error {
	return r.add("POST", path, handler)
}

func (r *Router) PUT(path string, handler func(ctx *Context)) error {
	return r.add("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler func(ctx *Context)) error {
	return r.add("DELETE", path, handler)
}
