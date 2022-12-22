package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/spf13/cobra"
)

// RootCommand sets up the CLI
func RootCommand(t *core.Track, version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "track",
		Short:         "track is a time tracking command line tool",
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
	root.AddCommand(switchCommand(t))

	// TODO: cancel
	// TODO: reports

	return root
}
