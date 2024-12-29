package cmdnew

import (
	"errors"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/manifoldco/promptui"
	"github.com/poteto0/poteto/cmd/engine"
)

func TestCommandNew(t *testing.T) {
	var (
		calledExit      bool
		calledGetWd     bool
		calledPromptRun bool
		calledRunNew    bool
	)

	// mock
	monkey.Patch(os.Exit, func(code int) {
		calledExit = true
	})
	monkey.Patch(os.Getwd, func() (dir string, err error) {
		calledGetWd = true
		return "test", nil
	})
	monkey.Patch(engine.RunNew, func(param engine.EngineNewParam) error {
		calledRunNew = true
		return nil
	})

	tests := []struct {
		name          string
		mockPromptRun func(p *promptui.Prompt) (string, error)
		args          []string
		expected      []bool
	}{
		{
			"test help & exit -h",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "-h"},
			[]bool{true, false, false, false},
		},
		{
			"test help & exit --help",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "-h"},
			[]bool{true, false, false, false},
		},
		{
			"test unexpected & exit",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "unexpected"},
			[]bool{true, false, false, false},
		},
		{
			"test -f case",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "-f"},
			[]bool{true, true, true, true},
		},
		{
			"test --fast case",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "--fast"},
			[]bool{true, true, true, true},
		},
		{
			"test -d case",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "-d"},
			[]bool{true, true, true, true},
		},
		{
			"test --docker case",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "--docker"},
			[]bool{true, true, true, true},
		},
		{
			"test -j case",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "-j"},
			[]bool{true, true, true, true},
		},
		{
			"test --jsonrpc case",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", nil
			},
			[]string{"poteto-cli", "new", "--jsonrpc"},
			[]bool{true, true, true, true},
		},
		{
			"test prompt.Run throw error",
			func(p *promptui.Prompt) (string, error) {
				calledPromptRun = true
				return "test", errors.New("error")
			},
			[]string{"poteto-cli", "new"},
			[]bool{true, true, true, false},
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledExit = false
				calledGetWd = false
				calledPromptRun = false
				calledRunNew = false
			}()

			os.Args = it.args

			// mock
			monkey.Patch((*promptui.Prompt).Run, it.mockPromptRun)

			CommandNew()

			// assert
			if calledExit != it.expected[0] {
				t.Error("Unmatched call called os.Exit")
			}
			if calledGetWd != it.expected[1] {
				t.Error("Unmatched call os.Getwd")
			}
			if calledPromptRun != it.expected[2] {
				t.Error("Unmatched call prompt.Run")
			}
			if calledRunNew != it.expected[3] {
				t.Error("Unmatched call engine.RunNew")
			}
		})
	}
}

// just run
func TestHelp(t *testing.T) {
	help()
}
