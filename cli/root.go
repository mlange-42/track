package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/spf13/cobra"
)

// RootCommand sets up the CLI
func RootCommand(t *core.Track, version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "track",
		Short: "Track is a time tracking command line tool",
		Long: `Track is a time tracking command line tool

Getting started
---------------

Create a project:
$ track create project my-project

Start tracking the project:
$ track start my-project

Stop tracking:
$ track stop

Show today's records:
$ track list records today

Show a daily timeline:
$ track report timeline days

`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(statusCommand(t))
	root.AddCommand(listCommand(t))
	root.AddCommand(createCommand(t))
	root.AddCommand(startCommand(t))
	root.AddCommand(stopCommand(t))
	root.AddCommand(resumeCommand(t))
	root.AddCommand(switchCommand(t))
	root.AddCommand(reportCommand(t))
	root.AddCommand(editCommand(t))

	// TODO: colors
	// TODO: cancel
	// TODO: reports
	// TODO: undo
	// TODO: config
	// TODO: rename "project" to "activity" or "topic"?
	// TODO: project categories?
	// TODO: use project hierarchy
	// TODO: allow multiple running projects? -> should be covered by hierarchy
	// TODO: simplified JSON time format for reports

	return root
}
