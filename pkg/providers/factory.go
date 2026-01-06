package providers

import (
	"fmt"
	"os"
	"strings"
)

func NewProvider() (LLMProvider, error) {
	providerType := os.Getenv("LLM_PROVIDER")
	model := os.Getenv("MODEL")
	if model == "" {
		model = "gpt-oss:20b" // Default for Ollama if not set
	}

	switch strings.ToLower(providerType) {
	case "openai":
		key := os.Getenv("OPENAI_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is required for openai provider")
		}
		if model == "" {
			model = "gpt-4"
		}
		return &OpenAIProvider{APIKey: key, Model: model}, nil

	case "gemini":
		key := os.Getenv("GEMINI_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY is required for gemini provider")
		}
		if model == "" {
			model = "gemini-pro"
		}
		return &GeminiProvider{APIKey: key, Model: model}, nil

	case "claude":
		key := os.Getenv("ANTHROPIC_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY is required for claude provider")
		}
		if model == "" {
			model = "claude-3-opus-20240229"
		}
		return &ClaudeProvider{APIKey: key, Model: model}, nil

	case "ollama", "": // Default
		url := os.Getenv("OLLAMA_URL")
		if url == "" {
			url = "http://localhost:11434/api/generate"
		}
		return &OllamaProvider{URL: url, Model: model}, nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", providerType)
	}
}
