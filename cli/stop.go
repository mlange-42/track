package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func stopCommand(t *core.Track) *cobra.Command {
	var deleteRecord bool

	stop := &cobra.Command{
		Use:     "stop",
		Short:   "Stop the current record",
		Aliases: []string{"x"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			open, ok := t.OpenRecord()
			if !ok {
				out.Err("failed to stop record: no record running")
				return
			}

			if deleteRecord && !confirmDeleteRecord(open) {
				out.Err("failed to stop record: aborted by user")
				return
			}

			record, err := t.StopRecord(time.Now())
			if err != nil {
				out.Err("failed to stop record: %s", err)
				return
			}
			out.Success("Stopped record in '%s' at %02d:%02d", record.Project, record.End.Hour(), record.End.Minute())

			if !deleteRecord {
				return
			}

			out.Print("\n")
			err = t.DeleteRecord(record)
			if err != nil {
				out.Err("failed to delete record: %s", err)
				return
			}
			out.Success("Deleted record %s from '%s'", record.Start.Format(util.DateTimeFormat), record.Project)
		},
	}

	stop.Flags().BoolVarP(&deleteRecord, "delete", "D", false, "Delete the running record.")

	return stop
}
