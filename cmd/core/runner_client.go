package core

import (
	stdContext "context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/poteto0/poteto/utils"
)

type RunnerOption struct {
	isBuildScript bool   `yaml:"is_build_script"`
	buildScript   string `yaml:"build_script"`
}

var DefaultRunnerOption = RunnerOption{
	isBuildScript: true,
	buildScript:   "go run main.go",
}

type runnerClient struct {
	runnerDir    string
	watcher      *fsnotify.Watcher
	startupMutex sync.RWMutex
	process      *os.Process
	option       RunnerOption
}

type IRunnerClient interface {
	FileWatcher(ctx stdContext.Context, fileChangeStream chan<- struct{}) func() error
	BuildRunner(ctx stdContext.Context, fileChangeStream chan struct{}) func() error
	AsyncBuild(ctx stdContext.Context, errChan chan<- error)
	Build(ctx stdContext.Context) error
	killProcess() error
	Close() error
}

func NewRunnerClient() IRunnerClient {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	wd, _ := os.Getwd()
	watcher.Add(wd) // TODO: recursive

	return &runnerClient{
		runnerDir: wd,
		watcher:   watcher,
		option:    DefaultRunnerOption,
	}
}

func (client *runnerClient) FileWatcher(ctx stdContext.Context, fileChangeStream chan<- struct{}) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			// ファイル変更
			case event, ok := <-client.watcher.Events:
				if !ok { // event無し
					return nil
				}

				utils.PotetoPrint(
					fmt.Sprintf("poteto-cli detect event: %s", event.Op),
				)

				switch {
				// reload event
				// write, create, remove, rename
				case event.Has(fsnotify.Write),
					event.Has(fsnotify.Create),
					event.Has(fsnotify.Remove),
					event.Has(fsnotify.Rename):

					fileChangeStream <- struct{}{}

				// skip just chmod
				case event.Has(fsnotify.Chmod):
					continue

				default:
					return errors.New("unsupported event")
				}

			case err, ok := <-client.watcher.Errors:
				if !ok { // event無し
					return nil
				}
				return err
			}
		}
	}
}

func (client *runnerClient) BuildRunner(ctx stdContext.Context, fileChangeStream chan struct{}) func() error {
	return func() error {

		errChan := make(chan error, 1)
		go func() {
			client.AsyncBuild(ctx, errChan)
		}()

		for {
			select {
			// error occur in run
			case err := <-errChan:
				return err

			case <-ctx.Done():
				return ctx.Err()

			// rebuild
			case <-fileChangeStream:
				go func() {
					client.AsyncBuild(ctx, errChan)
				}()
			}
		}
	}
}

func (client *runnerClient) AsyncBuild(ctx stdContext.Context, errChan chan<- error) {
	if err := client.Build(ctx); err != nil {
		errChan <- err
	}
}

func (client *runnerClient) Build(ctx stdContext.Context) error {
	client.startupMutex.Lock()

	if err := client.killProcess(); err != nil {
		return err
	}

	// run build script
	cmd := exec.Command(client.option.buildScript)
	if err := cmd.Start(); err != nil {
		return err
	}

	// save process for kill
	client.process = cmd.Process
	client.startupMutex.Unlock()

	return nil
}

// syscall.Kill is not defined in Windows
// https://pkg.go.dev/syscall
func (client *runnerClient) killProcess() error {
	if client.process == nil {
		return nil
	}

	if err := client.process.Kill(); err != nil {
		return err
	}

	return nil
}

func (client *runnerClient) Close() error {
	return client.watcher.Close()
}
