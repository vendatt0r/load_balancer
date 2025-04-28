package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type RateLimitConfig struct {
	Capacity   int `yaml:"capacity"`
	RefillRate int `yaml:"refill_rate"`
}

type Config struct {
	Port             string                     `yaml:"port"`
	Backends         []string                   `yaml:"backends"`
	RateLimits       map[string]RateLimitConfig `yaml:"rate_limits"`
	DefaultRateLimit RateLimitConfig            `yaml:"default_rate_limit"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
