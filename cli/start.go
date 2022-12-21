package cli

import (
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/spf13/cobra"
)

func startCommand(t *core.Track) *cobra.Command {
	start := &cobra.Command{
		Use:   "start <project> [message]",
		Short: "Start a record",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			project := args[0]

			if !t.ProjectExists(project) {
				out.Err("failed to start record: project '%s' does not exist", project)
				return
			}

			if rec, ok := t.OpenRecord(); ok {
				out.Err("failed to start record: record in '%s' still running", rec.Project)
				return
			}

			note := strings.Join(args[1:], " ")

			record := core.Record{
				Project: project,
				Note:    note,
				Start:   time.Now(),
				End:     time.Time{},
			}

			if err := t.SaveRecord(record, false); err != nil {
				out.Err("failed to create record: %s", err.Error())
				return
			}

			out.Success("Started record at %02d:%02d", record.Start.Hour(), record.Start.Minute())
		},
	}

	return start
}
