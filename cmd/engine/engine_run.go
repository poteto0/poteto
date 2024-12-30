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

	fmt.Println("hello world")

	clientContext := stdContext.Background()
	fileChangeStream := make(chan struct{}, 1)
	fileWatcher := runnerClient.FileWatcher(
		clientContext,
		fileChangeStream,
	)
	buildRunner := runnerClient.BuildRunner(
		clientContext,
		fileChangeStream,
	)

	errBuildChan := make(chan error, 1)
	errWatcherChan := make(chan error, 1)
	go func() {
		if err := buildRunner(); err != nil {
			errBuildChan <- err
		}

		// FileWatcher watch file system
		if err := fileWatcher(); err != nil {
			errWatcherChan <- err
		}
	}()

	for {
		select {
		case err := <-errBuildChan:
			return err
		case err := <-errWatcherChan:
			return err
		}
	}
}
