package watcher

import (
	"bufio"
	"encoding/json"
	"os/exec"
)

type JournalEntry struct {
	Message    string `json:"MESSAGE"`
	SystemUnit string `json:"_SYSTEMD_UNIT"`
}

type LinuxJournalWatcher struct {
	cmd *exec.Cmd
}

func NewLinuxJournalWatcher() *LinuxJournalWatcher {
	return &LinuxJournalWatcher{}
}

func (w *LinuxJournalWatcher) Start() (<-chan LogEntry, <-chan error, error) {
	logs := make(chan LogEntry)
	errs := make(chan error)

	// -f: follow
	// -p 3: priority 3 (error) and above
	// -o json: output json
	// -n 0: start from now
	w.cmd = exec.Command("journalctl", "-f", "-p", "3", "-o", "json", "-n", "0")
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

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			var entry JournalEntry
			if err := json.Unmarshal(line, &entry); err == nil {
				if entry.Message != "" {
					logs <- LogEntry{
						Message:     entry.Message,
						ServiceName: entry.SystemUnit,
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			errs <- err
		}
		w.cmd.Wait()
	}()

	return logs, errs, nil
}

func (w *LinuxJournalWatcher) Stop() error {
	if w.cmd != nil && w.cmd.Process != nil {
		return w.cmd.Process.Kill()
	}
	return nil
}
