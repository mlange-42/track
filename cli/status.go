package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func statusCommand(t *core.Track) *cobra.Command {
	var maxBreakStr string

	status := &cobra.Command{
		Use:   "status [PROJECT]",
		Short: "Reports the status of the running or given project",
		Long: `Reports the status of the running or given project

Columns of the status are:

* curr  - Time of the current record
* total - Recorded time since the last break longer than --max-break
* break - Break time since the last break longer than --max-break
* today - Total recorded time since midnight
`,
		Aliases: []string{"s"},
		Args:    util.WrappedArgs(cobra.MaximumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			maxBreak, err := time.ParseDuration(maxBreakStr)
			if err != nil {
				out.Err("failed to show status: %s", err)
				return
			}

			var project string
			open, hasOpenRecord := t.OpenRecord()
			if len(args) > 0 {
				project = args[0]
				if hasOpenRecord && open.Project != project {
					hasOpenRecord = false
				}
			} else {
				if !hasOpenRecord {
					last, err := t.LatestRecord()
					out.Warn("No running record. Start tracking or specify a project.")
					if err != nil {
						return
					}
					out.Print("\n")
					out.Warn(
						"Stopped project '%s' %s ago\n",
						last.Project,
						util.FormatDuration(time.Now().Sub(last.End)),
					)
					project = last.Project
				} else {
					project = open.Project
				}
			}

			now := time.Now()
			start := util.ToDate(now)
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
			cumTime := time.Second * 0
			breakTime := time.Second * 0
			totalTime := time.Second * 0

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

				if !prevEnd.IsZero() {
					bt := startTime.Sub(prevEnd)
					if bt < maxBreak {
						breakTime += bt
					} else {
						cumTime = time.Second * 0
						breakTime = time.Second * 0
					}
				}

				totalTime += worked
				cumTime += worked
				if rec.End.IsZero() {
					currTime += worked
				}

				prevEnd = endTime
			}
			out.Print("+------------------+-------+-------+-------+-------+\n")
			out.Print("|          project |  curr | total | break | today |\n")
			out.Print(
				"| %16s | %s | %s | %s | %s |\n",
				project,
				util.FormatDuration(currTime),
				util.FormatDuration(cumTime),
				util.FormatDuration(breakTime),
				util.FormatDuration(totalTime),
			)
			out.Print("+------------------+-------+-------+-------+-------+")
		},
	}
	status.Flags().StringVar(
		&maxBreakStr,
		"max-break",
		t.Config.MaxBreakDuration.String(),
		"Maximum length of breaks to consider them in daily break time.\nThe default can be set in the config file",
	)

	return status
}
