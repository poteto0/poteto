package engine

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"testing"

	"bou.ke/monkey"
)

func TestRunCommandNew(t *testing.T) {
	var calledChdir bool
	var calledMkdir bool
	var calledIsExist bool
	var calledExecCommand bool
	var calledCreateMain bool

	tests := []struct {
		name            string
		mockChdir       func(dir string) error
		mockMkdir       func(name string, param fs.FileMode) error
		mockIsExist     func(err error) bool
		mockExecCommand func(name string, args ...string) *exec.Cmd
		mockExecRun     func(*exec.Cmd) error
		mockCreateMain  func() error
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
			func() error {
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
			func() error {
				calledCreateMain = true
				return nil
			},
			[]bool{true, true, true, false, false},
		},
		{
			"Test chdir fail",
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
				return true
			},
			func(name string, args ...string) *exec.Cmd {
				calledExecCommand = true
				return nil
			},
			func(*exec.Cmd) error {
				return nil
			},
			func() error {
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
			func() error {
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
			func() error {
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

			monkey.Patch(os.Chdir, it.mockChdir)
			monkey.Patch(os.Mkdir, it.mockMkdir)
			monkey.Patch(exec.Command, it.mockExecCommand)
			monkey.Patch((*exec.Cmd).Run, it.mockExecRun)
			monkey.Patch(createMain, it.mockCreateMain)
		})

		run("hello")

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

		if !calledCreateMain != it.expectCalled[4] {
			t.Errorf("Unmatched call for createMain")
		}
	}

}

func TestCommandNew(t *testing.T) {
	monkey.Unpatch(CommandNew)
	// Mock
	var capture int
	monkey.Patch(os.Exit, func(code int) { capture = code })
	monkey.Patch(run, func(a string) error {
		return nil
	})

	tests := []struct {
		name   string
		arg    string
		IsExit bool
	}{
		{
			"test unexpected option case",
			"-hello",
			true,
		},
		{
			"test -h case",
			"-h",
			true,
		},
		{
			"test --help case",
			"--help",
			true,
		},
		{
			"test -f case",
			"-f",
			false,
		},
		{
			"test --fast case",
			"--fast",
			false,
		},
		{
			"test -d case",
			"-d",
			false,
		},
		{
			"test --docker case",
			"--docker",
			false,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				capture = 0
			}()

			os.Args = []string{"poteto-cli", "new", it.arg}
			CommandNew()

			if it.IsExit {
				if capture != -1 {
					t.Errorf("Not exited")
				}
			} else {
				if capture != 0 {
					t.Errorf("Not go through")
				}
			}
		})
	}
}

func TestCreateMain(t *testing.T) {
	var calledCreate bool
	var calledWrite bool

	monkey.Patch((*os.File).Close, func(f *os.File) error {
		return nil
	})

	tests := []struct {
		name         string
		mockCreate   func(string) (*os.File, error)
		mockWite     func(*os.File, []byte) (int, error)
		expectCalled []bool
	}{
		{
			"Test right case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, nil
			},
			[]bool{true, true},
		},
		{
			"Test create fail case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, errors.New("error")
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, nil
			},
			[]bool{true, false},
		},
		{
			"Test write fail case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, errors.New("error")
			},
			[]bool{true, true},
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledCreate = false
				calledWrite = false
			}()

			monkey.Patch(os.Create, it.mockCreate)
			monkey.Patch((*os.File).Write, it.mockWite)

			createMain()

			if calledCreate != it.expectCalled[0] {
				t.Errorf("Unmatched call of os.Create")
			}

			if calledWrite != it.expectCalled[1] {
				t.Errorf("Unmatched call of f.Write")
			}
		})
	}
}

func TestCreateDocker(t *testing.T) {
	var calledCreate bool
	var calledWrite bool

	monkey.Patch((*os.File).Close, func(f *os.File) error {
		return nil
	})

	tests := []struct {
		name         string
		mockCreate   func(string) (*os.File, error)
		mockWite     func(*os.File, []byte) (int, error)
		expectCalled []bool
	}{
		{
			"Test right case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, nil
			},
			[]bool{true, true},
		},
		{
			"Test create fail case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, errors.New("error")
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, nil
			},
			[]bool{true, false},
		},
		{
			"Test write fail case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, errors.New("error")
			},
			[]bool{true, true},
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledCreate = false
				calledWrite = false
			}()

			monkey.Patch(os.Create, it.mockCreate)
			monkey.Patch((*os.File).Write, it.mockWite)

			createDockerfile()

			if calledCreate != it.expectCalled[0] {
				t.Errorf("Unmatched call of os.Create")
			}

			if calledWrite != it.expectCalled[1] {
				t.Errorf("Unmatched call of f.Write")
			}
		})
	}
}

func TestCreateDockerCompose(t *testing.T) {
	var calledCreate bool
	var calledWrite bool

	monkey.Patch((*os.File).Close, func(f *os.File) error {
		return nil
	})

	tests := []struct {
		name         string
		mockCreate   func(string) (*os.File, error)
		mockWite     func(*os.File, []byte) (int, error)
		expectCalled []bool
	}{
		{
			"Test right case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, nil
			},
			[]bool{true, true},
		},
		{
			"Test create fail case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, errors.New("error")
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, nil
			},
			[]bool{true, false},
		},
		{
			"Test write fail case",
			func(name string) (*os.File, error) {
				calledCreate = true
				return &os.File{}, nil
			},
			func(f *os.File, b []byte) (int, error) {
				calledWrite = true
				return 0, errors.New("error")
			},
			[]bool{true, true},
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledCreate = false
				calledWrite = false
			}()

			monkey.Patch(os.Create, it.mockCreate)
			monkey.Patch((*os.File).Write, it.mockWite)

			createDockerCompose()

			if calledCreate != it.expectCalled[0] {
				t.Errorf("Unmatched call of os.Create")
			}

			if calledWrite != it.expectCalled[1] {
				t.Errorf("Unmatched call of f.Write")
			}
		})
	}
}
