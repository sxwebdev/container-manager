package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	BindAddress string          `yaml:"bind_address"`
	Port        string          `yaml:"port"`
	Services    []ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	Name        string `yaml:"name"`
	ProjectPath string `yaml:"project_path"`
	ComposeFile string `yaml:"compose_file"`
	Enabled     bool   `yaml:"enabled"`
}

func Load(configPath string) (Config, error) {
	// Default configuration values
	config := Config{
		BindAddress: "127.0.0.1",
		Port:        "8090",
		Services:    []ServiceConfig{},
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	// Environment variables override config
	if addr := os.Getenv("CONTAINER_MANAGER_BIND_ADDRESS"); addr != "" {
		config.BindAddress = addr
	}
	if port := os.Getenv("CONTAINER_MANAGER_PORT"); port != "" {
		config.Port = port
	}

	return config, nil
}
