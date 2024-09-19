package middleware

import (
	"net/http"

	"github.com/poteto0/poteto"
	"github.com/poteto0/poteto/constant"
)

type CORSConfig struct {
	AllowOrigins []string `yaml:"allow_origins"`
	AllowMethods []string `yaml:"allow_methods"`
}

var DefaultCORSConfig = CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
}

func CORSWithConfig(config CORSConfig) poteto.MiddlewareFunc {
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}

	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	allowOriginPatterns := []string{}
	for _, origin := range config.AllowOrigins {
		pattern := wrapRegExp(origin)
		allowOriginPatterns = append(allowOriginPatterns, pattern)
	}

	return func(next poteto.HandlerFunc) poteto.HandlerFunc {
		return func(ctx poteto.Context) error {
			req := ctx.GetRequest()
			res := ctx.GetResponse()
			origin := req.Header.Get(constant.HEADER_ORIGIN)

			res.Header().Add(constant.HEADER_VARY, constant.HEADER_ORIGIN)
			preflight := req.Method == http.MethodOptions

			// Not From Browser
			if origin == "" {
				if !preflight {
					return next(ctx)
				}
				return ctx.NoContent()
			}

			allowOrigin := getAllowOrigin(origin, allowOriginPatterns)

			// Origin not allowed
			if allowOrigin == "" {
				if !preflight {
					return next(ctx)
				}
				return ctx.NoContent()
			}

			// allowed method
			if matchMethod(req.Method, config.AllowMethods) {
				return next(ctx)
			}

			return ctx.NoContent()
		}
	}
}

func getAllowOrigin(origin string, allowOrigins []string) string {
	for _, o := range allowOrigins {
		if o == "*" || o == origin {
			return origin
		}
		if matchSubdomain(origin, o) {
			return origin
		}
	}

	return ""
}

func matchMethod(method string, allowMethods []string) bool {
	for _, m := range allowMethods {
		if m == method {
			return true
		}
	}

	return false
}
