package main

import (
	"fmt"
	"log_watcher/pkg/notifier"
	"log_watcher/pkg/providers"
	"log_watcher/pkg/watcher"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	// 1. Initialize LLM Provider
	provider, err := providers.NewProvider()
	if err != nil {
		fmt.Printf("Error initializing provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[*] Go Journal Watcher started using %s\n", provider.Name())

	// 2. Initialize Watcher & Notifier based on OS
	var logWatcher watcher.LogWatcher
	var notify notifier.Notifier

	switch runtime.GOOS {
	case "linux":
		logWatcher = watcher.NewLinuxJournalWatcher()
		notify = notifier.NewLinuxNotifier()
	case "darwin":
		logWatcher = watcher.NewMacOSLogWatcher()
		notify = notifier.NewMacOSNotifier()
	default:
		fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
		os.Exit(1)
	}

	// 3. Start Watching
	logs, errs, err := logWatcher.Start()
	if err != nil {
		fmt.Printf("Error starting watcher: %v\n", err)
		os.Exit(1)
	}
	defer logWatcher.Stop()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("[*] Watching logs...")

	go func() {
		<-sigChan
		fmt.Println("\n[*] Shutting down...")
		logWatcher.Stop()
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
