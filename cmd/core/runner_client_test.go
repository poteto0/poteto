package core

import (
	"bufio"
	stdContext "context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/fsnotify/fsnotify"
)

func TestDefineRunnerOption(t *testing.T) {
	option := DefaultRunnerOption
	if !option.isBuildScript {
		t.Errorf("Unmatched optin")
	}
}

func TestNewRunnerClient(t *testing.T) {
	defer monkey.UnpatchAll()

	NewRunnerClient()

	monkey.Patch(fsnotify.NewWatcher, func() (*fsnotify.Watcher, error) {
		return &fsnotify.Watcher{}, errors.New("error")
	})

	NewRunnerClient()
}

func TestStartLogTransporter(t *testing.T) {
	runnerClient := NewRunnerClient()

	fileChangeStream := make(chan struct{}, 1)

	clientContext := stdContext.Background()
	logTransporter := runnerClient.LogTransporter(clientContext, fileChangeStream)

	errLogChan := make(chan error, 1)
	go func() {
		if err := logTransporter(); err != nil {
			errLogChan <- err
		}
	}()

	close(fileChangeStream)
}

func TestContextFinLogTransporter(t *testing.T) {
	runnerClient := NewRunnerClient()

	fileChangeStream := make(chan struct{}, 1)

	clientContext := stdContext.Background()
	logTransporter := runnerClient.LogTransporter(clientContext, fileChangeStream)

	errLogChan := make(chan error, 1)
	go func() {
		if err := logTransporter(); err != nil {
			errLogChan <- err
		}
	}()

	clientContext.Done()
}

func TestStreamLogTransporter(t *testing.T) {
	runnerClient := NewRunnerClient().(*runnerClient)
	reader := bufio.NewReader(strings.NewReader("test"))
	runnerClient.reader = reader

	fileChangeStream := make(chan struct{}, 1)

	clientContext := stdContext.Background()
	logTransporter := runnerClient.LogTransporter(clientContext, fileChangeStream)

	errLogChan := make(chan error, 1)
	go func() {
		if err := logTransporter(); err != nil {
			errLogChan <- err
		}
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		clientContext.Done()
	}
}

func TestStreamErrorLogTransporter(t *testing.T) {
	defer monkey.UnpatchAll()

	runnerClient := NewRunnerClient().(*runnerClient)
	reader := bufio.NewReader(strings.NewReader("test"))
	runnerClient.reader = reader
	monkey.Patch((*bufio.Reader).ReadLine, func(b *bufio.Reader) (line []byte, isPrefix bool, err error) {
		return nil, false, errors.New("error")
	})

	fileChangeStream := make(chan struct{}, 1)
	clientContext := stdContext.Background()
	logTransporter := runnerClient.LogTransporter(clientContext, fileChangeStream)

	errLogChan := make(chan error, 1)
	go func() {
		if err := logTransporter(); err != nil {
			errLogChan <- err
		}
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Unmatched not occur error")
	case <-errLogChan:
		return
	}
}

func TestClientWatcherUnexpectedEvent(t *testing.T) {
	client := NewRunnerClient().(*runnerClient)
	fileChangeStream := make(chan struct{}, 1)
	clientContext := stdContext.Background()
	fileWatcher := client.FileWatcher(clientContext, fileChangeStream)

	var err error
	go func() {
		err = fileWatcher()
	}()

	client.watcher.Events <- fsnotify.Event{}

	select {
	case <-time.After(time.Millisecond * 100):
		clientContext.Done()
		if err == nil {
			t.Errorf("Unmatched not throw expected error")
		}
	}
}

func TestClientWatcherWatchEvent(t *testing.T) {
	client := NewRunnerClient().(*runnerClient)
	fileChangeStream := make(chan struct{}, 1)
	clientContext := stdContext.Background()
	fileWatcher := client.FileWatcher(clientContext, fileChangeStream)

	go func() {
		fileWatcher()
	}()

	client.watcher.Events <- fsnotify.Event{
		Op: fsnotify.Chmod,
	}

	client.watcher.Events <- fsnotify.Event{
		Op: fsnotify.Write,
	}

	client.watcher.Events <- fsnotify.Event{
		Op: fsnotify.Create,
	}

	client.watcher.Events <- fsnotify.Event{
		Op: fsnotify.Remove,
	}

	client.watcher.Events <- fsnotify.Event{
		Op: fsnotify.Rename,
	}

	select {
	case <-time.After(time.Millisecond * 100):
		clientContext.Done()
	}
}

func TestBuildRunnerCallAsyncBuildIfFileChange(t *testing.T) {
	defer monkey.UnpatchAll()

	calledBuild := 0
	monkey.Patch((*runnerClient).AsyncBuild, func(client *runnerClient, ctx stdContext.Context, errChan chan<- error) {
		calledBuild++
	})

	client := NewRunnerClient().(*runnerClient)
	fileChangeStream := make(chan struct{}, 1)
	clientContext := stdContext.Background()
	buildRunner := client.BuildRunner(clientContext, fileChangeStream)
	go func() {
		buildRunner()
	}()
	fileChangeStream <- struct{}{}

	select {
	case <-time.After(time.Millisecond * 100):
		if calledBuild != 2 {
			t.Errorf("Unmatched call num client.AsyncBuild %v", calledBuild)
		}
	}
}

func TestBuildRunnerIfAsyncBuildError(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch((*runnerClient).AsyncBuild, func(client *runnerClient, ctx stdContext.Context, errChan chan<- error) {
		errChan <- errors.New("error")
	})

	client := NewRunnerClient().(*runnerClient)
	fileChangeStream := make(chan struct{}, 1)
	clientContext := stdContext.Background()

	var err error
	buildRunner := client.BuildRunner(clientContext, fileChangeStream)
	go func() {
		err = buildRunner()
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		if err == nil {
			t.Error("Unmatched not throw expected error")
		}
	}
}

func TestBuildRunnerCallAsyncBuildFirst(t *testing.T) {
	defer monkey.UnpatchAll()

	calledBuild := 0
	monkey.Patch((*runnerClient).AsyncBuild, func(client *runnerClient, ctx stdContext.Context, errChan chan<- error) {
		calledBuild++
	})

	client := NewRunnerClient().(*runnerClient)
	fileChangeStream := make(chan struct{}, 1)
	clientContext := stdContext.Background()

	buildRunner := client.BuildRunner(clientContext, fileChangeStream)
	go func() {
		buildRunner()
	}()

	select {
	case <-time.After(time.Millisecond * 100):
		clientContext.Done()
		if calledBuild != 1 {
			t.Error("Unmatched call num client.AsyncBuild")
		}
	}
}

func TestAsyncBuild(t *testing.T) {
	defer monkey.UnpatchAll()

	client := NewRunnerClient().(*runnerClient)

	monkey.Patch((*runnerClient).Build, func(client *runnerClient, ctx stdContext.Context) error {
		return nil
	})

	clientContext := stdContext.Background()
	errChan1 := make(chan error, 1)
	client.AsyncBuild(clientContext, errChan1)

	errChan2 := make(chan error, 1)
	monkey.Patch((*runnerClient).Build, func(client *runnerClient, ctx stdContext.Context) error {
		return errors.New("error")
	})
	client.AsyncBuild(clientContext, errChan2)

	select {
	case <-errChan1:
		t.Errorf("Unmatched occur error")
	case <-errChan2:
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Unmatched not occur error")
	}
}

func TestBuild(t *testing.T) {
	defer monkey.UnpatchAll()

	var (
		calledKill  bool
		calledStart bool
	)

	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{
			Process: &os.Process{
				Pid: 100,
			},
		}
	})
	monkey.Patch((*exec.Cmd).StdoutPipe, func(c *exec.Cmd) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("test")), nil
	})

	client := NewRunnerClient().(*runnerClient)

	tests := []struct {
		name            string
		mockKillProcess func(client *runnerClient) error
		mockStart       func(c *exec.Cmd) error
		expected        []bool
	}{
		{
			"test good case",
			func(client *runnerClient) error {
				calledKill = true
				return nil
			},
			func(c *exec.Cmd) error {
				calledStart = true
				return nil
			},
			[]bool{true, true},
		},
		{
			"test kill process cause error case",
			func(client *runnerClient) error {
				calledKill = true
				return errors.New("error")
			},
			func(c *exec.Cmd) error {
				calledStart = true
				return nil
			},
			[]bool{true, false},
		},
		{
			"test start cause error case",
			func(client *runnerClient) error {
				calledKill = true
				return nil
			},
			func(c *exec.Cmd) error {
				calledStart = true
				return errors.New("error")
			},
			[]bool{true, true},
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledKill = false
				calledStart = false
			}()

			monkey.Patch((*runnerClient).killProcess, it.mockKillProcess)
			monkey.Patch((*exec.Cmd).Start, it.mockStart)

			clientContext := stdContext.Background()
			client.Build(clientContext)

			if calledKill != it.expected[0] {
				t.Error("Unmatched call client.KillProcess")
			}

			if calledStart != it.expected[1] {
				t.Error("Unmatched call cmd.Start")
			}
		})
	}
}

func TestKillProcess(t *testing.T) {
	defer monkey.UnpatchAll()

	client := NewRunnerClient().(*runnerClient)
	if err := client.killProcess(); err != nil {
		t.Errorf("Unmatched error (not expected)")
	}

	client.pid = 100
	monkey.Patch((*runnerClient).killByOS, func(client *runnerClient) error {
		return nil
	})
	if err := client.killProcess(); err != nil {
		t.Errorf("Unmatched error (not expected)")
	}

	monkey.Patch((*runnerClient).killByOS, func(client *runnerClient) error {
		return errors.New("error")
	})
	if err := client.killProcess(); err == nil {
		t.Errorf("Unmatched error (expected)")
	}
}

// TODO: other os test
// TODO: mock GOOS
func TestKillByOS(t *testing.T) {
	defer monkey.UnpatchAll()

	var (
		calledCmd  bool
		calledKill bool
	)

	monkey.Patch((*exec.Cmd).Run, func(c *exec.Cmd) error {
		calledCmd = true
		return nil
	})

	monkey.Patch(syscall.Kill, func(pid int, sig syscall.Signal) (err error) {
		calledKill = true
		return nil
	})

	runnerClient := NewRunnerClient().(*runnerClient)

	tests := []struct {
		name         string
		runOs        string
		expectedCmd  bool
		expectedKill bool
	}{
		{"test call by linux", "linux", false, true},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			defer func() {
				calledCmd = false
				calledKill = false
			}()

			runnerClient.killByOS()

			if calledCmd != it.expectedCmd {
				t.Errorf("Unmatched call cmd")
			}

			if calledKill != it.expectedKill {
				t.Errorf("Unmatched call kill")
			}
		})
	}
}

func TestClose(t *testing.T) {
	defer monkey.UnpatchAll()

	var (
		calledKill  bool
		calledClose bool
	)

	monkey.Patch((*runnerClient).killProcess, func(client *runnerClient) error {
		calledKill = true
		return nil
	})

	monkey.Patch((*fsnotify.Watcher).Close, func(w *fsnotify.Watcher) error {
		calledClose = true
		return nil
	})

	runnerClient := NewRunnerClient().(*runnerClient)

	runnerClient.Close()

	if !calledKill {
		t.Errorf("Unmatched call client.killProcess")
	}

	if !calledClose {
		t.Errorf("Unmatched call fsnotify.watcher.Close")
	}
}
