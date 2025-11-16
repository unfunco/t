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
		Text:      "bogus",
		Muted:     "#333",
		Highlight: "#58C5C7",
		Success:   "#99CC00",
		Worry:     "#ff7676",
	}

	if _, err := FromConfig(cfg); err == nil {
		t.Fatal("expected FromConfig to error for invalid colour")
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
