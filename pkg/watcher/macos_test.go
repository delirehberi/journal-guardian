package watcher

import (
	"bufio"
	"strings"
	"testing"
)

func TestMacOSLogParsing(t *testing.T) {
	input := `[
{
  "eventMessage" : "Something bad happened",
  "processImagePath" : "/usr/bin/bad_service",
  "messageType" : "Error"
},
{
  "eventMessage" : "Another error",
  "processImagePath" : "/Applications/MyApp.app/Contents/MacOS/MyApp",
  "messageType" : "Error"
}
]
`
	w := NewMacOSLogWatcher()
	logs := make(chan LogEntry, 10)
	errs := make(chan error, 10)
	
	reader := bufio.NewReader(strings.NewReader(input))
	
	go func() {
		defer close(logs)
		defer close(errs)
		w.parseStream(reader, logs, errs)
	}()
	
	// Collect results
	var entries []LogEntry
	for entry := range logs {
		entries = append(entries, entry)
	}
	
	// Check for errors
	select {
	case err, ok := <-errs:
		if ok && err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	default:
	}
	
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}
	
	if entries[0].Message != "Something bad happened" {
		t.Errorf("Unexpected message: %s", entries[0].Message)
	}
	if entries[0].ServiceName != "bad_service" {
		t.Errorf("Unexpected service: %s", entries[0].ServiceName)
	}
	
	if entries[1].ServiceName != "MyApp" {
		t.Errorf("Unexpected service for app: %s", entries[1].ServiceName)
	}
}
