// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unfunco/t/internal/theme"
)

func TestLoadFromDirReturnsDefaultsWhenMissing(t *testing.T) {
	dir := t.TempDir()

	cfg, err := LoadFromDir(dir)
	if err != nil {
		t.Fatalf("LoadFromDir() error = %v", err)
	}

	if want := theme.DefaultConfig(); cfg.Theme != want {
		t.Fatalf("theme mismatch, want %+v got %+v", want, cfg.Theme)
	}
}

func TestLoadFromDirReadsConfigFile(t *testing.T) {
	dir := t.TempDir()
	content := []byte(`{
  "theme": {
    "text": "#010101",
    "muted": "#020202",
    "highlight": "#030303",
    "success": "#040404",
    "worry": "#050505"
  }
}`)

	if err := os.WriteFile(filepath.Join(dir, "config.json"), content, 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadFromDir(dir)
	if err != nil {
		t.Fatalf("LoadFromDir() error = %v", err)
	}

	want := theme.Config{
		Text:      "#010101",
		Muted:     "#020202",
		Highlight: "#030303",
		Success:   "#040404",
		Worry:     "#050505",
	}

	if cfg.Theme != want {
		t.Fatalf("theme mismatch, want %+v got %+v", want, cfg.Theme)
	}
}
