package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	ServiceType         string `yaml:"service_type"`
	ServicesFile        string `yaml:"services_file"`
	MongoURI            string `yaml:"mongo_uri"`
	MongoDatabaseName   string `yaml:"mongo_db_name"`
	MongoCollectionName string `yaml:"mongo_collection_name"`
}

type GatewayConfig struct {
	Addr string `yaml:"addr"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	GlobalConfig  GlobalConfig  `yaml:"global"`
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
