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

var timelineModes = map[string]func(*core.Reporter) string{
	"days":   timelineDays,
	"weeks":  timelineWeeks,
	"months": timelineMonths,
	"d":      timelineDays,
	"w":      timelineWeeks,
	"m":      timelineMonths,
}

func reportCommand(t *core.Track) *cobra.Command {
	var projects []string
	var tags []string
	var start string
	var end string
	var timeline string

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

			var timelineFunc func(*core.Reporter) string
			if timeline != "" {
				var ok bool
				if timelineFunc, ok = timelineModes[timeline]; !ok {
					out.Err("failed to generate report: invalid timeline argument '%s'", timeline)
					return
				}
			}

			reporter, err := core.NewReporter(t, projects, filters)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			if timeline == "" {
				for name, dur := range reporter.ProjectTime {
					out.Success("%-15s %s\n", name, util.FormatDuration(dur))
				}
				return
			}
			out.Success(timelineFunc(reporter))
		},
	}

	report.Flags().StringSliceVarP(&projects, "projects", "p", []string{}, "Projects to include. Includes all projects if not specified")
	report.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Tags to include. Includes records with any of the given tags")
	report.Flags().StringVarP(&start, "start", "s", "", "Start date")
	report.Flags().StringVarP(&end, "end", "e", "", "End date")
	report.Flags().StringVarP(&timeline, "timeline", "l", "", "Timeline mode. One of [days weeks months] (first letter is sufficient).")

	return report
}

func timelineDays(r *core.Reporter) string {
	minDate := util.Date(r.TimeRange.Start.Date())
	maxDate := util.Date(r.TimeRange.End.Add(time.Hour * 24).Date())
	numDays := int(maxDate.Sub(minDate).Hours() / 24)

	dates := make([]time.Time, numDays, numDays)
	currDate := minDate
	for i := range dates {
		dates[i] = currDate
		currDate = currDate.Add(time.Hour * 24)
	}

	values := make([]float64, numDays, numDays)
	for _, rec := range r.Records {
		d := int(rec.Start.Sub(minDate).Hours() / 24)
		values[d] = values[d] + rec.Duration().Hours()
	}

	return renderTimeline(dates, values)
}

func timelineWeeks(r *core.Reporter) string {
	return ""
}

func timelineMonths(r *core.Reporter) string {
	return ""
}

func renderTimeline(dates []time.Time, values []float64) string {
	sb := strings.Builder{}
	for i := range dates {
		d := dates[i]
		v := values[i]
		fmt.Fprintf(&sb, "%s %s ", d.Weekday().String()[:2], d.Format(util.DateFormat))
		for i := 0; i < int(v); i++ {
			fmt.Fprint(&sb, "#")
		}
		fmt.Fprintf(&sb, "\n")
	}

	return sb.String()
}
