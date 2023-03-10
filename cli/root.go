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

Show the records of the current week:
$ track report week

Subcommands
-----------`,
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
	root.AddCommand(pauseCommand(t))
	root.AddCommand(reportCommand(t))
	root.AddCommand(editCommand(t))
	root.AddCommand(deleteCommand(t))
	root.AddCommand(exportCommand(t))
	root.AddCommand(workspaceCommand(t))
	root.AddCommand(moveCommand(t))

	root.Long += "\n\n" + formatCmdTree(root)

	return root
}
