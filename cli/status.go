package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func statusCommand(t *core.Track) *cobra.Command {
	status := &cobra.Command{
		Use:     "status [PROJECT]",
		Short:   "Reports the status of the running or given project",
		Aliases: []string{"s"},
		Args:    util.WrappedArgs(cobra.MaximumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			var project string
			open, hasOpenRecord := t.OpenRecord()
			if len(args) > 0 {
				project = args[0]
				if hasOpenRecord && open.Project != project {
					hasOpenRecord = false
				}
			} else {
				if !hasOpenRecord {
					out.Err("No running record. Start tracking or specify a project.")
					return
				}
				project = open.Project
			}

			now := time.Now()
			start := util.Date(now.Date())
			filterStart := start.Add(-time.Hour * 24)

			filters := core.FilterFunctions{
				core.FilterByTime(filterStart, time.Time{}),
			}

			reporter, err := core.NewReporter(t, []string{project}, filters)
			if err != nil {
				out.Err("failed to show status: %s", err)
				return
			}

			prevEnd := time.Time{}
			currTime := time.Second * 0
			totalTime := time.Second * 0
			breakTime := time.Second * 0

			for _, rec := range reporter.Records {
				endTime := rec.End
				if endTime.IsZero() {
					endTime = now
				} else {
					if endTime.Before(start) {
						continue
					}
				}
				startTime := rec.Start
				if startTime.Before(start) {
					startTime = start
				}

				worked := endTime.Sub(startTime)

				totalTime += worked
				if rec.End.IsZero() {
					currTime += worked
				}

				if !prevEnd.IsZero() {
					breakTime += startTime.Sub(prevEnd)
				}

				prevEnd = endTime
			}
			out.Print("+------------------+-------+-------+-------+\n")
			out.Print("|          project |  curr | today | break |\n")
			out.Print(
				"| %16s | %s | %s | %s |\n",
				project,
				util.FormatDuration(currTime),
				util.FormatDuration(totalTime),
				util.FormatDuration(breakTime),
			)
			out.Print("+------------------+-------+-------+-------+")
		},
	}

	return status
}
