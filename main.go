package main

import (
	"fmt"
	"log_watcher/pkg/config"
	"log_watcher/pkg/notifier"
	"log_watcher/pkg/providers"
	"log_watcher/pkg/watcher"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

func fanIn(watchers []watcher.LogWatcher) (<-chan watcher.LogEntry, <-chan error, func() error) {
	aggLogs := make(chan watcher.LogEntry)
	aggErrs := make(chan error)
	var wg sync.WaitGroup

	for _, w := range watchers {
		logs, errs, err := w.Start()
		if err != nil {
			fmt.Printf("Error starting specific watcher: %v\n", err)
			continue
		}

		wg.Add(1)
		go func(l <-chan watcher.LogEntry, e <-chan error) {
			defer wg.Done()
			for {
				select {
				case entry, ok := <-l:
					if !ok {
						return
					}
					aggLogs <- entry
				case err, ok := <-e:
					if !ok {
						return
					}
					aggErrs <- err
				}
			}
		}(logs, errs)
	}

	stopFunc := func() error {
		for _, w := range watchers {
			w.Stop()
		}
		// wait for all goroutines to finish
		// Note from developer: In a strictly proper fan-in we should close aggLogs/aggErrs
		// after wg.Wait(), but since we are main looping, we just stop.
		// For proper cleanup in a library, we'd handle channel closing more carefully.
		return nil
	}

	return aggLogs, aggErrs, stopFunc
}

func main() {
	// 1. Initialize LLM Provider
	provider, err := providers.NewProvider()
	if err != nil {
		fmt.Printf("Error initializing provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[*] Go Journal Watcher started using %s\n", provider.Name())

	// 2. Load Config
	configPath := config.FindConfig()
	var cfg *config.Config
	if configPath != "" {
		fmt.Printf("[*] Loading configuration from %s\n", configPath)
		c, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
		} else {
			cfg = c
		}
	}

	var watchers []watcher.LogWatcher

	if cfg != nil {
		fmt.Println("[*] Configuration found, initializing configured watchers...")
		for _, source := range cfg.Sources {
			w, err := watcher.NewWatcher(source)
			if err != nil {
				fmt.Printf("Error initializing watcher for %s: %v\n", source.Type, err)
				continue
			}
			watchers = append(watchers, w)
		}
	} else {
		fmt.Println("[*] No configuration found, defaulting to OS watcher...")
		// Default behavior
		switch runtime.GOOS {
		case "linux":
			watchers = append(watchers, watcher.NewLinuxJournalWatcher())
		case "darwin":
			watchers = append(watchers, watcher.NewMacOSLogWatcher())
		default:
			fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
			os.Exit(1)
		}
	}

	if len(watchers) == 0 {
		fmt.Println("No watchers initialized. Exiting.")
		os.Exit(1)
	}

	// 3. Start Watching (Fan-In)
	logs, errs, stopWatchers := fanIn(watchers)
	defer stopWatchers()

	// Initialize Notifier
	var notify notifier.Notifier
	switch runtime.GOOS {
	case "linux":
		notify = notifier.NewLinuxNotifier()
	case "darwin":
		notify = notifier.NewMacOSNotifier()
	default:
		// Fallback or simple print notifier
		notify = notifier.NewLinuxNotifier() // Assuming interface is compatible or valid fallback
	}


	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("[*] Watching logs from %d source(s)...\n", len(watchers))

	go func() {
		<-sigChan
		fmt.Println("\n[*] Shutting down...")
		stopWatchers()
		os.Exit(0)
	}()

	// 4. Process Logs
	for {
		select {
		case err := <-errs:
			if err != nil {
				fmt.Printf("Watcher error: %v\n", err)
			}
		case entry, ok := <-logs:
			if !ok {
				return
			}
			
			fmt.Printf("\n[!] Error in %s:\n    %s\n", entry.ServiceName, entry.Message)
			fmt.Println("    Asking AI...")

			suggestion, err := provider.Generate(fmt.Sprintf("Service: %s. Log: %s", entry.ServiceName, entry.Message))
			if err != nil {
				fmt.Printf("Error generating suggestion: %v\n", err)
				continue
			}

			fmt.Printf("\n--- 💡 FIX SUGGESTION ---\n%s\n---------------------\n", suggestion)

			if err := notify.Send(fmt.Sprintf("Error: %s", entry.ServiceName), suggestion); err != nil {
				fmt.Printf("Error sending notification: %v\n", err)
			}
		}
	}
}
