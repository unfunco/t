// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package icons

import "testing"

func TestFromConfigUsesNerdIconsWhenForced(t *testing.T) {
	t.Setenv(iconModeEnvVar, "1")

	set, err := FromConfig(Config{})
	if err != nil {
		t.Fatalf("FromConfig() error = %v", err)
	}

	if set.HelpSeparator != "    " {
		t.Fatalf("expected nerd help separator, got %q", set.HelpSeparator)
	}

	if !set.ShowHelpIcons {
		t.Fatalf("expected nerd set to show help icons")
	}
}

func TestFromConfigFallsBackToPlainWhenDisabled(t *testing.T) {
	t.Setenv(iconModeEnvVar, "0")

	set, err := FromConfig(Config{})
	if err != nil {
		t.Fatalf("FromConfig() error = %v", err)
	}

	if set.HelpSeparator != "  " {
		t.Fatalf("expected plain help separator, got %q", set.HelpSeparator)
	}

	if set.ShowHelpIcons {
		t.Fatalf("expected plain set to hide help icons")
	}
}

func TestFromConfigErrorsOnUnknownMode(t *testing.T) {
	_, err := FromConfig(Config{Mode: "bogus"})
	if err == nil {
		t.Fatalf("expected error for unknown mode")
	}
}
