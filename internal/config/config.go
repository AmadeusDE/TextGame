package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	APIKey     string `json:"api_key"`
	ServerPort string `json:"server_port"`
	DataPath   string `json:"data_path"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	// Try loading from file
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		cfg.ServerPort = port
	}
	if key := os.Getenv("API_KEY"); key != "" {
		cfg.APIKey = key
	}
	if dPath := os.Getenv("DATA_PATH"); dPath != "" {
		cfg.DataPath = dPath
	}

	// Defaults
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	if cfg.APIKey == "" {
		cfg.APIKey = "dev-key"
	}
	if cfg.DataPath == "" {
		cfg.DataPath = "./data"
	}

	return &cfg, nil
}
