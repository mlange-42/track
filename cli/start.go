package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func startCommand(t *core.Track) *cobra.Command {
	start := &cobra.Command{
		Use:   "start <PROJECT> [note...]",
		Short: "Start a record for a project",
		Long: fmt.Sprintf(`Start a record for a project
		
Everything after the project name is considered a note for the record.
Notes can contain tags, denoted by the prefix "%s", like "%stag"`, core.TagPrefix, core.TagPrefix),
		Aliases: []string{"+"},
		Args:    util.WrappedArgs(cobra.MinimumNArgs(1)),
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
