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
	var copy bool
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
		RunE: func(cmd *cobra.Command, args []string) error {
			project := args[0]

			if !t.ProjectExists(project) {
				return fmt.Errorf("failed to start record: project '%s' does not exist", project)
			}

			if copy && len(args) > 1 {
				return fmt.Errorf("failed to start record: can't use note arguments with flag --copy")
			}

			proj, err := t.LoadProjectByName(project)
			if err != nil {
				return fmt.Errorf("failed to start record: %s", err)
			}
			if proj.Archived {
				return fmt.Errorf("failed to start record: project '%s' is archived", proj.Name)
			}

			rec, err := t.OpenRecord()
			if err != nil {
				return fmt.Errorf("failed to start record: %s", err)
			}
			if rec != nil {
				return fmt.Errorf("failed to start record: record in '%s' still running", rec.Project)
			}

			var startTime time.Time

			latest, err := t.LatestRecord()
			if err != nil {
				return fmt.Errorf("failed to start record: %s", err.Error())
			}
			if latest != nil {
				startTime, err = getStartTime(latest.End, ago, atTime)
				if err != nil {
					return fmt.Errorf("failed to start record: %s", err.Error())
				}
			} else {
				startTime, err = getStartTime(util.NoTime, ago, atTime)
				if err != nil {
					return fmt.Errorf("failed to start record: %s", err.Error())
				}
			}

			note := ""
			tags := map[string]string{}

			if copy {
				latest, err := t.FindLatestRecord(core.FilterByProjects([]string{project}))
				if err != nil {
					return fmt.Errorf("failed to start record with copy: %s", err.Error())
				}
				if latest != nil {
					note = latest.Note
					tags = latest.Tags
				} else {
					return fmt.Errorf("failed to create record with copy: no previous record in '%s'", project)
				}
			} else {
				note = strings.Join(args[1:], " ")
				tags, err = core.ExtractTagsSlice(args[1:])
				if err != nil {
					return fmt.Errorf("failed to create record: %s", err.Error())
				}
			}

			record, err := t.StartRecord(&proj, note, tags, startTime)
			if err != nil {
				return fmt.Errorf("failed to create record: %s", err.Error())
			}

			out.Success("Started record in '%s' at %02d:%02d", project, record.Start.Hour(), record.Start.Minute())
			return nil
		},
	}

	start.Flags().BoolVarP(&copy, "copy", "c", false, "Copy note and tags from the last record of the project.")

	start.Flags().StringVar(&atTime, "at", "", "Start the record at a different time than now.")
	start.Flags().DurationVar(&ago, "ago", 0*time.Second, "Start the record at a different time than now, given as a duration.")

	start.MarkFlagsMutuallyExclusive("at", "ago")

	return start
}
