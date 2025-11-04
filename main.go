package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/unfunco/t/internal/cmd"
	"github.com/unfunco/t/internal/theme"
)

// customColorScheme returns a custom color scheme with teal headings.
func customColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	scheme := fang.AnsiColorScheme(c)
	// Override the title color to match the teal used in the TUI
	scheme.Title = theme.Default().Highlight
	return scheme
}

func main() {
	if err := fang.Execute(
		context.Background(),
		cmd.NewDefaultTCommand(),
		fang.WithColorSchemeFunc(customColorScheme),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		os.Exit(1)
	}
}
