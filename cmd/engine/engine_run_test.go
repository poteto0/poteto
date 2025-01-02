package engine

import (
	"testing"
	"time"
)

func TestSuccessRunRun(t *testing.T) {
	param := EngineRunParam{}

	go func() {
		RunRun(param)
	}()

	select {
	case <-time.After(time.Microsecond * 100):
		return
	}
}
