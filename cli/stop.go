package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/spf13/cobra"
)

func stopCommand(t *core.Track) *cobra.Command {
	stop := &cobra.Command{
		Use:     "stop",
		Short:   "Stop the current record",
		Aliases: []string{"x"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			record, err := t.StopRecord(time.Now())

			if err != nil {
				out.Err("failed to stop record: %s", err)
				return
			}

			out.Success("Stopped record in '%s' at %02d:%02d", record.Project, record.End.Hour(), record.End.Minute())
		},
	}

	return stop
}
