package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestNewTCommandRejectsBlankTitle(t *testing.T) {
	cmd := NewTCommand(strings.NewReader(""), io.Discard, io.Discard)
	cmd.SetArgs([]string{"   "})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected an error for blank title, got nil")
	}

	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestNewTCommandRejectsLongTitle(t *testing.T) {
	cmd := NewTCommand(strings.NewReader(""), io.Discard, io.Discard)
	longTitle := strings.Repeat("a", titleCharLimit+1)
	cmd.SetArgs([]string{longTitle})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected an error for long title, got nil")
	}

	expected := fmt.Sprintf("todo title must be %d characters or fewer", titleCharLimit)
	if err.Error() != expected {
		t.Fatalf("expected error %q, got %q", expected, err.Error())
	}
}
