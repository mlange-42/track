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
		Use:     "start <project> [message]",
		Short:   "Start a record",
		Aliases: []string{"+"},
		Args:    cobra.MinimumNArgs(1),
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
			tags := t.ExtractTags(args[1:])

			record, err := t.StartRecord(project, note, tags, time.Now())
			if err != nil {
				out.Err("failed to create record: %s", err.Error())
				return
			}

			out.Success("Started record in '%s' at %02d:%02d", project, record.Start.Hour(), record.Start.Minute())
		},
	}

	return start
}
