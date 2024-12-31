package engine

import (
	stdContext "context"
	"fmt"

	"github.com/poteto0/poteto/cmd/core"
)

type EngineRunParam struct{}

func RunRun(param EngineRunParam) error {
	runnerClient := core.NewRunnerClient()
	defer runnerClient.Close()

	clientContext := stdContext.Background()
	fileChangeStream := make(chan struct{}, 1)

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

	errLogChan := make(chan error, 1)
	logTransporter := runnerClient.LogTransporter(clientContext)
	go func() {
		// log transport
		if err := logTransporter(); err != nil {
			errLogChan <- err
		}
	}()

	for {
		select {
		case err := <-errBuildChan:
			fmt.Println(err)
			return err
		case err := <-errWatcherChan:
			fmt.Println(err)
			return err
		case err := <-errLogChan:
			fmt.Println(err)
			return err
		}
	}
}
