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
    "mode": "auto",
    "dark": {
      "text": "#111111",
      "muted": "#222222",
      "highlight": "#333333",
      "success": "#444444",
      "worry": "#555555"
    },
    "light": {
      "text": "#aaaaaa",
      "muted": "#bbbbbb",
      "highlight": "#cccccc",
      "success": "#dddddd",
      "worry": "#eeeeee"
    }
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
		Mode: theme.ModeAuto,
		Dark: theme.PaletteConfig{
			Text:      "#111111",
			Muted:     "#222222",
			Highlight: "#333333",
			Success:   "#444444",
			Worry:     "#555555",
		},
		Light: theme.PaletteConfig{
			Text:      "#aaaaaa",
			Muted:     "#bbbbbb",
			Highlight: "#cccccc",
			Success:   "#dddddd",
			Worry:     "#eeeeee",
		},
	}

	if cfg.Theme != want {
		t.Fatalf("theme mismatch, want %+v got %+v", want, cfg.Theme)
	}
}
