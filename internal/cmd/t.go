package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/unfunco/t/internal/model"
	"github.com/unfunco/t/internal/storage"
	"github.com/unfunco/t/internal/tui"
	"github.com/unfunco/t/internal/version"
)

var ErrAmbiguousDateFlags = errors.New("only one of --today and --tomorrow may be specified")

// NewDefaultTCommand returns a new t command configured with the standard
// input, output, and error file descriptors.
func NewDefaultTCommand() *cobra.Command {
	return NewTCommand(os.Stdin, os.Stdout, os.Stderr)
}

// NewTCommand returns a new t command configured with the given input, output,
// and error file descriptors.
func NewTCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	var (
		today    bool
		tomorrow bool
	)

	t := &cobra.Command{
		Use:   "t [title] [--flags]",
		Short: "Interactive TODO manager for the CLI.",
		Long: heredoc.Doc(`
			Interactive TODO manager for the command line.
			Run without arguments to open the interactive interface.
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
			# Open interactive interface.
			t

			# Add new TODOs.
			t "Do something"
			t "Do something today" --today
			t "Do something tomorrow" --tomorrow
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if today && tomorrow {
				return ErrAmbiguousDateFlags
			}

			store, err := storage.New()
			if err != nil {
				return fmt.Errorf("failed to initialize storage: %w", err)
			}

			// Launch the TUI if no title argument is provided.
			if len(args) == 0 {
				todayList, err := store.LoadTodayList()
				if err != nil {
					return fmt.Errorf("failed to load today list: %w", err)
				}

				tomorrowList, err := store.LoadTomorrowList()
				if err != nil {
					return fmt.Errorf("failed to load tomorrow list: %w", err)
				}

				todoList, err := store.LoadTodoList()
				if err != nil {
					return fmt.Errorf("failed to load todo list: %w", err)
				}

				m := tui.New(todayList, tomorrowList, todoList)
				p := tea.NewProgram(&m)

				tuiModel, err := p.Run()
				if err != nil {
					return fmt.Errorf("error running TUI: %w", err)
				}

				if m, ok := tuiModel.(*tui.Model); ok && m.WasSubmitted() {
					if err := store.SaveAll(m.GetTodayList(), m.GetTomorrowList(), m.GetTodosList()); err != nil {
						return fmt.Errorf("failed to save changes: %w", err)
					}
				}

				return nil
			}

			title := args[0]
			todo := model.NewTodo(title, "")

			var targetList *model.TodoList
			if today {
				targetList, err = store.LoadTodayList()
				if err != nil {
					return fmt.Errorf("failed to load today list: %w", err)
				}
				targetList.Todos = append(targetList.Todos, todo)
				if err := store.SaveToday(targetList); err != nil {
					return fmt.Errorf("failed to save today list: %w", err)
				}
				_, _ = fmt.Fprintf(out, "Added to Today: %s\n", title)
			} else if tomorrow {
				targetList, err = store.LoadTomorrowList()
				if err != nil {
					return fmt.Errorf("failed to load tomorrow list: %w", err)
				}
				targetList.Todos = append(targetList.Todos, todo)
				if err := store.SaveTomorrow(targetList); err != nil {
					return fmt.Errorf("failed to save tomorrow list: %w", err)
				}
				_, _ = fmt.Fprintf(out, "Added to Tomorrow: %s\n", title)
			} else {
				targetList, err = store.LoadTodoList()
				if err != nil {
					return fmt.Errorf("failed to load todo list: %w", err)
				}
				targetList.Todos = append(targetList.Todos, todo)
				if err := store.SaveTodo(targetList); err != nil {
					return fmt.Errorf("failed to save todo list: %w", err)
				}
				_, _ = fmt.Fprintf(out, "Added to Todos: %s\n", title)
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

	t.Flags().BoolVar(&today, "today", false, "Add a TODO for today")
	t.Flags().BoolVar(&tomorrow, "tomorrow", false, "Add a TODO for tomorrow")

	return t
}
