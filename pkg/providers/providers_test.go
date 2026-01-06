package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOllamaGenerate(t *testing.T) {
	expectedResponse := "Use sudo to fix permissions"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		var req OllamaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}
		
		if !strings.Contains(req.Prompt, "test log entry") {
			t.Errorf("Expected prompt to contain test log, got: %s", req.Prompt)
		}

		resp := OllamaResponse{Response: expectedResponse}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	provider := &OllamaProvider{URL: ts.URL, Model: "llama2"}
	got, err := provider.Generate("test log entry")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if got != expectedResponse {
		t.Errorf("Expected %q, got %q", expectedResponse, got)
	}
}

func TestOpenAIGenerate(t *testing.T) {
	expectedResponse := "Check your config file"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Missing or invalid Authorization header")
		}

		var req OpenAIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request")
		}
		
		resp := OpenAIResponse{
			Choices: []struct {
				Message OpenAIMessage `json:"message"`
			}{
				{Message: OpenAIMessage{Content: expectedResponse}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	provider := &OpenAIProvider{APIKey: "test-key", Model: "gpt-4", APIURL: ts.URL}
	got, err := provider.Generate("test log")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if got != expectedResponse {
		t.Errorf("Expected %q, got %q", expectedResponse, got)
	}
}

func TestGeminiGenerate(t *testing.T) {
	expectedResponse := "Restart the service"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL contains key if we were testing full URL logic, but here we override APIURL
		// so we just check Body
		var req GeminiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request")
		}

		resp := GeminiResponse{
			Candidates: []struct {
				Content GeminiContent `json:"content"`
			}{
				{
					Content: GeminiContent{
						Parts: []GeminiPart{{Text: expectedResponse}},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	provider := &GeminiProvider{APIKey: "gemini-key", Model: "gemini-pro", APIURL: ts.URL}
	got, err := provider.Generate("test log")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if got != expectedResponse {
		t.Errorf("Expected %q, got %q", expectedResponse, got)
	}
}

func TestClaudeGenerate(t *testing.T) {
	expectedResponse := "Check logs in /var/log"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "claude-key" {
			t.Errorf("Missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("Missing anthropic-version header")
		}

		resp := ClaudeResponse{
			Content: []struct {
				Text string `json:"text"`
			}{
				{Text: expectedResponse},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	provider := &ClaudeProvider{APIKey: "claude-key", Model: "claude-3", APIURL: ts.URL}
	got, err := provider.Generate("test log")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if got != expectedResponse {
		t.Errorf("Expected %q, got %q", expectedResponse, got)
	}
}
