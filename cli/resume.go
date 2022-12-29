package cli

import (
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func resumeCommand(t *core.Track) *cobra.Command {
	var continueIt bool

	resume := &cobra.Command{
		Use:   "resume [NOTE...]",
		Short: "Resume the last project",
		Long: `Resume the last project

If no note is provided, the note and tags from the last record are applied.

For details on notes and tags, see command 'start'.`,
		Aliases: []string{"re"},
		Args:    util.WrappedArgs(cobra.MinimumNArgs(0)),
		Run: func(cmd *cobra.Command, args []string) {
			last, err := t.LatestRecord()
			if err != nil {
				out.Err("failed to resume: %s", err)
				return
			}
			if !last.HasEnded() {
				out.Err("failed to resume: record running in '%s'", last.Project)
				return
			}

			project := last.Project

			if !t.ProjectExists(project) {
				out.Err("failed to resume: project '%s' does not exist", project)
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

			note := last.Note
			tags := last.Tags
			if len(args) > 0 {
				note = strings.Join(args[1:], " ")
				tags = core.ExtractTagsSlice(args[1:])
			}

			if continueIt {
				oldEnd := last.End
				last.Note = note
				last.Tags = tags
				last.End = time.Time{}

				err = t.SaveRecord(last, true)
				if err != nil {
					out.Err("failed to resume: %s", err)
					return
				}
				out.Success(
					"Continue record in '%s' at %s - skipping %s break.",
					project,
					time.Now().Format(util.TimeFormat),
					util.FormatDuration(time.Now().Sub(oldEnd)),
				)
				return
			}

			record, err := t.StartRecord(project, note, tags, time.Now())
			if err != nil {
				out.Err("failed to resume: %s", err.Error())
				return
			}

			out.Success("Resume record in '%s' at %s", project, record.Start.Format(util.TimeFormat))
		},
	}

	resume.Flags().BoolVarP(&continueIt, "continue", "C", false, "Continue the last record instead of starting a new one. Skips the breaks.")

	return resume
}
