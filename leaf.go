package poteto

/*
/ Leaf Make Router Great
/ p.Leaf("/users", func(leaf Leaf) {
/   leaf.Register(sampleMiddleware)
/   leaf.GET("/", getAllUserForTest)
/   leaf.POST("/create", getAllUserForTest)
/   leaf.PUT("/change", getAllUserForTest)
/   leaf.DELETE("/delete", getAllUserForTest)
/ })
*/

type leaf struct {
	poteto   Poteto
	basePath string
}

type Leaf interface {
	Register(middlewares ...MiddlewareFunc) *middlewareTree
	GET(addPath string, handler HandlerFunc) error
	POST(addPath string, handler HandlerFunc) error
	PUT(addPath string, handler HandlerFunc) error
	PATCH(path string, handler HandlerFunc) error
	DELETE(addPath string, handler HandlerFunc) error
	HEAD(path string, handler HandlerFunc) error
	OPTIONS(path string, handler HandlerFunc) error
	TRACE(path string, handler HandlerFunc) error
	CONNECT(path string, handler HandlerFunc) error
}

func NewLeaf(poteto Poteto, basePath string) Leaf {
	return &leaf{
		poteto:   poteto,
		basePath: basePath,
	}
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

// internal call Poteto.PATCH w/ bp
func (l *leaf) PATCH(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.PATCH(path, handler)
}

// internal call Poteto.DELETE w/ bp
func (l *leaf) DELETE(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.DELETE(path, handler)
}

// internal call Poteto.HEAD w/ bp
func (l *leaf) HEAD(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.HEAD(path, handler)
}

// internal call Poteto.OPTIONS w/ bp
func (l *leaf) OPTIONS(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.OPTIONS(path, handler)
}

// internal call Poteto.TRACE w/ bp
func (l *leaf) TRACE(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.TRACE(path, handler)
}

// internal call Poteto.CONNECT w/ bp
func (l *leaf) CONNECT(addPath string, handler HandlerFunc) error {
	path := l.basePath + addPath
	return l.poteto.CONNECT(path, handler)
}
