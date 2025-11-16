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

// Config represents the raw theme configuration values.
type Config struct {
	Text      string
	Muted     string
	Highlight string
	Success   string
	Worry     string
}

// Theme defines the theme for the TUI and CLI help docs.
type Theme struct {
	Text      PaletteColor
	Muted     PaletteColor
	Highlight PaletteColor
	Success   PaletteColor
	Worry     PaletteColor

	// UI characters
	CursorChar string // Character shown next to selected todo item
}

// PaletteColor keeps a hex colour string for LipGloss and an RGBA version for Fang.
type PaletteColor struct {
	raw  string
	rgba color.RGBA
}

// LipGloss converts the color to a lipgloss.Color.
func (c PaletteColor) LipGloss() lipgloss.Color {
	return lipgloss.Color(c.raw)
}

// RGBA exposes an image/color compliant representation of the color.
func (c PaletteColor) RGBA() color.RGBA {
	return c.rgba
}

// Default returns the default theme.
func Default() Theme {
	return MustFromConfig(Config{
		Text:      "#FFFFFF",
		Muted:     "#696969",
		Highlight: "#58C5C7",
		Success:   "#99CC00",
		Worry:     "#FF7676",
	})
}

// MustFromConfig creates a Theme from a Config and panics if parsing fails.
func MustFromConfig(cfg Config) Theme {
	t, err := FromConfig(cfg)
	if err != nil {
		panic(fmt.Sprintf("invalid theme configuration: %v", err))
	}
	return t
}

// FromConfig converts a Config into a ready-to-use Theme.
func FromConfig(cfg Config) (Theme, error) {
	text, err := parsePaletteColor(cfg.Text)
	if err != nil {
		return Theme{}, fmt.Errorf("parse text colour: %w", err)
	}

	muted, err := parsePaletteColor(cfg.Muted)
	if err != nil {
		return Theme{}, fmt.Errorf("parse muted colour: %w", err)
	}

	highlight, err := parsePaletteColor(cfg.Highlight)
	if err != nil {
		return Theme{}, fmt.Errorf("parse highlight colour: %w", err)
	}

	success, err := parsePaletteColor(cfg.Success)
	if err != nil {
		return Theme{}, fmt.Errorf("parse success colour: %w", err)
	}

	worry, err := parsePaletteColor(cfg.Worry)
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
		// Sound.
	default:
		return "", fmt.Errorf("hex colour must be 3 or 6 characters, got %d", len(s))
	}

	for _, r := range s {
		if !(r >= '0' && r <= '9') && !(r >= 'a' && r <= 'f') {
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
