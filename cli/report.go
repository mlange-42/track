package cli

import (
	"fmt"
	"math"
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
	startDate := util.Date(r.TimeRange.Start.Date())
	return timeline(r, startDate, time.Hour*24, time.Minute*30)
}

func timelineWeeks(r *core.Reporter) string {
	startDate := util.Date(r.TimeRange.Start.Date())
	weekDay := (int(startDate.Weekday()) + 6) % 7
	startDate = startDate.Add(time.Duration(-weekDay * 24 * int(time.Hour)))
	return timeline(r, startDate, time.Hour*24*7, time.Hour)
}

func timelineMonths(r *core.Reporter) string {
	y1, m1, _ := r.TimeRange.Start.Date()
	y2, m2, _ := r.TimeRange.End.Date()
	numBins := (y2-y1)*12 + int(m2) - int(m1) + 1

	dates := make([]time.Time, numBins, numBins)
	year, month := y1, m1
	for i := range dates {
		dates[i] = util.Date(year, time.Month(month), 1)
		month++
		if month > 12 {
			year++
			month = 1
		}
	}

	values := make([]float64, numBins, numBins)
	for _, rec := range r.Records {
		y2, m2, _ := rec.Start.Date()
		d := (y2-y1)*12 + int(m2) - int(m1)
		values[d] = values[d] + rec.Duration().Hours()
	}

	return renderTimeline(dates, values, 8)
}
func timeline(r *core.Reporter, startDate time.Time, delta time.Duration, unit time.Duration) string {
	minDate := startDate
	maxDate := util.Date(r.TimeRange.End.Add(delta).Date())
	numBins := int(maxDate.Sub(minDate).Hours() / delta.Hours())

	dates := make([]time.Time, numBins, numBins)
	currDate := minDate
	for i := range dates {
		dates[i] = currDate
		currDate = currDate.Add(delta)
	}

	values := make([]float64, numBins, numBins)
	for _, rec := range r.Records {
		d := int(rec.Start.Sub(minDate).Hours() / delta.Hours())
		values[d] = values[d] + rec.Duration().Hours()
	}

	return renderTimeline(dates, values, unit.Hours())
}

func renderTimeline(dates []time.Time, values []float64, unit float64) string {
	sb := strings.Builder{}
	for i := range dates {
		d := dates[i]
		v := values[i] / unit
		fmt.Fprintf(&sb, "%s %s ", d.Weekday().String()[:2], d.Format(util.DateFormat))
		for i := 0; i < int(v); i++ {
			fmt.Fprint(&sb, "#")
		}
		if v > math.Floor(v)+0.1 {
			fmt.Fprint(&sb, ".")
		}
		fmt.Fprintf(&sb, "\n")
	}

	return sb.String()
}
