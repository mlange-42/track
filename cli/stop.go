package cli

import (
	"fmt"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func stopCommand(t *core.Track) *cobra.Command {
	var deleteRecord bool
	var atTime string
	var ago time.Duration

	stop := &cobra.Command{
		Use:     "stop",
		Short:   "Stop the current record",
		Aliases: []string{"x"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			open, err := t.OpenRecord()
			if err != nil {
				return fmt.Errorf("failed to stop record: %s", err)
			}
			if open == nil {
				return fmt.Errorf("failed to stop record: no record running")
			}

			if deleteRecord && !confirm(
				fmt.Sprintf(
					"Really delete record %s (%s) from project '%s' (y/n): ",
					open.Start.Format(util.DateTimeFormat),
					util.FormatDuration(open.Duration(util.NoTime, util.NoTime)),
					open.Project,
				),
				"y",
			) {
				return fmt.Errorf("failed to stop record: aborted by user")
			}

			stopTime, err := getStopTime(open, ago, atTime)
			if err != nil {
				return fmt.Errorf("failed to stop record: %s", err)
			}

			record, err := t.StopRecord(stopTime)
			if err != nil {
				return fmt.Errorf("failed to stop record: %s", err)
			}
			out.Success("Stopped record in '%s' at %s", record.Project, record.End.Format(util.TimeFormat))

			if !deleteRecord {
				return nil
			}

			out.Print("\n")
			err = t.DeleteRecord(record)
			if err != nil {
				return fmt.Errorf("failed to delete record: %s", err)
			}
			out.Success("Deleted record %s from '%s'", record.Start.Format(util.DateTimeFormat), record.Project)
			return nil
		},
	}

	stop.Flags().BoolVarP(&deleteRecord, "delete", "D", false, "Delete the running record.")
	stop.Flags().StringVar(&atTime, "at", "", "Stop the record at a different time than now.")
	stop.Flags().DurationVar(&ago, "ago", 0*time.Second, "Stop the record at a different time than now, given as a duration.")

	stop.MarkFlagsMutuallyExclusive("at", "ago")

	return stop
}
