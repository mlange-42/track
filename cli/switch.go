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

func switchCommand(t *core.Track) *cobra.Command {
	var copy bool
	var force bool
	var atTime string
	var ago time.Duration

	switchCom := &cobra.Command{
		Use:   "switch PROJECT [NOTE...]",
		Short: "Start a record and stop any running record",
		Long: fmt.Sprintf(`Start a record and stop any running record

Everything after the project name is considered a note for the record.
Notes can contain tags, denoted by the prefix "%s", like "%stag"`, core.TagPrefix, core.TagPrefix),
		Aliases: []string{"sw"},
		Args:    util.WrappedArgs(cobra.MinimumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			project := args[0]

			if !t.ProjectExists(project) {
				out.Err("failed to start record: project '%s' does not exist", project)
				return
			}

			if copy && len(args) > 1 {
				out.Err("failed to start record: can't use note arguments with flag --copy")
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

			var startStopTime time.Time
			open, err := t.OpenRecord()
			if err != nil {
				out.Err("failed to start record: %s", err)
				return
			}
			if open != nil {
				var err error
				startStopTime, err = getStopTime(open, ago, atTime)
				if err != nil {
					out.Err("failed to stop record: %s", err)
					return
				}

				record, err := t.StopRecord(startStopTime)
				if err != nil {
					out.Err("failed to create record: %s", err.Error())
					return
				}

				if !force && record.Project == project {
					out.Warn("Already working on project '%s'. Use --force to start a new record anyway", project)
					return
				}

				out.Success("Stopped record in '%s' at %s\n", record.Project, record.End.Format(util.TimeFormat))
			} else {
				latest, err := t.LatestRecord()
				if err != nil {
					out.Err("failed to create record: %s", err.Error())
					return
				}
				if latest != nil {
					startStopTime, err = getStartTime(latest.End, ago, atTime)
					if err != nil {
						out.Err("failed to create record: %s", err.Error())
						return
					}
				} else {
					startStopTime, err = getStartTime(util.NoTime, ago, atTime)
					if err != nil {
						out.Err("failed to create record: %s", err.Error())
						return
					}
				}
			}

			note := ""
			tags := map[string]string{}

			if copy {
				latest, err := t.FindLatestRecord(core.FilterByProjects([]string{project}))
				if err != nil {
					out.Err("failed to start record with copy: %s", err.Error())
					return
				}
				if latest != nil {
					note = latest.Note
					tags = latest.Tags
				} else {
					out.Err("failed to create record with copy: no previous record in '%s'", project)
					return
				}
			} else {
				note = strings.Join(args[1:], " ")
				tags, err = core.ExtractTagsSlice(args[1:])
				if err != nil {
					out.Err("failed to start record: %s", err.Error())
					return
				}
			}

			record, err := t.StartRecord(&proj, note, tags, startStopTime)
			if err != nil {
				out.Err("failed to create record: %s", err.Error())
				return
			}

			out.Success("Started record in '%s' at %s", project, record.Start.Format(util.TimeFormat))
		},
	}
	switchCom.Flags().BoolVarP(&copy, "copy", "c", false, "Copy note and tags from the last record of the project.")

	switchCom.Flags().BoolVarP(&force, "force", "f", false, "Force start of a new record if the project is already running")
	switchCom.Flags().StringVar(&atTime, "at", "", "Switch at a different time than now.")
	switchCom.Flags().DurationVar(&ago, "ago", 0*time.Second, "Switch at a different time than now, given as a duration.")

	switchCom.MarkFlagsMutuallyExclusive("at", "ago")

	return switchCom
}
