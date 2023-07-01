package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Debug  string `yaml:"debug"`
	Notify Notify `yaml:"notify"`
	Paths  Paths  `yaml:"paths"`
	Terms  Terms  `yaml:"terms"`
	Type   string `yaml:"type"`
}

type Notify struct {
	Urls []string `yaml:"urls"`
}

type Paths struct {
	Allow []string `yaml:"allow"`
	Block []string `yaml:"block"`
	Check []string `yaml:"check"`
}

type Terms struct {
	Allow []string `yaml:"allow"`
	Block []string `yaml:"block"`
}

func parseConfig(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read configuration file: %v", err)
	}

	config := &Config{}

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, fmt.Errorf("unable to parse configuration file: %v", err)
	}

	return config, nil
}
