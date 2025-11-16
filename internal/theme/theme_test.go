package theme

import (
	"image/color"
	"testing"
)

func TestParsePaletteColorHex(t *testing.T) {
	c, err := parsePaletteColor("#abcdef")
	if err != nil {
		t.Fatalf("parsePaletteColor returned error: %v", err)
	}

	want := color.RGBA{R: 0xAB, G: 0xCD, B: 0xEF, A: 0xFF}
	if c.RGBA() != want {
		t.Fatalf("expected RGBA %v, got %v", want, c.RGBA())
	}

	if c.LipGloss() != "#ABCDEF" {
		t.Fatalf("expected lipgloss colour #ABCDEF, got %s", c.LipGloss())
	}
}

func TestFromConfigInvalid(t *testing.T) {
	cfg := Config{
		Mode: ModeDark,
		Dark: PaletteConfig{
			Text:      "bogus",
			Muted:     "#333333",
			Highlight: "#58C5C7",
			Success:   "#99CC00",
			Worry:     "#FF7676",
		},
	}

	if _, err := FromConfig(cfg, true); err == nil {
		t.Fatal("expected FromConfig to error for invalid colour")
	}
}

func TestFromConfigAutoUsesLightPalette(t *testing.T) {
	cfg := Config{
		Mode: ModeAuto,
		Dark: PaletteConfig{
			Text:      "#010101",
			Muted:     "#020202",
			Highlight: "#030303",
			Success:   "#040404",
			Worry:     "#050505",
		},
		Light: PaletteConfig{
			Text:      "#A0A0A0",
			Muted:     "#A1A1A1",
			Highlight: "#A2A2A2",
			Success:   "#A3A3A3",
			Worry:     "#A4A4A4",
		},
	}

	th, err := FromConfig(cfg, false)
	if err != nil {
		t.Fatalf("FromConfig() error = %v", err)
	}

	if got := th.Text.LipGloss(); got != "#A0A0A0" {
		t.Fatalf("expected light palette text #A0A0A0, got %s", got)
	}
}

func TestDefaultTheme(t *testing.T) {
	th := Default()

	if th.Highlight.LipGloss() != "#58C5C7" {
		t.Fatalf("expected default highlight #58C5C7, got %s", th.Highlight.LipGloss())
	}

	if th.Success.RGBA() == (color.RGBA{}) {
		t.Fatal("expected success colour to be initialised")
	}

	if th.Worry.RGBA() == (color.RGBA{}) {
		t.Fatal("expected worry colour to be initialised")
	}
}

func TestParsePaletteColorShorthand(t *testing.T) {
	c, err := parsePaletteColor("#abc")
	if err != nil {
		t.Fatalf("parsePaletteColor returned error: %v", err)
	}

	want := color.RGBA{R: 0xAA, G: 0xBB, B: 0xCC, A: 0xFF}
	if c.RGBA() != want {
		t.Fatalf("expected RGBA %v, got %v", want, c.RGBA())
	}

	if c.LipGloss() != "#AABBCC" {
		t.Fatalf("expected lipgloss colour #AABBCC, got %s", c.LipGloss())
	}
}

func TestParsePaletteColorRejectsANSI(t *testing.T) {
	if _, err := parsePaletteColor("15"); err == nil {
		t.Fatal("expected ANSI value to be rejected")
	}
}
