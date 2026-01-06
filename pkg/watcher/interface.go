package watcher

// LogEntry represents a generic log entry
type LogEntry struct {
	Message     string
	ServiceName string
}

// LogWatcher interface for different OS log providers
type LogWatcher interface {
	// Start begins watching logs and sends them to the returned channel
	Start() (<-chan LogEntry, <-chan error, error)
	// Stop stops the watcher
	Stop() error
}
