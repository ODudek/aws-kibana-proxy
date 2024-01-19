package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type AppConfig struct {
	Port       string `yaml:"port"`
	EsEndpoint string `yaml:"es_endpoint"`
}

func LoadConfig(filename string) (AppConfig, error) {
	config := AppConfig{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return AppConfig{}, fmt.Errorf("error reading config file: %s", err.Error())
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing config file: %s", err.Error())
	}

	return config, nil
}

func NewAppConfig() *AppConfig {
	config, err := LoadConfig("./config.yaml")
	if err != nil {
		// stop application when there's problem with loading config
		panic(err)
	}
	return &config
}
