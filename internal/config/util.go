package config

import (
	"fmt"
	"io"
	"os"
	"strings"

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

	// 'Bot ' Prefix for Discord-API
	if !strings.HasPrefix(cfg.Discord_Settings.Token, "Bot ") {
		cfg.Discord_Settings.Token = "Bot " + cfg.Discord_Settings.Token
	}

	return &cfg, nil
}
