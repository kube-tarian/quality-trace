package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ClickhouseConnUrl string `yaml:"clickhouseConnUrl"`
}

func FromFile(file string) (Config, error) {
	var config Config

	filename, _ := filepath.Abs(file)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("read file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, fmt.Errorf("yaml unmarshal: %w", err)
	}

	return config, nil
}
