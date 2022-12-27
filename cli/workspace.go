package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func workspaceCommand(t *core.Track) *cobra.Command {

	workspace := &cobra.Command{
		Use:     "workspace WORKSPACE",
		Short:   "Switch to another workspace",
		Long:    `Switch to another workspace`,
		Aliases: []string{"w"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			ws := args[0]
			err := t.SwitchWorkspace(ws)
			if err != nil {
				out.Err("failed to switch workspace: %s", err.Error())
				return
			}

			out.Success("Switched to workspace '%s'", ws)
		},
	}

	return workspace
}
