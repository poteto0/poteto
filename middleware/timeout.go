package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/poteto0/poteto"
)

type TimeoutConfig struct {
	Limit           time.Duration `yaml:"limit"`
	TimeoutResponse any
}

type TimeoutResponseEx struct {
	Message string `json:"message"`
}

var DefaultTimeoutResponse = TimeoutResponseEx{
	Message: "Gateway Time Out",
}

var DefaultTimeoutConfig = TimeoutConfig{
	Limit:           time.Second * 10,
	TimeoutResponse: DefaultTimeoutResponse,
}

func TimeoutWithConfig(config TimeoutConfig) poteto.MiddlewareFunc {
	if config.Limit == 0 {
		config.Limit = DefaultTimeoutConfig.Limit
	}

	if config.TimeoutResponse == nil {
		config.TimeoutResponse = DefaultTimeoutConfig.TimeoutResponse
	}

	return func(next poteto.HandlerFunc) poteto.HandlerFunc {
		return func(ctx poteto.Context) error {
			var result error

			done := make(chan struct{})
			go func() {
				defer func() {
					// in case of panic
					if r := recover(); r != nil {
						result = fmt.Errorf("%v", r)
					}

					close(done)
				}()

				// do
				result = next(ctx)
			}()

			select {
			case <-done:
				return result
			// this loaded
			case <-time.After(config.Limit):
				// escape double response
				if ctx.GetResponse().IsCommitted {
					return nil
				}
				return ctx.JSON(http.StatusGatewayTimeout, config.TimeoutResponse)
			}
		}
	}
}
