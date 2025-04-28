package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     string   `json:"port"`
	Backends []string `json:"backends"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
