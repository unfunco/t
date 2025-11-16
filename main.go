package main

import (
	"context"
	"fmt"
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/unfunco/t/internal/cmd"
	"github.com/unfunco/t/internal/config"
	"github.com/unfunco/t/internal/icons"
	"github.com/unfunco/t/internal/theme"
)

func main() {
	configPath := ""
	if path, err := config.Path(); err == nil {
		configPath = path
	}

	cfg := config.Config{
		Theme: theme.DefaultConfig(),
		Icons: icons.DefaultConfig(),
	}
	if loadedCfg, err := config.Load(); err != nil {
		logConfigWarning(configPath, err)
	} else {
		cfg = loadedCfg
	}

	hasDarkBackground := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)

	th, err := theme.FromConfig(cfg.Theme, hasDarkBackground)
	if err != nil {
		logThemeWarning(configPath, err)
		th = theme.MustFromConfig(theme.DefaultConfig(), hasDarkBackground)
	}

	iconSet, err := icons.FromConfig(cfg.Icons)
	if err != nil {
		logIconWarning(configPath, err)
		iconSet = icons.MustFromConfig(icons.DefaultConfig())
	}

	if err := fang.Execute(
		context.Background(),
		cmd.NewDefaultTCommandWithTheme(th, iconSet),
		fang.WithColorSchemeFunc(func(c lipgloss.LightDarkFunc) fang.ColorScheme {
			return customColorScheme(c, th)
		}),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		os.Exit(1)
	}
}

func logConfigWarning(configPath string, err error) {
	if configPath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "warning: failed to load config: %v; using defaults\n", err)
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "warning: failed to load config at %s: %v; using defaults\n", configPath, err)
}

func logThemeWarning(configPath string, err error) {
	if configPath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "warning: invalid theme configuration: %v; using default theme\n", err)
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "warning: invalid theme configuration in %s: %v; using default theme\n", configPath, err)
}

func logIconWarning(configPath string, err error) {
	if configPath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "warning: invalid icon configuration: %v; using default icons\n", err)
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "warning: invalid icon configuration in %s: %v; using default icons\n", configPath, err)
}

func customColorScheme(c lipgloss.LightDarkFunc, th theme.Theme) fang.ColorScheme {
	scheme := fang.AnsiColorScheme(c)

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
