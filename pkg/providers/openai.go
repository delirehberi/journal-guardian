package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIProvider struct {
	APIKey string
	Model  string
	APIURL string // Added for testing
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
}

func (p *OpenAIProvider) Name() string {
	return fmt.Sprintf("OpenAI (%s)", p.Model)
}

func (p *OpenAIProvider) Generate(prompt string) (string, error) {
	reqBody := OpenAIRequest{
		Model: p.Model,
		Messages: []OpenAIMessage{
			{Role: "system", Content: SYSTEM_PROMPT},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	url := p.APIURL
	if url == "" {
		url = "https://api.openai.com/v1/chat/completions"
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error contacting OpenAI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai error %s: %s", resp.Status, body)
	}

	var oResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return "", fmt.Errorf("error decoding OpenAI response: %v", err)
	}

	if len(oResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return oResp.Choices[0].Message.Content, nil
}
