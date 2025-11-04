package theme

import (
	"image/color"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color scheme for the application.
type Theme struct {
	// Primary highlight color (teal) - used for Fang CLI rendering
	Highlight color.RGBA

	// Lipgloss color values for TUI
	HighlightLipgloss lipgloss.Color
	ActiveText        lipgloss.Color
	MutedText         lipgloss.Color
	SuccessText       lipgloss.Color

	// UI characters
	CursorChar string // Character shown next to selected todo item
}

// Default returns the default teal theme.
func Default() Theme {
	teal := color.RGBA{R: 0x58, G: 0xC5, B: 0xC7, A: 0xFF}

	return Theme{
		Highlight:         teal,
		HighlightLipgloss: "#58C5C7",
		ActiveText:        "15",
		MutedText:         "240",
		SuccessText:       "42",
		CursorChar:        "❯",
	}
}

// ContainerStyle returns the style for the main container.
func (t *Theme) ContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2)
}

// TabStyle returns the style for inactive tabs.
func (t *Theme) TabStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(t.MutedText)
}

// ActiveTabStyle returns the style for the active tab.
func (t *Theme) ActiveTabStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(t.ActiveText).
		Background(t.HighlightLipgloss).
		Bold(true)
}

// ItemStyle returns the base style for todo items.
func (t *Theme) ItemStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

// HighlightedItemStyle returns the style for the currently selected todo item.
func (t *Theme) HighlightedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.HighlightLipgloss)
}

// CompletedTitleStyle returns the style for completed todo titles.
func (t *Theme) CompletedTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.MutedText).
		Strikethrough(true)
}

// DescriptionStyle returns the style for todo descriptions.
func (t *Theme) DescriptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.MutedText)
}

// HelpStyle returns the style for help text.
func (t *Theme) HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.MutedText)
}

// SuccessStyle returns the style for success indicators.
func (t *Theme) SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.SuccessText)
}
