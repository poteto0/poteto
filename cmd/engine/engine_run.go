package engine

import (
	stdContext "context"
	"os"
	"os/signal"

	"github.com/poteto0/poteto/cmd/core"
)

func RunRun(option core.RunnerOption) error {
	runnerClient := core.NewRunnerClient(option)
	defer runnerClient.Close()

	// Ctrl+Cで子プロセスをkillする
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	clientContext := stdContext.Background()
	fileChangeStream := make(chan struct{}, 1)
	defer close(fileChangeStream)

	errWatcherChan := make(chan error, 1)
	fileWatcher := runnerClient.FileWatcher(
		clientContext,
		fileChangeStream,
	)
	go func() {
		// FileWatcher watch file system
		if err := fileWatcher(); err != nil {
			errWatcherChan <- err
		}
	}()

	errBuildChan := make(chan error, 1)
	buildRunner := runnerClient.BuildRunner(
		clientContext,
		fileChangeStream,
	)
	go func() {
		// Build Runner
		if err := buildRunner(); err != nil {
			errBuildChan <- err
		}
	}()

	logChan := make(chan struct{}, 1)
	errLogChan := make(chan error, 1)
	logTransporter := runnerClient.LogTransporter(clientContext, fileChangeStream)
	go func() {
		defer close(logChan)
		// log transport
		if err := logTransporter(); err != nil {
			errLogChan <- err
		}
	}()

	for {
		select {
		case err := <-errBuildChan:
			return err
		case err := <-errWatcherChan:
			return err
		case err := <-errLogChan:
			return err
		case <-logChan:
			logChan = make(chan struct{}, 1)
			// re-watch log stream watcher
			go func() {
				defer close(logChan)
				// log transport
				if err := logTransporter(); err != nil {
					errLogChan <- err
				}
			}()
		case <-quit:
			return runnerClient.Close()
		}
	}
}
