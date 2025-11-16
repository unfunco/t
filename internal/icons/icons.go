// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package icons

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Mode determines which icon set should be used.
type Mode string

const (
	// ModeAuto chooses the icon set based on detected font support.
	ModeAuto Mode = "auto"
	// ModeNerdFont always uses the NerdFont Font friendly icon set.
	ModeNerdFont Mode = "nerd_font"
	// ModeASCII always uses the ASCII friendly icon set.
	ModeASCII Mode = "ascii"

	iconModeEnvVar = "T_ICON_MODE"
)

// Config represents the raw icon configuration values.
type Config struct {
	Mode     Mode      `json:"mode"`
	ASCII    SetConfig `json:"ascii"`
	NerdFont SetConfig `json:"nerd_font"`
}

// SetConfig describes a single icon set.
type SetConfig struct {
	HelpSeparator string `json:"helpSeparator"`
	Add           string `json:"add"`
	Cancel        string `json:"cancel"`
	Edit          string `json:"edit"`
	Navigate      string `json:"navigate"`
	Overdue       string `json:"overdue"`
	Select        string `json:"select"`
	Submit        string `json:"submit"`
	Cursor        string `json:"cursor"`

	ShowHelpIcons *bool `json:"showHelpIcons,omitempty"`
}

// Set captures the ready-to-use icons for the UI.
type Set struct {
	HelpSeparator string
	Add           string
	Cancel        string
	Edit          string
	Navigate      string
	Overdue       string
	Select        string
	Submit        string
	Cursor        string

	ShowHelpIcons bool
}

// DefaultConfig returns the built-in icon configuration.
func DefaultConfig() Config {
	return Config{
		Mode:     ModeAuto,
		ASCII:    defaultASCIIIconSet(),
		NerdFont: defaultNerdFontIconSet(),
	}
}

func defaultNerdFontIconSet() SetConfig {
	return SetConfig{
		HelpSeparator: "    ",
		Add:           "",
		Cancel:        "",
		Edit:          "",
		Navigate:      "",
		Overdue:       "",
		Select:        "",
		Submit:        "",
		Cursor:        "",
		ShowHelpIcons: boolPtr(true),
	}
}

func defaultASCIIIconSet() SetConfig {
	return SetConfig{
		HelpSeparator: "  ",
		Add:           "+",
		Cancel:        "x",
		Edit:          "~",
		Navigate:      "->",
		Overdue:       "!",
		Select:        "*",
		Submit:        "S",
		Cursor:        ">",
		ShowHelpIcons: boolPtr(false),
	}
}

func boolPtr(v bool) *bool {
	return &v
}

func (cfg *Config) withDefaults() Config {
	if cfg == nil {
		return DefaultConfig()
	}

	out := *cfg
	def := DefaultConfig()

	if strings.TrimSpace(string(out.Mode)) == "" {
		out.Mode = def.Mode
	}

	out.ASCII = out.ASCII.withDefaults(def.ASCII)
	out.NerdFont = out.NerdFont.withDefaults(def.NerdFont)

	return out
}

func (sc *SetConfig) withDefaults(fallback SetConfig) SetConfig {
	if sc == nil {
		return fallback
	}

	out := *sc

	if strings.TrimSpace(out.HelpSeparator) == "" {
		out.HelpSeparator = fallback.HelpSeparator
	}
	if strings.TrimSpace(out.Add) == "" {
		out.Add = fallback.Add
	}
	if strings.TrimSpace(out.Cancel) == "" {
		out.Cancel = fallback.Cancel
	}
	if strings.TrimSpace(out.Edit) == "" {
		out.Edit = fallback.Edit
	}
	if strings.TrimSpace(out.Navigate) == "" {
		out.Navigate = fallback.Navigate
	}
	if strings.TrimSpace(out.Overdue) == "" {
		out.Overdue = fallback.Overdue
	}
	if strings.TrimSpace(out.Select) == "" {
		out.Select = fallback.Select
	}
	if strings.TrimSpace(out.Submit) == "" {
		out.Submit = fallback.Submit
	}
	if strings.TrimSpace(out.Cursor) == "" {
		out.Cursor = fallback.Cursor
	}
	if out.ShowHelpIcons == nil {
		out.ShowHelpIcons = fallback.ShowHelpIcons
	}

	return out
}

// MustFromConfig creates an icon Set from a Config and panics if parsing fails.
func MustFromConfig(cfg Config) Set {
	set, err := FromConfig(cfg)
	if err != nil {
		panic(fmt.Sprintf("invalid icon configuration: %v", err))
	}

	return set
}

// FromConfig converts a Config into a ready-to-use Set.
func FromConfig(cfg Config) (Set, error) {
	cfg = cfg.withDefaults()

	mode := strings.ToLower(strings.TrimSpace(string(cfg.Mode)))
	switch mode {
	case "", string(ModeAuto):
		if nerdFontsEnabled() {
			return cfg.NerdFont.toSet(), nil
		}
		return cfg.ASCII.toSet(), nil
	case string(ModeNerdFont):
		return cfg.NerdFont.toSet(), nil
	case string(ModeASCII):
		return cfg.ASCII.toSet(), nil
	default:
		return Set{}, fmt.Errorf("unknown icon mode %q", cfg.Mode)
	}
}

func (sc *SetConfig) toSet() Set {
	if sc == nil {
		icons := defaultASCIIIconSet()
		return icons.toSet()
	}

	showHelp := false
	if sc.ShowHelpIcons != nil {
		showHelp = *sc.ShowHelpIcons
	}

	return Set{
		HelpSeparator: sc.HelpSeparator,
		Add:           sc.Add,
		Cancel:        sc.Cancel,
		Edit:          sc.Edit,
		Navigate:      sc.Navigate,
		Overdue:       sc.Overdue,
		Select:        sc.Select,
		Submit:        sc.Submit,
		Cursor:        sc.Cursor,
		ShowHelpIcons: showHelp,
	}
}

var (
	nerdFontsOnce     sync.Once
	nerdFontsDetected bool
)

func nerdFontsEnabled() bool {
	if forced, ok := forcedNerdFontPreference(); ok {
		return forced
	}

	nerdFontsOnce.Do(func() {
		nerdFontsDetected = detectNerdFontSupport()
	})

	return nerdFontsDetected
}

func forcedNerdFontPreference() (bool, bool) {
	val, ok := os.LookupEnv(iconModeEnvVar)
	if !ok {
		return false, false
	}

	val = strings.TrimSpace(val)
	forced, err := strconv.ParseBool(val)
	if err != nil {
		return false, false
	}

	return forced, true
}

func detectNerdFontSupport() bool {
	for _, dir := range fontSearchPaths() {
		if hasNerdFont(dir) {
			return true
		}
	}

	return false
}

func fontSearchPaths() []string {
	seen := make(map[string]struct{})
	add := func(path string, list []string) []string {
		path = strings.TrimSpace(path)
		if path == "" {
			return list
		}

		if _, ok := seen[path]; ok {
			return list
		}

		seen[path] = struct{}{}
		return append(list, path)
	}

	var dirs []string

	if home, err := os.UserHomeDir(); err == nil {
		dirs = add(filepath.Join(home, "Library", "Fonts"), dirs)
		dirs = add(filepath.Join(home, ".local", "share", "fonts"), dirs)
		dirs = add(filepath.Join(home, ".fonts"), dirs)
	}

	dirs = add("/Library/Fonts", dirs)
	dirs = add("/System/Library/Fonts", dirs)
	dirs = add("/System/Library/Fonts/Supplemental", dirs)
	dirs = add("/usr/share/fonts", dirs)
	dirs = add("/usr/local/share/fonts", dirs)

	return dirs
}

func hasNerdFont(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		name := strings.ToLower(entry.Name())
		if !entry.IsDir() && containsNerdFontMarker(name) {
			return true
		}

		if entry.IsDir() && containsNerdFontMarker(name) {
			return true
		}
	}

	return false
}

func containsNerdFontMarker(name string) bool {
	name = strings.ToLower(name)
	return strings.Contains(name, "nerd font") || strings.Contains(name, "nerdfont")
}
