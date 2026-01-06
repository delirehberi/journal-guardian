package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeminiProvider struct {
	APIKey string
	Model  string
	APIURL string // Added for testing
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content GeminiContent `json:"content"`
	} `json:"candidates"`
}

func (p *GeminiProvider) Name() string {
	return fmt.Sprintf("Gemini (%s)", p.Model)
}

func (p *GeminiProvider) Generate(prompt string) (string, error) {
	// Gemini api is weird, it doesn't take system prompt as a separate field easily in v1beta/generateContent for all models,
	// so we prepend it to the user prompt for simplicity.
	fullPrompt := fmt.Sprintf("%s\n\nError: %s", SYSTEM_PROMPT, prompt)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: fullPrompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	// In real usage, the URL includes the model. For testing, we might want to override completely.
	// If APIURL is set (testing), use it directly. Otherwise construct standard URL.
	url := p.APIURL
	if url == "" {
		url = fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.Model, p.APIKey)
	} else {
		// Mock logic or just append key if needed? For unit tests, we usually mock the exact path or ignore query params
		// Let's assume the test server handles the path provided by the caller or we just use the base URL for testing.
		// A simple way is: if testing, p.APIURL is the FULL endpoint.
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error contacting Gemini: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini error %s: %s", resp.Status, body)
	}

	var gResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return "", fmt.Errorf("error decoding Gemini response: %v", err)
	}

	if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content returned from Gemini")
	}

	return gResp.Candidates[0].Content.Parts[0].Text, nil
}
