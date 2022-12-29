package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func breakCommand(t *core.Track) *cobra.Command {
	breakCom := &cobra.Command{
		Use:     "break DURATION",
		Short:   "Inserts a break into the running recording",
		Long:    `Inserts a break into the running recording`,
		Aliases: []string{"b"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			open, err := t.OpenRecord()
			if err != nil {
				out.Err("failed to insert break: %s", err)
				return
			}
			if open == nil {
				out.Warn("failed to insert break: no running record")
				return
			}

			duration, err := time.ParseDuration(args[0])
			if err != nil {
				out.Err("failed to insert break: %s", err)
				return
			}
			endTime := time.Now().Add(-duration)
			if endTime.Before(open.Start) {
				out.Err("failed to insert break: break is longer than current record")
				return
			}

			open.End = endTime
			err = t.SaveRecord(open, true)
			if err != nil {
				out.Err("failed to stop record: %s", err)
				return
			}
			out.Success("Stopped record in '%s' at %s\n", open.Project, open.End.Format(util.TimeFormat))

			record, err := t.StartRecord(open.Project, open.Note, open.Tags, time.Now())
			if err != nil {
				out.Err("failed to create record: %s", err.Error())
				return
			}

			out.Success("Started record in '%s' at %s", record.Project, record.Start.Format(util.TimeFormat))
		},
	}

	return breakCom
}
