package main

import (
	"context"
	"image/color"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/unfunco/t/internal/cmd"
	"github.com/unfunco/t/internal/theme"
)

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

func customColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	scheme := fang.AnsiColorScheme(c)
	th := theme.Default()

	scheme.Base = th.Text.RGBA()
	scheme.Description = th.Muted.RGBA()
	scheme.Comment = th.Muted.RGBA()
	scheme.Flag = th.Highlight.RGBA()
	scheme.FlagDefault = th.Muted.RGBA()
	scheme.Command = th.Highlight.RGBA()
	scheme.Program = th.Highlight.RGBA()
	scheme.QuotedString = th.Success.RGBA()
	scheme.Argument = th.Text.RGBA()
	scheme.DimmedArgument = th.Muted.RGBA()
	scheme.Help = th.Muted.RGBA()
	scheme.Dash = th.Text.RGBA()
	scheme.Title = th.Highlight.RGBA()
	scheme.ErrorHeader = [2]color.Color{th.Text.RGBA(), th.Worry.RGBA()}
	scheme.ErrorDetails = th.Worry.RGBA()

	return scheme
}
