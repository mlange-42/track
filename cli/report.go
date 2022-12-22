package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func reportCommand(t *core.Track) *cobra.Command {
	var projects []string
	var tags []string
	var start string
	var end string

	report := &cobra.Command{
		Use:   "report",
		Short: "Generate reports of time tracking",
		Run: func(cmd *cobra.Command, args []string) {
			var filters core.FilterFunctions

			if len(projects) > 0 {
				filters = append(filters, core.FilterByProjects(projects))
			}
			if len(tags) > 0 {
				filters = append(filters, core.FilterByTagsAny(tags))
			}
			var err error
			var startTime time.Time
			var endTime time.Time
			if len(start) > 0 {
				startTime, err = util.ParseDate(start)
				if err != nil {
					out.Err("failed to generate report: %s", err)
					return
				}
			}
			if len(end) > 0 {
				endTime, err = util.ParseDate(end)
				if err != nil {
					out.Err("failed to generate report: %s", err)
					return
				}
			}
			if !(startTime.IsZero() && endTime.IsZero()) {
				filters = append(filters, core.FilterByTime(startTime, endTime))
			}

			reporter, err := core.NewReporter(t, projects, filters)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			for name, dur := range reporter.ProjectTime {
				out.Success("%-15s %s\n", name, util.FormatDuration(dur))
			}
		},
	}

	report.Flags().StringSliceVarP(&projects, "projects", "p", []string{}, "Projects to include. Includes all projects if not specified")
	report.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Tags to include. Includes records with any of the given tags")
	report.Flags().StringVarP(&start, "start", "s", "", "Start date")
	report.Flags().StringVarP(&end, "end", "e", "", "End date")

	return report
}
