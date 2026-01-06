package watcher

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"path/filepath"
)

// MacOSLogEntry represents the JSON structure from 'log stream --style json'
type MacOSLogEntry struct {
	EventMessage     string `json:"eventMessage"`
	ProcessImagePath string `json:"processImagePath"`
	MessageType      string `json:"messageType"`
}

type MacOSLogWatcher struct {
	cmd *exec.Cmd
}

func NewMacOSLogWatcher() *MacOSLogWatcher {
	return &MacOSLogWatcher{}
}

func (w *MacOSLogWatcher) Start() (<-chan LogEntry, <-chan error, error) {
	logs := make(chan LogEntry)
	errs := make(chan error)

	// log stream --style json --predicate 'type == "error"'
	w.cmd = exec.Command("log", "stream", "--style", "json", "--predicate", "type == \"error\"")
	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := w.cmd.Start(); err != nil {
		return nil, nil, err
	}

	go func() {
		defer close(logs)
		defer close(errs)

		w.parseStream(bufio.NewReader(stdout), logs, errs)

		w.cmd.Wait()
	}()

	return logs, errs, nil
}

func (w *MacOSLogWatcher) parseStream(reader *bufio.Reader, logs chan<- LogEntry, errs chan<- error) {
	decoder := json.NewDecoder(reader)

	// Consume the opening '['
	// note: log stream might output some text preamble, but usually that's stderr.
	// We assume strict JSON array start here.
	t, err := decoder.Token()
	if err != nil {
		errs <- err
		return
	}
	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		// If it's not '[', it might be headers or empty.
		// For now we treat it as error or just return.
		return
	}

	for decoder.More() {
		var entry MacOSLogEntry
		if err := decoder.Decode(&entry); err != nil {
			errs <- err
			return
		}

		if entry.EventMessage != "" {
			serviceName := filepath.Base(entry.ProcessImagePath)
			if serviceName == "." || serviceName == "/" {
				serviceName = "unknown"
			}

			logs <- LogEntry{
				Message:     entry.EventMessage,
				ServiceName: serviceName,
			}
		}
	}
	
	// Consume closing ']'
	_, _ = decoder.Token()
}

func (w *MacOSLogWatcher) Stop() error {
	if w.cmd != nil && w.cmd.Process != nil {
		return w.cmd.Process.Kill()
	}
	return nil
}
