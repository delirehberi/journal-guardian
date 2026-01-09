package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type LogSourceConfig struct {
	Type   string            `json:"type"`   // "file", "journalctl", "macos"
	Params map[string]string `json:"params"` // e.g., "path" for file
}

type Config struct {
	Sources []LogSourceConfig `json:"sources"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// FindConfig searches for config.json in standard locations
func FindConfig() string {
	// 1. Current directory
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json"
	}

	// 2. User config directory
	configDir, err := os.UserConfigDir()
	if err == nil {
		userPath := filepath.Join(configDir, "log_watcher", "config.json")
		if _, err := os.Stat(userPath); err == nil {
			return userPath
		}
	}

	// 3. System global config (Linux mainly)
	if _, err := os.Stat("/etc/log_watcher/config.json"); err == nil {
		return "/etc/log_watcher/config.json"
	}

	return ""
}
