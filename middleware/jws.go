package middleware

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poteto0/poteto"
	"github.com/poteto0/poteto/constant"
)

type potetoJWSConfig struct {
	AuthScheme string
	SignMethod string
	SignKey    any
	ContextKey string
	ClaimsFunc func(c poteto.Context) jwt.Claims
}

type PotetoJWSConfig interface {
	KeyFunc(token *jwt.Token) (any, error)
	ParseToken(ctx poteto.Context, auth string) (any, error)
}

var DefaultJWSConfig = &potetoJWSConfig{
	AuthScheme: constant.AUTH_SCHEME,
	SignMethod: constant.ALGORITHM_HS256,
	ContextKey: "user",
	ClaimsFunc: func(c poteto.Context) jwt.Claims {
		return jwt.MapClaims{}
	},
}

func NewJWSConfig(key any, contextKey string) poteto.MiddlewareFunc {
	cfg := DefaultJWSConfig
	cfg.ContextKey = contextKey
	cfg.SignKey = key
	return JWSWithConfig(cfg)
}

func (cfg *potetoJWSConfig) KeyFunc(token *jwt.Token) (any, error) {
	if token.Method.Alg() != cfg.SignMethod {
		return nil, errors.New("unexpected jwt signing method")
	}

	if cfg.SignKey == nil {
		return nil, errors.New("undefined sign key")
	}

	return cfg.SignKey, nil
}

func (cfg *potetoJWSConfig) ParseToken(ctx poteto.Context, auth string) (any, error) {
	token, err := jwt.ParseWithClaims(auth, cfg.ClaimsFunc(ctx), cfg.KeyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token, nil
}

func JWSWithConfig(cfg PotetoJWSConfig) poteto.MiddlewareFunc {
	config := cfg.(*potetoJWSConfig)
	if config.SignKey == nil {
		panic(config.SignKey)
	}

	return func(next poteto.HandlerFunc) poteto.HandlerFunc {
		return func(ctx poteto.Context) error {
			authValue := extractBearer(ctx)

			token, err := cfg.ParseToken(ctx, authValue)
			if err != nil {
				return err
			}

			ctx.Set(config.ContextKey, token)
			return next(ctx)
		}
	}
}

func extractBearer(ctx poteto.Context) string {
	authHeader := ExtractFromHeader(ctx, constant.HEADER_AUTHORIZATION)
	target := constant.AUTH_SCHEME + constant.PARAM_PREFIX
	return strings.Split(authHeader, target)[1]
}
