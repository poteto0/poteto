package main

import (
	"os"
	"testing"

	"bou.ke/monkey"
)

func TestPotetoCliMain(t *testing.T) {
	t.Errorf("Error")
	var isExit bool
	monkey.Patch(os.Exit, func(code int) {
		isExit = true
		os.Args = []string{"poteto-cli", "escape"}
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
