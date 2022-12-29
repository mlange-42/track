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
	var atTime string
	var ago time.Duration

	start := &cobra.Command{
		Use:   "start PROJECT [NOTE...]",
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

			proj, err := t.LoadProjectByName(project)
			if err != nil {
				out.Err("failed to start record: %s", err)
				return
			}
			if proj.Archived {
				out.Err("failed to start record: project '%s' is archived", proj.Name)
				return
			}

			rec, err := t.OpenRecord()
			if err != nil {
				out.Err("failed to start record: %s", err)
				return
			}
			if rec != nil {
				out.Err("failed to start record: record in '%s' still running", rec.Project)
				return
			}

			var startTime time.Time

			latest, err := t.LatestRecord()
			if err != nil {
				out.Err("failed to start record: %s", err.Error())
				return
			}
			if latest != nil {
				startTime, err = getStartTime(latest, ago, atTime)
				if err != nil {
					out.Err("failed to start record: %s", err.Error())
					return
				}
			} else {
				startTime, err = getStartTime(nil, ago, atTime)
				if err != nil {
					out.Err("failed to start record: %s", err.Error())
					return
				}
			}

			note := strings.Join(args[1:], " ")
			tags := core.ExtractTagsSlice(args[1:])

			record, err := t.StartRecord(project, note, tags, startTime)
			if err != nil {
				out.Err("failed to create record: %s", err.Error())
				return
			}

			out.Success("Started record in '%s' at %02d:%02d", project, record.Start.Hour(), record.Start.Minute())
		},
	}

	start.Flags().StringVar(&atTime, "at", "", "Stop the record at a different time than now.")
	start.Flags().DurationVar(&ago, "ago", 0*time.Second, "Stop the record at a different time than now, given as a duration.")

	start.MarkFlagsMutuallyExclusive("at", "ago")

	return start
}
