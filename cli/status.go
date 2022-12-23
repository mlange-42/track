package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func statusCommand(t *core.Track) *cobra.Command {
	status := &cobra.Command{
		Use:     "status",
		Short:   "Reports current status",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			if rec, ok := t.OpenRecord(); ok {
				out.Success(
					"Tracking project '%s' since %s (%s)",
					rec.Project,
					rec.Start.Format(util.TimeFormat),
					util.FormatDuration(rec.Duration()),
				)
				return
			}
			out.Success("No running record")
		},
	}

	return status
}
