package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ApiPort string `json:"api_port"`
	ApiURL  string `json:"api_url"`
	ApiKey  string `json:"api_key"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
