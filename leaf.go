package poteto

type leaf struct {
	poteto   Poteto
	basePath string
}

func NewLeaf(poteto Poteto, basePath string) Leaf {
	return &leaf{
		poteto:   poteto,
		basePath: basePath,
	}
}

type Leaf interface {
	Register(middlewares ...MiddlewareFunc) *middlewareTree
	GET(addPath string, handler HandlerFunc) error
	POST(addPath string, handler HandlerFunc) error
	PUT(addPath string, handler HandlerFunc) error
	DELETE(addPath string, handler HandlerFunc) error
}

// internal call Poteto.Combine w/ bp
func (l *leaf) Register(middlewares ...MiddlewareFunc) *middlewareTree {
	return l.poteto.Combine(l.basePath, middlewares...)
}

// internal call Poteto.GET w/ bp
func (l *leaf) GET(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.GET(path, handler)
}

// internal call Poteto.POST w/ bp
func (l *leaf) POST(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.POST(path, handler)
}

// internal call Poteto.PUT w/ bp
func (l *leaf) PUT(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.PUT(path, handler)
}

func (l *leaf) DELETE(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.DELETE(path, handler)
}
