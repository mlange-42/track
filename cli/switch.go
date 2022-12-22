package cli

import (
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/spf13/cobra"
)

func switchCommand(t *core.Track) *cobra.Command {
	var force bool

	switchCom := &cobra.Command{
		Use:   "switch <project> [message]",
		Short: "Start a record and stop any running record",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			project := args[0]

			if !t.ProjectExists(project) {
				out.Err("failed to start record: project '%s' does not exist", project)
				return
			}

			if _, ok := t.OpenRecord(); ok {
				record, err := t.StopRecord(time.Now())
				if err != nil {
					out.Err("failed to create record: %s", err.Error())
					return
				}

				if !force && record.Project == project {
					out.Warn("Already working on project '%s'. Use --force to start a new record anyway\n", project)
					return
				}

				out.Success("Stopped record in '%s' at %02d:%02d\n", record.Project, record.End.Hour(), record.End.Minute())
			}

			note := strings.Join(args[1:], " ")
			tags := t.ExtractTags(args[1:])

			record, err := t.StartRecord(project, note, tags, time.Now())
			if err != nil {
				out.Err("failed to create record: %s", err.Error())
				return
			}

			out.Success("Started record in '%s' at %02d:%02d\n", project, record.Start.Hour(), record.Start.Minute())
		},
	}

	switchCom.Flags().BoolVarP(&force, "force", "f", false, "Force start of a new record if the project is already running")

	return switchCom
}
