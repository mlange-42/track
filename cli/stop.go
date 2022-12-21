package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/spf13/cobra"
)

func stopCommand(t *core.Track) *cobra.Command {
	stop := &cobra.Command{
		Use:   "stop",
		Short: "Stop the current record",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			record, ok := t.OpenRecord()
			if !ok {
				out.Err("failed to stop record: no running record")
				return
			}

			record.End = time.Now()

			err := t.SaveRecord(record, true)
			if err != nil {
				out.Err("failed to stop record: %s", err)
				return
			}

			out.Success("Stopped record at %02d:%02d", record.End.Hour(), record.End.Minute())
		},
	}

	return stop
}
