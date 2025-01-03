package cmdrun

import (
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/poteto-go/poteto/cmd/core"
	"github.com/poteto-go/poteto/cmd/engine"
)

func TestHelp(t *testing.T) {
	help()
}

func TestCommandRun(t *testing.T) {
	defer monkey.UnpatchAll()

	calledExit := false
	monkey.Patch(os.Exit, func(code int) {
		calledExit = true
	})
	monkey.Patch(engine.RunRun, func(param core.RunnerOption) error {
		return nil
	})

	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			"test -h case",
			[]string{"poteto-cli", "run", "-h"},
			true,
		},
		{
			"test --help case",
			[]string{"poteto-cli", "run", "-h"},
			true,
		},
		{
			"test unknown case",
			[]string{"poteto-cli", "run", "unknown"},
			true,
		},
		{
			"test nothing case (work)",
			[]string{"poteto-cli", "run"},
			false,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledExit = false
			}()

			os.Args = it.args
			CommandRun()

			if calledExit != it.expected {
				t.Errorf("Unmatched call os.Exit")
			}
		})
	}
}
