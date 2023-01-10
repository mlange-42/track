package cli

import (
	"fmt"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			ws := args[0]
			err := t.SwitchWorkspace(ws)
			if err != nil {
				return fmt.Errorf("failed to switch workspace: %s", err.Error())
			}

			out.Success("Switched to workspace '%s'", ws)
			return nil
		},
	}

	return workspace
}
