package cli

import (
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func pauseCommand(t *core.Track) *cobra.Command {
	var duration time.Duration

	pauseCom := &cobra.Command{
		Use:     "pause [NOTE]",
		Short:   "Pauses or inserts a pause into the running recording",
		Long:    `Pauses or inserts a pause into the running recording`,
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ArbitraryArgs),
		Run: func(cmd *cobra.Command, args []string) {
			open, err := t.OpenRecord()
			if err != nil {
				out.Err("failed to insert pause: %s", err)
				return
			}
			if open == nil {
				out.Warn("failed to insert pause: no running record")
				return
			}
			if open.IsPaused() {
				out.Err("failed to insert pause: record is already paused")
				return
			}

			now := time.Now()

			endTime := time.Time{}
			if cmd.Flags().Changed("duration") {
				endTime = now
			}
			startTime := now.Add(-duration)
			note := strings.Join(args, " ")
			_, err = open.InsertPause(startTime, endTime, note)
			if err != nil {
				out.Err("failed to insert pause: %s", err)
				return
			}

			err = t.SaveRecord(open, true)
			if err != nil {
				out.Err("failed to pause record: %s", err)
				return
			}
			if endTime.IsZero() {
				out.Success("Paused record in '%s'\n", open.Project)
			} else {
				out.Success("Inserted pause of %s in '%s'\n", duration, open.Project)
			}
		},
	}
	pauseCom.Flags().DurationVarP(&duration, "duration", "d", 0*time.Hour, "Duration of the break. Inserts a finished break if given")

	return pauseCom
}
