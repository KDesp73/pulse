package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration values from config.yml
type Config struct {
	MQTT struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		Topic    string `yaml:"topic"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mqtt"`
	Web struct {
		Port int    `yaml:"port"`
		Page string `yaml:"page"`
	} `yaml:"web"`
}

var GlobalConfig *Config

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}
