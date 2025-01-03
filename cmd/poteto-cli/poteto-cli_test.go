package main

import (
	"os"
	"testing"

	"bou.ke/monkey"
	cmdnew "github.com/poteto0/poteto/cmd/cmd-new"
	cmdrun "github.com/poteto0/poteto/cmd/cmd-run"
)

func TestPotetoCliMai(t *testing.T) {
	defer monkey.UnpatchAll()

	var isExit bool
	monkey.Patch(os.Exit, func(code int) {
		isExit = true
		os.Args = []string{"poteto-cli", "escape"}
	})
	monkey.Patch(cmdnew.CommandNew, func() {
		return
	})
	monkey.Patch(cmdrun.CommandRun, func() {
		return
	})

	tests := []struct {
		name   string
		arg    string
		isExit bool
	}{
		{
			"Test arg len 0 case",
			"",
			true,
		},
		{
			"Test arg -h case",
			"-h",
			true,
		},
		{
			"Test arg --help case",
			"--help",
			true,
		},
		{
			"Test arg unknown case",
			"-hello",
			true,
		},
		{
			"Test arg new case",
			"new",
			true,
		},
		{
			"Test arg run case",
			"run",
			true,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				isExit = false
			}()

			os.Args = []string{"poteto-cli", it.arg}
			if len(it.arg) == 0 {
				os.Args = []string{"poteto-cli"}
			}

			main()

			if isExit != it.isExit {
				t.Errorf("Unmatched")
			}
		})
	}
}
