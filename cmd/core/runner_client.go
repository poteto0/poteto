package core

import (
	stdContext "context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/poteto0/poteto/utils"
)

type runnerClient struct {
	watcher *fsnotify.Watcher
}

type IRunnerClient interface {
	FileWatcher(ctx stdContext.Context, fileChangeStream chan<- struct{}) func() error
	Close() error
}

func NewRunnerClient() IRunnerClient {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	wd, _ := os.Getwd()
	watcher.Add(wd)

	return &runnerClient{
		watcher: watcher,
	}
}

func (client *runnerClient) FileWatcher(ctx stdContext.Context, fileChangeStream chan<- struct{}) func() error {
	return func() error {
		for {
			fmt.Println("go")
			select {
			case <-ctx.Done():
				return ctx.Err()

			// ファイル変更
			case event, ok := <-client.watcher.Events:
				if !ok { // event無し
					return nil
				}

				utils.PotetoPrint(
					fmt.Sprintf("poteto-cli detect event: %v", event),
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

				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
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

func (client *runnerClient) Close() error {
	return client.watcher.Close()
}
