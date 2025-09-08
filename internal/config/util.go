package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func UnmarshalConfigFile(file *os.File) (*Configuration, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("io ReadAll: %s", err.Error())
	}

	var cfg Configuration
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
