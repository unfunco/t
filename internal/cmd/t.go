package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
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

			// Launch the TUI if no title argument is provided.
			if len(args) == 0 {
				m := tui.New()
				p := tea.NewProgram(&m)

				model, err := p.Run()
				if err != nil {
					return fmt.Errorf("error running TUI: %w", err)
				}

				if tuiModel, ok := model.(*tui.Model); ok {
					//nolint:SA9003
					if tuiModel.WasSubmitted() {
						// TODO(unfunco): Save to persistent storage.
					}
				}

				return nil
			}

			_, _ = fmt.Fprintf(out, "Adding todo: %s\n", args[0])

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
