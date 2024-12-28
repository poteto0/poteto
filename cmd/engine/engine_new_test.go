package engine

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"testing"

	"bou.ke/monkey"
	"github.com/poteto0/poteto/cmd/template"
)

func TestRunNew(t *testing.T) {
	defer monkey.UnpatchAll()

	param := EngineNewParam{
		ProjectName: "test-project",
	}

	var (
		calledChdir       bool
		calledMkdir       bool
		calledIsExist     bool
		calledExecCommand bool
		calledCreateMain  bool
	)

	tests := []struct {
		name            string
		mockChdir       func(dir string) error
		mockMkdir       func(name string, param fs.FileMode) error
		mockIsExist     func(err error) bool
		mockExecCommand func(name string, args ...string) *exec.Cmd
		mockExecRun     func(*exec.Cmd) error
		mockCreateMain  func(EngineNewParam) error
		expectCalled    []bool
	}{
		{
			"Test good case",
			func(dir string) error {
				calledChdir = true
				return nil
			},
			func(name string, param fs.FileMode) error {
				calledMkdir = true
				return nil
			},
			func(err error) bool {
				calledIsExist = true
				return true
			},
			func(name string, args ...string) *exec.Cmd {
				calledExecCommand = true
				return nil
			},
			func(*exec.Cmd) error {
				return nil
			},
			func(param EngineNewParam) error {
				calledCreateMain = true
				return nil
			},
			[]bool{true, true, false, true, true},
		},
		{
			"Test fail make dir",
			func(dir string) error {
				calledChdir = true
				return nil
			},
			func(name string, param fs.FileMode) error {
				calledMkdir = true
				return errors.New("error")
			},
			func(err error) bool {
				calledIsExist = true
				return false
			},
			func(name string, args ...string) *exec.Cmd {
				calledExecCommand = true
				return nil
			},
			func(*exec.Cmd) error {
				return nil
			},
			func(param EngineNewParam) error {
				calledCreateMain = true
				return nil
			},
			[]bool{true, true, true, false, false},
		},
		{
			"Test chdir fail",
			func(dir string) error {
				calledChdir = true
				return errors.New("error")
			},
			func(name string, param fs.FileMode) error {
				calledMkdir = true
				return nil
			},
			func(err error) bool {
				calledIsExist = true
				return true
			},
			func(name string, args ...string) *exec.Cmd {
				calledExecCommand = true
				return nil
			},
			func(*exec.Cmd) error {
				return nil
			},
			func(param EngineNewParam) error {
				calledCreateMain = true
				return nil
			},
			[]bool{true, true, false, false, false},
		},
		{
			"Test exec run fail",
			func(dir string) error {
				calledChdir = true
				return nil
			},
			func(name string, param fs.FileMode) error {
				calledMkdir = true
				return nil
			},
			func(err error) bool {
				calledIsExist = true
				return true
			},
			func(name string, args ...string) *exec.Cmd {
				calledExecCommand = true
				return nil
			},
			func(*exec.Cmd) error {
				return errors.New("error")
			},
			func(param EngineNewParam) error {
				calledCreateMain = true
				return nil
			},
			[]bool{true, true, false, true, false},
		},
		{
			"Test exec createMain fail",
			func(dir string) error {
				calledChdir = true
				return nil
			},
			func(name string, param fs.FileMode) error {
				calledMkdir = true
				return nil
			},
			func(err error) bool {
				calledIsExist = true
				return true
			},
			func(name string, args ...string) *exec.Cmd {
				calledExecCommand = true
				return nil
			},
			func(*exec.Cmd) error {
				return nil
			},
			func(param EngineNewParam) error {
				calledCreateMain = true
				return errors.New("error")
			},
			[]bool{true, true, false, true, true},
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledChdir = false
				calledMkdir = false
				calledIsExist = false
				calledExecCommand = false
				calledCreateMain = false
			}()

			// Mock
			monkey.Patch(os.Chdir, it.mockChdir)
			monkey.Patch(os.Mkdir, it.mockMkdir)
			monkey.Patch(os.IsExist, it.mockIsExist)
			monkey.Patch(exec.Command, it.mockExecCommand)
			monkey.Patch((*exec.Cmd).Run, it.mockExecRun)
			monkey.Patch(createMain, it.mockCreateMain)

			// Act
			RunNew(param)

			// Assert
			if calledChdir != it.expectCalled[0] {
				t.Errorf("Unmatched call for os.Chdir")
			}
			if calledMkdir != it.expectCalled[1] {
				t.Errorf("Unmatched call for os.Mkdir")
			}
			if calledIsExist != it.expectCalled[2] {
				t.Errorf("Unmatched call for os.IsExit")
			}
			if calledExecCommand != it.expectCalled[3] {
				t.Errorf("Unmatched call for exec.Command")
			}
			if calledCreateMain != it.expectCalled[4] {
				t.Errorf("Unmatched call for createMain")
			}
		})
	}
}

func TestCreateTemplateFile(t *testing.T) {
	defer monkey.UnpatchAll()

	selected := ""

	monkey.Patch(createAndWrite, func(filename, temp string) error {
		selected = temp
		return nil
	})

	monkey.Patch(choiceTemplateFile, func(param EngineNewParam) string {
		return template.DefaultTemplate
	})

	createDockerCompose()
	if selected != template.DockerComposeTemplate {
		t.Error("Unmatched docker-compose")
	}

	createDockerfile()
	if selected != template.DockerTemplate {
		t.Error("Unmatched Dockerfile")
	}

	createMain(EngineNewParam{})
	if selected != template.DefaultTemplate {
		t.Error("Unmatched main.go")
	}
}

func TestCreateAndWrite(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch((*os.File).Close, func(f *os.File) error {
		return nil
	})

	tests := []struct {
		name          string
		mockCreate    func(name string) (*os.File, error)
		mockWrite     func(f *os.File, b []byte) (n int, err error)
		expectedError bool
	}{
		{
			"test cannot create file throw error",
			func(name string) (*os.File, error) {
				return &os.File{}, errors.New("error")
			},
			func(f *os.File, b []byte) (n int, err error) {
				return 1, nil
			},
			true,
		},
		{
			"test cannot write file throw error",
			func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (n int, err error) {
				return 1, errors.New("error")
			},
			true,
		},
		{
			"test work case",
			func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (n int, err error) {
				return 1, nil
			},
			false,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			monkey.Patch(os.Create, it.mockCreate)
			monkey.Patch((*os.File).Write, it.mockWrite)

			result := createAndWrite("test", "template")
			if it.expectedError && result == nil {
				t.Error("Unmatched")
			}
			if !it.expectedError && result != nil {
				t.Error("Unmatched")
			}
		})
	}
}

func TestChoiceTemplateFile(t *testing.T) {
	tests := []struct {
		name     string
		param    EngineNewParam
		expected string
	}{
		{
			"Test choice fast api template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      true,
				IsJSONRPC:   false,
			},
			template.FastTemplate,
		},
		{
			"Test choice fast jsonrpc template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      true,
				IsJSONRPC:   true,
			},
			template.JSONRPCFastTemplate,
		},
		{
			"Test choice normal jsonrpc template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      false,
				IsJSONRPC:   true,
			},
			template.JSONRPCTemplate,
		},
		{
			"Test choice normal api template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      false,
				IsJSONRPC:   false,
			},
			template.DefaultTemplate,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			result := choiceTemplateFile(it.param)
			if result != it.expected {
				t.Error("Unmatched")
			}
		})
	}
}
