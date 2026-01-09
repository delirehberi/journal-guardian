package watcher

import (
	"fmt"
	"path/filepath"

	"github.com/nxadm/tail"
)

type FileLogWatcher struct {
	path string
	tail *tail.Tail
}

func NewFileLogWatcher(path string) *FileLogWatcher {
	return &FileLogWatcher{path: path}
}

func (w *FileLogWatcher) Start() (<-chan LogEntry, <-chan error, error) {
	logs := make(chan LogEntry)
	errs := make(chan error)

	// Configure tail
	config := tail.Config{
		Follow:    true,
		ReOpen:    true, // Handle log rotation
		MustExist: false,
		Poll:      true, // Polling is often more reliable for various filesystems/docker mounts
		Logger:    tail.DiscardingLogger,
	}

	t, err := tail.TailFile(w.path, config)
	if err != nil {
		close(logs)
		close(errs)
		return nil, nil, err
	}
	w.tail = t

	go func() {
		defer close(logs)
		defer close(errs)
		// Cleanup when done (though Stop() handles main cleanup)
		// w.tail.Cleanup() 

		for line := range w.tail.Lines {
			if line.Err != nil {
				errs <- line.Err
				continue
			}
			logs <- LogEntry{
				Message:     line.Text,
				ServiceName: fmt.Sprintf("file:%s", filepath.Base(w.path)),
			}
		}
		
		// If lines channel closes, check for error (though tail typically stays open with Follow)
		if err := w.tail.Err(); err != nil {
			errs <- err
		}
	}()

	return logs, errs, nil
}

func (w *FileLogWatcher) Stop() error {
	if w.tail != nil {
		return w.tail.Stop()
	}
	return nil
}
