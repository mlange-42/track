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

func timelineReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	timeline := &cobra.Command{
		Use:     "timeline (days|weeks|months)",
		Short:   "Timeline reports of time tracking",
		Aliases: []string{"t"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			mode := args[0]

			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			filters, err := createFilters(options, projects, false)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			timelineFunc, ok := timelineModes[mode]
			if !ok {
				out.Err("failed to generate report: invalid timeline argument '%s'", mode)
				return
			}

			startTime, endTime, err := parseStartEnd(options)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			reporter, err := core.NewReporter(
				t, options.projects, filters,
				options.includeArchived, startTime, endTime,
			)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			out.Print(timelineFunc(reporter))
		},
	}
	timeline.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	timeline.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	return timeline
}

func timelineDays(r *core.Reporter) string {
	startDate := util.ToDate(r.TimeRange.Start)
	return timeline(r, startDate, time.Hour*24, time.Minute*30)
}

func timelineWeeks(r *core.Reporter) string {
	startDate := util.ToDate(r.TimeRange.Start)
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
		dates[i] = util.Date(year, month, 1)
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
		values[d] = values[d] + rec.Duration(r.TimeRange.Start, r.TimeRange.End).Hours()
	}

	return renderTimeline(dates, values, 8)
}
func timeline(r *core.Reporter, startDate time.Time, delta time.Duration, unit time.Duration) string {
	minDate := startDate
	maxDate := util.ToDate(r.TimeRange.End.Add(delta))
	numBins := int(maxDate.Sub(minDate).Hours() / delta.Hours())

	dates := make([]time.Time, numBins, numBins)
	currDate := minDate
	for i := range dates {
		dates[i] = currDate
		currDate = currDate.Add(delta)
	}

	values := make([]float64, numBins, numBins)
	for _, rec := range r.Records {
		// TODO: split if over increment
		d := int(rec.Start.Sub(minDate).Hours() / delta.Hours())
		values[d] = values[d] + rec.Duration(r.TimeRange.Start, r.TimeRange.End).Hours()
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
