package middleware

import (
	"net/http"
	"regexp"

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

			allowSubDomain := getAllowSubDomain(origin, config.AllowOrigins)
			// allowed origin path
			allowOrigin := getAllowOrigin(allowSubDomain, allowOriginPatterns)

			// Origin not allowed
			if allowOrigin == "" {
				if !preflight {
					return next(ctx)
				}
				return ctx.NoContent()
			}

			// allowed method
			if matchMethod(req.Method, config.AllowMethods) {
				res.Header().Set(constant.HEADER_ACCESS_CONTROL_ORIGIN, allowOrigin)
				return next(ctx)
			}

			return ctx.NoContent()
		}
	}
}

func getAllowSubDomain(origin string, allowOrigins []string) string {
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

func getAllowOrigin(origin string, allowOriginPatterns []string) string {
	for _, pattern := range allowOriginPatterns {
		if match, _ := regexp.MatchString(pattern, origin); match {
			return origin
		}
	}
	return ""
}
