package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClaudeProvider struct {
	APIKey string
	Model  string
	APIURL string // Added for testing
}

type ClaudeRequest struct {
	Model     string          `json:"model"`
	Messages  []ClaudeMessage `json:"messages"`
	System    string          `json:"system,omitempty"`
	MaxTokens int             `json:"max_tokens"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (p *ClaudeProvider) Name() string {
	return fmt.Sprintf("Claude (%s)", p.Model)
}

func (p *ClaudeProvider) Generate(prompt string) (string, error) {
	reqBody := ClaudeRequest{
		Model:     p.Model,
		MaxTokens: 1024,
		System:    SYSTEM_PROMPT,
		Messages: []ClaudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	url := p.APIURL
	if url == "" {
		url = "https://api.anthropic.com/v1/messages"
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error contacting Claude: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("claude error %s: %s", resp.Status, body)
	}

	var cResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return "", fmt.Errorf("error decoding Claude response: %v", err)
	}

	if len(cResp.Content) == 0 {
		return "", fmt.Errorf("no content returned from Claude")
	}

	return cResp.Content[0].Text, nil
}
