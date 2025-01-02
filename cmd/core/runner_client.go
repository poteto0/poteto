package core

import (
	"bufio"
	stdContext "context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/poteto0/poteto/utils"
)

type RunnerOption struct {
	isBuildScript bool     `yaml:"is_build_script"`
	buildScript   []string `yaml:"build_script"`
}

var DefaultRunnerOption = RunnerOption{
	isBuildScript: true,
	buildScript:   []string{"go", "run", "main.go"},
}

type runnerClient struct {
	runnerDir    string
	watcher      *fsnotify.Watcher
	startupMutex sync.RWMutex
	option       RunnerOption
	logStream    io.ReadCloser
	pid          int
	reader       *bufio.Reader
}

type IRunnerClient interface {
	LogTransporter(ctx stdContext.Context, fileChangeStream chan struct{}) func() error
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

func (client *runnerClient) LogTransporter(ctx stdContext.Context, fileChangeStream chan struct{}) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			// re-watch log stream watcher
			case <-fileChangeStream:
				return nil

			// log streamed
			default:
				if client.reader == nil {
					continue
				}

				line, _, err := client.reader.ReadLine()
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return err
				}

				utils.PotetoPrint(
					fmt.Sprintf("%s %s\n", color.HiGreenString("poteto |"), string(line)),
				)
			}
		}
	}
}

func (client *runnerClient) FileWatcher(ctx stdContext.Context, fileChangeStream chan<- struct{}) func() error {
	return func() error {
		var (
			timer     *time.Timer
			lastEvent fsnotify.Event
		)
		timer = time.NewTimer(time.Millisecond)
		<-timer.C // timer should be expired at first

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			// ファイル変更
			case event, ok := <-client.watcher.Events:
				if !ok { // event無し
					return nil
				}

				lastEvent = event
				timer.Reset(time.Millisecond)

			// 複数回イベントが発行されるため、timerを上で作り出して、一定時間後に処理する
			case <-timer.C:
				utils.PotetoPrint(
					fmt.Sprintf("%s poteto-cli detect event: %s\n", color.HiBlueString("pdebug |"), lastEvent.Op),
				)

				switch {
				// reload event
				// write, create, remove, rename
				case lastEvent.Has(fsnotify.Write),
					lastEvent.Has(fsnotify.Create),
					lastEvent.Has(fsnotify.Remove),
					lastEvent.Has(fsnotify.Rename):

					fileChangeStream <- struct{}{}

				// skip just chmod
				case lastEvent.Has(fsnotify.Chmod):
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
		fmt.Println(err)
		client.startupMutex.Unlock()
		return err
	}

	// run build script
	cmd := exec.Command("go", "run", "main.go")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	client.logStream, _ = cmd.StdoutPipe()
	client.reader = bufio.NewReader(client.logStream)

	// async start
	if err := cmd.Start(); err != nil {
		client.startupMutex.Unlock()
		return err
	}

	// save process for kill
	client.pid = cmd.Process.Pid
	client.startupMutex.Unlock()

	return nil
}

func (client *runnerClient) killProcess() error {
	if client.pid == 0 {
		return nil
	}

	if err := client.killByOS(); err != nil {
		return err
	}
	return nil
}

func (client *runnerClient) killByOS() error {
	switch runtime.GOOS {
	case "windows":
		// syscall.Kill is not defined in Windows
		// https://pkg.go.dev/syscall
		cmd := exec.Command("taskkill", "/pid", fmt.Sprintf("%d %s", client.pid, "/F"))
		return cmd.Run()

	case "linux", "ubuntu":
		// -pid
		// https://makiuchi-d.github.io/2020/05/10/go-kill-child-process.ja.html
		if err := syscall.Kill(-client.pid, syscall.SIGTERM); err != nil {
			return err
		}
		return nil

	default:
		if err := syscall.Kill(-client.pid, syscall.SIGTERM); err != nil {
			return err
		}
		return nil
	}
}

func (client *runnerClient) Close() error {
	client.killProcess()
	return client.watcher.Close()
}
