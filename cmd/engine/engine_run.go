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
	go func() {
		// FileWatcher watch file system
		err := fileWatcher()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}()

	for {
		select {
		// TODO: Rebuild
		case <-fileChangeStream:
			return nil
		}
	}
}
