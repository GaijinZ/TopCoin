package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Hostname       string `json:"hostname"`
	Port           string `json:"port"`
	ApiPort        string `json:"api_port"`
	ApiServiceName string `json:"api_service_name"`
	ApiURL         string `json:"api_url"`
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
