package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaProvider struct {
	URL   string
	Model string
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func (p *OllamaProvider) Name() string {
	return fmt.Sprintf("Ollama (%s)", p.Model)
}

func (p *OllamaProvider) Generate(prompt string) (string, error) {
	fullPrompt := fmt.Sprintf("%s\n\nError: %s", SYSTEM_PROMPT, prompt)
	fmt.Println("    Prompt to Ollama:", fullPrompt)

	reqBody := OllamaRequest{
		Model:  p.Model,
		Prompt: fullPrompt,
		Stream: false,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(p.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error contacting Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned non-OK status: %s", resp.Status)
	}

	var oResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return "", fmt.Errorf("error decoding Ollama response: %v", err)
	}

	return oResp.Response, nil
}
