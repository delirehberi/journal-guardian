package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

const (
	OLLAMA_URL_DEFAULT    = "http://localhost:11434/api/generate"
	MODEL_DEFAULT         = "llama3"
	SYSTEM_PROMPT = `
    You are a Linux SysAdmin. 
    I will provide you service and the log response from journalctl. You will understand the error and you give me suggestions on how to fix it. If its well known issues like missing files, permission denied, or wrong password with sudo. Do not give suggestions, just show the error about it like 'Permission denied' or 'No such file or directory for FILENAME that asked for SERVICE'. You are a problem solver, so do not make it complex simple problems but be careful about real problems. 
    `
)

var (
	OLLAMA_URL string
	MODEL      string
)

func init() {
	OLLAMA_URL = os.Getenv("OLLAMA_URL")
	if OLLAMA_URL == "" {
		OLLAMA_URL = OLLAMA_URL_DEFAULT
	}

	MODEL = os.Getenv("MODEL")
	if MODEL == "" {
		MODEL = MODEL_DEFAULT
	}
}

// Structures for JSON parsing
type JournalEntry struct {
	Message    string `json:"MESSAGE"`
	SystemUnit string `json:"_SYSTEMD_UNIT"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func getFixSuggestion(logMsg string) string {
    prompt:= fmt.Sprintf("%s\n\nError: %s", SYSTEM_PROMPT, logMsg)
    fmt.Println("    Prompt to Ollama:", prompt)
	reqBody := OllamaRequest{
		Model:  MODEL,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, _ := json.Marshal(reqBody)

    fmt.Sprintln("Data is %", string(jsonData))


	resp, err := http.Post(OLLAMA_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Sprintf("Error contacting Ollama: %v", err)
	}
	defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Sprintf("Ollama returned non-OK status: %s", resp.Status)
    }
    fmt.Println("    Received response from Ollama")
	var oResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return "Error decoding Ollama response"
	}
    fmt.Println("    Decoded Ollama response")
    fmt.Sprintln("    Suggestion:", oResp)
	return oResp.Response
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
	fmt.Printf("[*] Go Journal Watcher started (Model: %s)\n", MODEL)

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
			fmt.Println("    Asking Ollama...")

			suggestion := getFixSuggestion(fmt.Sprintf("Service: %s. Log: %s", entry.SystemUnit, entry.Message))
			fmt.Printf("\n--- 💡 OLLAMA FIX ---\n%s\n---------------------\n", suggestion)
			
			sendNotification(fmt.Sprintf("Error: %s", entry.SystemUnit), suggestion)
		}
	}
}
