package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	HostName string `json:"hostname"`
	ApiPort  string `json:"api_port"`
	ApiURL   string `json:"api_url"`
	ApiKey   string `json:"api_key"`
}

func DefaultConfig() *Config {
	return &Config{
		HostName: "localhost",
		ApiPort:  "8080",
		ApiURL:   "https://data-api.coindesk.com",
		ApiKey:   "",
	}
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return DefaultConfig(), nil
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
