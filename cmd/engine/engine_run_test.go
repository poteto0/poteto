package engine

import (
	"testing"
	"time"

	"github.com/poteto-go/poteto/cmd/core"
)

func TestSuccessRunRun(t *testing.T) {
	param := core.DefaultRunnerOption

	go func() {
		RunRun(param)
	}()

	select {
	case <-time.After(time.Microsecond * 100):
		return
	}
}
