package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log_watcher/pkg/providers"
	"os"
	"os/exec"
)

// Structures for JSON parsing
type JournalEntry struct {
	Message    string `json:"MESSAGE"`
	SystemUnit string `json:"_SYSTEMD_UNIT"`
}

func sendNotification(title, message string) {
	if message == "" {
		return
	}
	cmd := exec.Command("notify-send", title, message)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}
}

func main() {
	provider, err := providers.NewProvider()
	if err != nil {
		fmt.Printf("Error initializing provider: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[*] Go Journal Watcher started using %s\n", provider.Name())

	cmd := exec.Command("journalctl", "-f", "-p", "3", "-o", "json", "-n", "0")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Bytes()
		var entry JournalEntry

		// Attempt to parse JSON line
		if err := json.Unmarshal(line, &entry); err == nil {
			// Filter out empty messages
			if entry.Message == "" {
				continue
			}

			fmt.Printf("\n[!] Error in %s:\n    %s\n", entry.SystemUnit, entry.Message)
			fmt.Println("    Asking AI...")

			suggestion, err := provider.Generate(fmt.Sprintf("Service: %s. Log: %s", entry.SystemUnit, entry.Message))
			if err != nil {
				fmt.Printf("Error generating suggestion: %v\n", err)
				continue
			}

			fmt.Printf("\n--- 💡 FIX SUGGESTION ---\n%s\n---------------------\n", suggestion)

			sendNotification(fmt.Sprintf("Error: %s", entry.SystemUnit), suggestion)
		}
	}
}
