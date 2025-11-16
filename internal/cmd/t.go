// SPDX-FileCopyrightText: 2025 Daniel Morris <daniel@honestempire.com>
// SPDX-License-Identifier: MIT

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/MakeNowJust/heredoc/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/unfunco/t/internal/automation"
	"github.com/unfunco/t/internal/icons"
	"github.com/unfunco/t/internal/list"
	"github.com/unfunco/t/internal/model"
	"github.com/unfunco/t/internal/storage"
	"github.com/unfunco/t/internal/theme"
	"github.com/unfunco/t/internal/tui"
	"github.com/unfunco/t/internal/version"
)

const titleCharLimit = 100

var (
	ErrAmbiguousDateFlags = errors.New("only one of --today and --tomorrow may be specified")
	ErrEmptyTitle         = errors.New("todo title cannot be blank")
)

// NewDefaultTCommandWithTheme returns a new t command using the provided theme
// and icons, configured with the standard IO file descriptors.
func NewDefaultTCommandWithTheme(th theme.Theme, iconSet icons.Set) *cobra.Command {
	return NewTCommandWithTheme(os.Stdin, os.Stdout, os.Stderr, th, iconSet)
}

// NewTCommand returns a new t command configured with the given input, output,
// and error file descriptors.
func NewTCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	defIcons := icons.MustFromConfig(icons.DefaultConfig())
	return NewTCommandWithTheme(in, out, errOut, theme.Default(), defIcons)
}

// NewTCommandWithTheme returns a new t command configured with the provided
// input, output, error descriptors, theme, and icons.
func NewTCommandWithTheme(in io.Reader, out, errOut io.Writer, th theme.Theme, iconSet icons.Set) *cobra.Command {
	var (
		today    bool
		tomorrow bool
	)

	t := &cobra.Command{
		Use:   "t [title] [--flags]",
		Short: "Manage your todo lists in the CLI.",
		Long: heredoc.Doc(`
			The t command manages todo lists directly from the command line.
			Add new todos, or launch an interactive interface to view and manage
			your todos.
		`),
		Args:    cobra.MaximumNArgs(1),
		Version: version.SemanticVersion,
		Annotations: map[string]string{
			"versionInfo": fmt.Sprintf(
				"t %s (%s)",
				version.SemanticVersion,
				version.BuildDate,
			),
		},
		Example: heredoc.Doc(`
			# Add some todos.
			t "Do something"
			t "Do something today" --today
			t "Do something tomorrow" --tomorrow

			# Open the interactive interface.
			t
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if today && tomorrow {
				return ErrAmbiguousDateFlags
			}

			// Launch the TUI if no title argument is provided.
			if len(args) == 0 {
				store, err := storage.NewFileStorage()
				if err != nil {
					return fmt.Errorf("failed to initialise storage: %w", err)
				}

				lists, err := automation.Sync(store, time.Now())
				if err != nil {
					return fmt.Errorf("failed to prepare lists: %w", err)
				}

				m := tui.New(
					th,
					iconSet,
					lists[list.TodayID],
					lists[list.TomorrowID],
					lists[list.TodosID],
				)

				p := tea.NewProgram(&m)

				tuiModel, err := p.Run()
				if err != nil {
					return fmt.Errorf("error running TUI: %w", err)
				}

				if m, ok := tuiModel.(*tui.Model); ok && m.WasSubmitted() {
					if err := saveLists(store, m); err != nil {
						return err
					}
				}

				return nil
			}

			title := strings.TrimSpace(args[0])
			if err := validateTitle(title); err != nil {
				return err
			}

			store, err := storage.NewFileStorage()
			if err != nil {
				return fmt.Errorf("failed to initialise storage: %w", err)
			}

			if _, err := automation.Sync(store, time.Now()); err != nil {
				return fmt.Errorf("failed to prepare lists: %w", err)
			}

			var def list.Definition
			switch {
			case today:
				def = list.Today()
			case tomorrow:
				def = list.Tomorrow()
			default:
				def = list.Todos()
			}

			todo := model.NewTodo(title, "", list.DefaultDueDate(def.ID, time.Now()))

			if err := appendToList(store, def, &todo); err != nil {
				return err
			}

			return nil
		},
	}

	t.SetIn(in)
	t.SetOut(out)
	t.SetErr(errOut)

	t.CompletionOptions.DisableDefaultCmd = true
	t.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	t.Flags().BoolVar(&today, "today", false, "Add a todo for today")
	t.Flags().BoolVar(&tomorrow, "tomorrow", false, "Add a todo for tomorrow")

	return t
}

func saveLists(store storage.Storage, m *tui.Model) error {
	for _, def := range list.Default() {
		l := m.ListByID(def.ID)
		if l == nil {
			continue
		}
		if err := store.SaveList(def, l); err != nil {
			return fmt.Errorf("failed to save %s list: %w", def.Name, err)
		}
	}

	return nil
}

func appendToList(store storage.Storage, def list.Definition, todo *model.Todo) error {
	targetList, err := store.LoadList(def)
	if err != nil {
		return fmt.Errorf("failed to load %s list: %w", def.Name, err)
	}

	targetList.Todos = append(targetList.Todos, *todo)

	if err := store.SaveList(def, targetList); err != nil {
		return fmt.Errorf("failed to save %s list: %w", def.Name, err)
	}

	return nil
}

func validateTitle(t string) error {
	if t == "" {
		return ErrEmptyTitle
	}

	if utf8.RuneCountInString(t) > titleCharLimit {
		return fmt.Errorf("todo title must be %d characters or fewer", titleCharLimit)
	}

	return nil
}
