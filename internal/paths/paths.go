// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

// Package paths centralises platform-specific filesystem locations.
package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

const app = "t"

// DefaultDataDir returns the standard directory used for mutable data.
func DefaultDataDir() (string, error) {
	return resolveDir("XDG_DATA_HOME", func(home string) string {
		return filepath.Join(home, ".local", "share")
	})
}

// DefaultConfigDir returns the standard directory used for configuration.
func DefaultConfigDir() (string, error) {
	return resolveDir("XDG_CONFIG_HOME", func(home string) string {
		return filepath.Join(home, ".config")
	})
}

func resolveDir(envVar string, fallback func(home string) string) (string, error) {
	if dir := os.Getenv(envVar); dir != "" {
		return filepath.Join(dir, app), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(fallback(home), app), nil
}
