// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package theme

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Mode determines which palette should be active.
type Mode string

const (
	// ModeAuto chooses the palette based on the terminal background.
	ModeAuto Mode = "auto"
	// ModeDark always uses the dark palette.
	ModeDark Mode = "dark"
	// ModeLight always uses the light palette.
	ModeLight Mode = "light"
)

// Config represents the raw theme configuration values.
type Config struct {
	Mode  Mode          `json:"mode"`
	Light PaletteConfig `json:"light"`
	Dark  PaletteConfig `json:"dark"`
}

// PaletteConfig describes a single colour palette.
type PaletteConfig struct {
	Text      string `json:"text"`
	Muted     string `json:"muted"`
	Highlight string `json:"highlight"`
	Success   string `json:"success"`
	Worry     string `json:"worry"`
}

// DefaultConfig returns the built-in theme configuration.
func DefaultConfig() Config {
	return Config{
		Mode:  ModeAuto,
		Light: DefaultLightPalette(),
		Dark:  DefaultDarkPalette(),
	}
}

// DefaultDarkPalette returns the built-in palette optimised for dark backgrounds.
func DefaultDarkPalette() PaletteConfig {
	return PaletteConfig{
		Text:      "#FFFFFF",
		Muted:     "#696969",
		Highlight: "#58C5C7",
		Success:   "#99CC00",
		Worry:     "#FF7676",
	}
}

// DefaultLightPalette returns the built-in palette optimised for light backgrounds.
func DefaultLightPalette() PaletteConfig {
	return PaletteConfig{
		Text:      "#121417",
		Muted:     "#61646B",
		Highlight: "#205CBE",
		Success:   "#007A3B",
		Worry:     "#C62828",
	}
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

	out.Light = out.Light.withDefaults(def.Light)
	out.Dark = out.Dark.withDefaults(def.Dark)

	return out
}

func (pc *PaletteConfig) withDefaults(fallback PaletteConfig) PaletteConfig {
	if pc == nil {
		return fallback
	}

	out := *pc

	if strings.TrimSpace(out.Text) == "" {
		out.Text = fallback.Text
	}
	if strings.TrimSpace(out.Muted) == "" {
		out.Muted = fallback.Muted
	}
	if strings.TrimSpace(out.Highlight) == "" {
		out.Highlight = fallback.Highlight
	}
	if strings.TrimSpace(out.Success) == "" {
		out.Success = fallback.Success
	}
	if strings.TrimSpace(out.Worry) == "" {
		out.Worry = fallback.Worry
	}

	return out
}

func (cfg *Config) paletteForMode(hasDarkBackground bool) (PaletteConfig, error) {
	if cfg == nil {
		def := DefaultConfig()
		return def.paletteForMode(hasDarkBackground)
	}

	mode := strings.ToLower(strings.TrimSpace(string(cfg.Mode)))
	switch mode {
	case "", string(ModeAuto):
		if hasDarkBackground {
			return cfg.Dark, nil
		}
		return cfg.Light, nil
	case string(ModeDark):
		return cfg.Dark, nil
	case string(ModeLight):
		return cfg.Light, nil
	default:
		return PaletteConfig{}, fmt.Errorf("unknown theme mode %q", cfg.Mode)
	}
}

func (pc *PaletteConfig) toTheme() (Theme, error) {
	if pc == nil {
		return Theme{}, fmt.Errorf("palette configuration cannot be nil")
	}

	text, err := parsePaletteColor(pc.Text)
	if err != nil {
		return Theme{}, fmt.Errorf("parse text colour: %w", err)
	}

	muted, err := parsePaletteColor(pc.Muted)
	if err != nil {
		return Theme{}, fmt.Errorf("parse muted colour: %w", err)
	}

	highlight, err := parsePaletteColor(pc.Highlight)
	if err != nil {
		return Theme{}, fmt.Errorf("parse highlight colour: %w", err)
	}

	success, err := parsePaletteColor(pc.Success)
	if err != nil {
		return Theme{}, fmt.Errorf("parse success colour: %w", err)
	}

	worry, err := parsePaletteColor(pc.Worry)
	if err != nil {
		return Theme{}, fmt.Errorf("parse worry colour: %w", err)
	}

	return Theme{
		Text:       text,
		Muted:      muted,
		Highlight:  highlight,
		Success:    success,
		Worry:      worry,
		CursorChar: "❯",
	}, nil
}

// Theme defines the theme for the TUI and CLI help docs.
type Theme struct {
	Text      PaletteColor
	Muted     PaletteColor
	Highlight PaletteColor
	Success   PaletteColor
	Worry     PaletteColor

	CursorChar string // Character shown next to selected todo item.
}

// PaletteColor keeps a hex colour string for LipGloss and an RGBA version
// for Fang.
type PaletteColor struct {
	raw  string
	rgba color.RGBA
}

// LipGloss converts the color to a lipgloss.Color.
func (c *PaletteColor) LipGloss() lipgloss.Color {
	if c == nil {
		return lipgloss.Color("")
	}
	return lipgloss.Color(c.raw)
}

// RGBA exposes an image/color compliant representation of the color.
func (c *PaletteColor) RGBA() color.RGBA {
	if c == nil {
		return color.RGBA{}
	}
	return c.rgba
}

// Default returns the default theme optimised for dark backgrounds.
func Default() Theme {
	return MustFromConfig(DefaultConfig(), true)
}

// MustFromConfig creates a Theme from a Config and panics if parsing fails.
func MustFromConfig(cfg Config, hasDarkBackground bool) Theme {
	t, err := FromConfig(cfg, hasDarkBackground)
	if err != nil {
		panic(fmt.Sprintf("invalid theme configuration: %v", err))
	}

	return t
}

// FromConfig converts a Config into a ready-to-use Theme.
func FromConfig(cfg Config, hasDarkBackground bool) (Theme, error) {
	cfg = cfg.withDefaults()

	palette, err := cfg.paletteForMode(hasDarkBackground)
	if err != nil {
		return Theme{}, err
	}

	return palette.toTheme()
}

func parsePaletteColor(input string) (PaletteColor, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return PaletteColor{}, fmt.Errorf("colour cannot be blank")
	}

	hexValue, err := normalizeHex(trimmed)
	if err != nil {
		return PaletteColor{}, err
	}

	rgba, err := hexToRGBA(hexValue)
	if err != nil {
		return PaletteColor{}, err
	}

	return PaletteColor{raw: hexValue, rgba: rgba}, nil
}

func normalizeHex(input string) (string, error) {
	s := strings.TrimPrefix(input, "#")
	s = strings.ToLower(s)

	switch len(s) {
	case 3:
		s = fmt.Sprintf("%c%c%c%c%c%c", s[0], s[0], s[1], s[1], s[2], s[2])
	case 6:
		// All good in the hood.
	default:
		return "", fmt.Errorf("hex colour must be 3 or 6 characters, got %d", len(s))
	}

	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return "", fmt.Errorf("invalid hex digit %q", r)
		}
	}

	return "#" + strings.ToUpper(s), nil
}

func hexToRGBA(hex string) (color.RGBA, error) {
	s := strings.TrimPrefix(hex, "#")
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("hex colour must be 6 characters, got %d", len(s))
	}

	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("parse red component: %w", err)
	}

	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("parse green component: %w", err)
	}

	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("parse blue component: %w", err)
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xFF}, nil
}

// ContainerStyle returns the style for the main container.
func (t *Theme) ContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2)
}

// TabStyle returns the style for inactive tabs.
func (t *Theme) TabStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(t.Muted.LipGloss())
}

// ActiveTabStyle returns the style for the active tab.
func (t *Theme) ActiveTabStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(t.Text.LipGloss()).
		Background(t.Highlight.LipGloss()).
		Bold(true)
}

// ItemStyle returns the base style for todo items.
func (t *Theme) ItemStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

// HighlightedItemStyle returns the style for the currently selected todo item.
func (t *Theme) HighlightedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Highlight.LipGloss())
}

// CompletedTitleStyle returns the style for completed todo titles.
func (t *Theme) CompletedTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Muted.LipGloss()).
		Strikethrough(true)
}

// DescriptionStyle returns the style for todo descriptions.
func (t *Theme) DescriptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Muted.LipGloss())
}

// HelpStyle returns the style for help text.
func (t *Theme) HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Muted.LipGloss())
}

// SuccessStyle returns the style for success indicators.
func (t *Theme) SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Success.LipGloss())
}

// WorryStyle returns the style for warning/error indicators.
func (t *Theme) WorryStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Worry.LipGloss())
}
