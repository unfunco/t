// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unfunco/t/internal/icons"
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

	if want := icons.DefaultConfig(); !iconConfigEqual(cfg.Icons, want) {
		t.Fatalf("icon config mismatch, want %+v got %+v", want, cfg.Icons)
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
  },
  "icons": {
    "mode": "plain",
    "nerd": {
      "helpSeparator": "  nerd  ",
      "add": "na",
      "cancel": "nc",
      "edit": "ne",
      "navigate": "nn",
      "overdue": "no",
      "select": "ns",
      "submit": "nS",
      "cursor": "n>",
      "showHelpIcons": true
    },
    "plain": {
      "helpSeparator": "  p  ",
      "add": "pa",
      "cancel": "pc",
      "edit": "pe",
      "navigate": "pn",
      "overdue": "po",
      "select": "ps",
      "submit": "pS",
      "cursor": "p>",
      "showHelpIcons": false
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

	wantIcons := icons.Config{
		Mode: icons.ModeASCII,
		NerdFont: icons.SetConfig{
			HelpSeparator: "  nerd  ",
			Add:           "na",
			Cancel:        "nc",
			Edit:          "ne",
			Navigate:      "nn",
			Overdue:       "no",
			Select:        "ns",
			Submit:        "nS",
			Cursor:        "n>",
			ShowHelpIcons: boolPtr(true),
		},
		ASCII: icons.SetConfig{
			HelpSeparator: "  p  ",
			Add:           "pa",
			Cancel:        "pc",
			Edit:          "pe",
			Navigate:      "pn",
			Overdue:       "po",
			Select:        "ps",
			Submit:        "pS",
			Cursor:        "p>",
			ShowHelpIcons: boolPtr(false),
		},
	}

	if !iconConfigEqual(cfg.Icons, wantIcons) {
		t.Fatalf("icons mismatch, want %+v got %+v", wantIcons, cfg.Icons)
	}
}

func boolPtr(v bool) *bool {
	return &v
}

func iconConfigEqual(a, b icons.Config) bool {
	if a.Mode != b.Mode {
		return false
	}

	return iconSetEqual(a.NerdFont, b.NerdFont) && iconSetEqual(a.ASCII, b.ASCII)
}

func iconSetEqual(a, b icons.SetConfig) bool {
	if a.HelpSeparator != b.HelpSeparator ||
		a.Add != b.Add ||
		a.Cancel != b.Cancel ||
		a.Edit != b.Edit ||
		a.Navigate != b.Navigate ||
		a.Overdue != b.Overdue ||
		a.Select != b.Select ||
		a.Submit != b.Submit ||
		a.Cursor != b.Cursor {
		return false
	}

	return boolValue(a.ShowHelpIcons) == boolValue(b.ShowHelpIcons)
}

func boolValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}
