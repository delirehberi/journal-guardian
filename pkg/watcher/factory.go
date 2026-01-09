package watcher

import (
	"fmt"
	"log_watcher/pkg/config"
	"runtime"
)

func NewWatcher(cfg config.LogSourceConfig) (LogWatcher, error) {
	switch cfg.Type {
	case "file":
		path, ok := cfg.Params["path"]
		if !ok || path == "" {
			return nil, fmt.Errorf("missing 'path' parameter for file watcher")
		}
		return NewFileLogWatcher(path), nil
	case "journalctl":
		if runtime.GOOS != "linux" {
			return nil, fmt.Errorf("journalctl watcher is only supported on Linux")
		}
		return NewLinuxJournalWatcher(), nil
	case "macos":
		if runtime.GOOS != "darwin" {
			return nil, fmt.Errorf("macos watcher is only supported on macOS")
		}
		return NewMacOSLogWatcher(), nil
	default:
		return nil, fmt.Errorf("unknown watcher type: %s", cfg.Type)
	}
}
