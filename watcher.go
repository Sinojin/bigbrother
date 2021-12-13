package bigbrother

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"sync"
)

var InvalidWatcher = errors.New("Please Provide watcher")
var AlreadyWatcherStarted = errors.New("Watcher already started")

type Command interface {
	ID() string // unique identifier for this command
	Command(info FileInfo)
}

type Watcher struct {
	watcher  *fsnotify.Watcher
	commands map[string]Command

	m *sync.Mutex
}

func (w *Watcher) AddCommand(cmd Command) string {
	w.m.Lock()
	defer w.m.Lock()
	if _, ok := w.commands[cmd.ID()]; !ok {
		w.commands[cmd.ID()] = cmd
	}
	return cmd.ID()
}

//AddPath Adds directories to watcher doesn't care files
func (w *Watcher) AddPath(path string) error {
	if w.watcher == nil {
		return InvalidWatcher
	}
	folders, err := FilePathWalkDir(path)
	if err != nil {
		return err
	}
	for _, folder := range folders {
		err := w.AddPath(folder)
		if err != nil {
			return err
		}
	}

	return nil
}
func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (w *Watcher) Start() error {
	if w.watcher != nil {
		return AlreadyWatcherStarted
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	w.watcher = watcher

	return nil
}

//todo : eventlerin sonucunda tetiklenen islemi karar verip kullaniciya actioni aciklamak gerek
func (w *Watcher) start() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {

			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			//todo:remove code below
			fmt.Println(err)
			//log.Println("error:", err)
		}
	}
}

func (w *Watcher) Close() {
	if w.watcher != nil {
		w.watcher.Close()
	}
}
