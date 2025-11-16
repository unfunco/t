// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

// Package config handles loading and saving application configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/unfunco/t/internal/icons"
	"github.com/unfunco/t/internal/paths"
	"github.com/unfunco/t/internal/theme"
)

const configFilename = "config.json"

// Config captures the configurable application properties.
type Config struct {
	Icons icons.Config `json:"icons"`
	Theme theme.Config `json:"theme"`
}

// Load retrieves the configuration from the default data directory.
func Load() (Config, error) {
	dataDir, err := paths.DefaultConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("determine config directory: %w", err)
	}

	return LoadFromDir(dataDir)
}

// LoadFromDir retrieves the configuration from the provided directory.
func LoadFromDir(configDir string) (Config, error) {
	if configDir == "" {
		return Config{}, fmt.Errorf("config directory cannot be empty")
	}

	cfg := Config{
		Theme: theme.DefaultConfig(),
		Icons: icons.DefaultConfig(),
	}

	configPath := filepath.Join(configDir, configFilename)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}

		return Config{}, fmt.Errorf("read config: %w", err)
	}

	if len(data) == 0 {
		return cfg, nil
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	return cfg, nil
}

// Path returns the default configuration file path.
func Path() (string, error) {
	configDir, err := paths.DefaultConfigDir()
	if err != nil {
		return "", fmt.Errorf("determine config directory: %w", err)
	}

	return filepath.Join(configDir, configFilename), nil
}
