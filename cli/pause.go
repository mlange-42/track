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
	var atTime string
	var ago time.Duration

	pauseCom := &cobra.Command{
		Use:     "pause [NOTE...]",
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

			minTime := open.Start
			if pause, ok := open.LastPause(); ok {
				minTime = pause.End
			}

			var startTime, endTime time.Time
			var nowCorr = time.Now()
			var timeChanged = cmd.Flags().Changed("ago") || cmd.Flags().Changed("at")
			if timeChanged {
				nowCorr, err = getStartTime(minTime, ago, atTime)
				if err != nil {
					out.Err("failed to insert pause: %s", err)
					return
				}
			}
			if cmd.Flags().Changed("duration") {
				if timeChanged {
					startTime = nowCorr
					endTime = startTime.Add(duration)
				} else {
					endTime = nowCorr
					startTime = endTime.Add(-duration)
				}
				if endTime.After(time.Now()) {
					out.Err("failed to insert pause: end of pause would be in the future")
					return
				}
				if startTime.Before(minTime) {
					out.Err("can't start at a time before the last stop/pause")
					return
				}
			} else {
				startTime = nowCorr
				endTime = time.Time{}
			}

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
	pauseCom.Flags().DurationVarP(&duration, "duration", "d", 0*time.Hour, "Duration of the pause. Inserts a finished pause if given.\nOtherwise, a pause with an open end is inserted")

	pauseCom.Flags().StringVar(&atTime, "at", "", "Pause the record at a different time than now.\nRefers to the start time of the pause.")
	pauseCom.Flags().DurationVar(&ago, "ago", 0*time.Second, "Pause the record at a different time than now, given as a duration.\nRefers refers to the start time of the pause.")

	pauseCom.MarkFlagsMutuallyExclusive("at", "ago")

	return pauseCom
}
