package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type GatewayConfig struct {
	Addr string `yaml:"addr"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	GatewayConfig GatewayConfig `yaml:"gateway"`
	LoggingConfig LoggingConfig `yaml:"logging"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
