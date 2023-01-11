package cli

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var timelineModes = map[string]func(*core.Reporter, bool, bool) string{
	"days":   timelineDays,
	"weeks":  timelineWeeks,
	"months": timelineMonths,
	"d":      timelineDays,
	"w":      timelineWeeks,
	"m":      timelineMonths,
}

func timelineReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var csv bool
	var table bool

	timeline := &cobra.Command{
		Use:     "timeline (days|weeks|months)",
		Short:   "Timeline reports of time tracking",
		Aliases: []string{"l"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := args[0]

			if table && !csv {
				return fmt.Errorf("failed to generate report: flag --table can only be used together with --csv")
			}

			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err)
			}

			filters, err := createFilters(options, projects, false)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err)
			}

			timelineFunc, ok := timelineModes[mode]
			if !ok {
				return fmt.Errorf("failed to generate report: invalid timeline argument '%s'", mode)
			}

			startTime, endTime, err := parseStartEnd(options)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err)
			}
			reporter, err := core.NewReporter(
				t, options.projects, filters,
				options.includeArchived, startTime, endTime,
			)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err)
			}

			out.Print(timelineFunc(reporter, csv, table))
			return nil
		},
	}
	timeline.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	timeline.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	timeline.Flags().BoolVar(&csv, "csv", false, "Report in CSV format")
	timeline.Flags().BoolVar(&table, "table", false, "For report in CSV format, reports one column per project")

	return timeline
}

func timelineDays(r *core.Reporter, csv bool, table bool) string {
	startDate := util.ToDate(r.TimeRange.Start)
	if table {
		return timelineTable(r, startDate, time.Hour*24)
	}
	return timeline(r, startDate, time.Hour*24, 30*time.Minute, csv)
}

func timelineWeeks(r *core.Reporter, csv bool, table bool) string {
	startDate := util.ToDate(r.TimeRange.Start)
	weekDay := (int(startDate.Weekday()) + 6) % 7
	startDate = startDate.Add(time.Duration(-weekDay * 24 * int(time.Hour)))
	if table {
		return timelineTable(r, startDate, time.Hour*24*7)
	}
	return timeline(r, startDate, time.Hour*24*7, 2*time.Hour, csv)
}

func timelineMonths(r *core.Reporter, csv bool, table bool) string {
	y1, m1, _ := r.TimeRange.Start.Date()
	y2, m2, _ := r.TimeRange.End.Date()
	numBins := (y2-y1)*12 + int(m2) - int(m1) + 1

	dates := make([]time.Time, numBins)
	year, month := y1, m1
	for i := range dates {
		dates[i] = util.Date(year, month, 1)
		month++
		if month > 12 {
			year++
			month = 1
		}
	}

	values := make([]time.Duration, numBins)
	projectValues := make(map[string][]time.Duration)
	for p := range r.Projects {
		projectValues[p] = make([]time.Duration, numBins)
	}
	for _, rec := range r.Records {
		y2, m2, _ := rec.Start.Date()
		d := (y2-y1)*12 + int(m2) - int(m1)
		dur := rec.Duration(r.TimeRange.Start, r.TimeRange.End)
		values[d] += dur
		projectValues[rec.Project][d] += dur
	}
	if table {
		return renderTimelineTable(dates, values, projectValues)
	}
	if csv {
		return renderTimelineCsv(dates, values)
	}
	return renderTimeline(dates, values, 8*time.Hour)
}

func timeline(r *core.Reporter, startDate time.Time, delta time.Duration, perBox time.Duration, csv bool) string {
	minDate := startDate
	maxDate := util.ToDate(r.TimeRange.End.Add(delta))
	numBins := int(maxDate.Sub(minDate).Hours() / delta.Hours())

	dates := make([]time.Time, numBins)
	currDate := minDate
	for i := range dates {
		dates[i] = currDate
		currDate = currDate.Add(delta)
	}

	values := make([]time.Duration, numBins)
	for _, rec := range r.Records {
		// TODO: split if over increment
		d := int(rec.Start.Sub(minDate).Hours() / delta.Hours())
		values[d] = values[d] + rec.Duration(r.TimeRange.Start, r.TimeRange.End)
	}
	if csv {
		return renderTimelineCsv(dates, values)
	}
	return renderTimeline(dates, values, perBox)
}

func timelineTable(r *core.Reporter, startDate time.Time, delta time.Duration) string {
	minDate := startDate
	maxDate := util.ToDate(r.TimeRange.End.Add(delta))
	numBins := int(maxDate.Sub(minDate).Hours() / delta.Hours())

	dates := make([]time.Time, numBins)
	currDate := minDate
	for i := range dates {
		dates[i] = currDate
		currDate = currDate.Add(delta)
	}

	values := make([]time.Duration, numBins)
	projectValues := make(map[string][]time.Duration)
	for p := range r.Projects {
		projectValues[p] = make([]time.Duration, numBins)
	}

	for _, rec := range r.Records {
		// TODO: split if over increment
		d := int(rec.Start.Sub(minDate).Hours() / delta.Hours())
		dur := rec.Duration(r.TimeRange.Start, r.TimeRange.End)
		values[d] += dur
		projectValues[rec.Project][d] += dur
	}
	return renderTimelineTable(dates, values, projectValues)
}

func renderTimeline(dates []time.Time, values []time.Duration, perBox time.Duration) string {
	sb := strings.Builder{}
	for i := range dates {
		d := dates[i]
		v := values[i]
		fmt.Fprintf(&sb, "%s %s  %s  ", d.Weekday().String()[:2], d.Format(util.DateFormat), util.FormatDuration(v))

		boxes := float64(v) / float64(perBox)
		for i := 0; i < int(boxes); i++ {
			fmt.Fprint(&sb, "|")
		}
		if boxes > math.Floor(boxes)+0.5 {
			fmt.Fprint(&sb, ":")
		} else if boxes > math.Floor(boxes)+0.1 {
			fmt.Fprint(&sb, ".")
		}

		fmt.Fprintf(&sb, "\n")
	}

	return sb.String()
}

func renderTimelineCsv(dates []time.Time, values []time.Duration) string {
	sb := strings.Builder{}

	fmt.Fprintf(&sb, "date,weekday,duration\n")
	for i := range dates {
		d := dates[i]
		v := values[i]
		fmt.Fprintf(&sb, "%s,%s,%s\n", d.Format(util.DateFormat), d.Weekday().String()[:2], util.FormatDuration(v))
	}

	return sb.String()
}

func renderTimelineTable(dates []time.Time, values []time.Duration, projectValues map[string][]time.Duration) string {
	sb := strings.Builder{}

	projects := maps.Keys(projectValues)
	sort.Strings(projects)

	fmt.Fprintf(&sb, "date,weekday,total,%s\n", strings.Join(projects, ","))
	for i := range dates {
		d := dates[i]
		v := values[i]
		fmt.Fprintf(&sb, "%s,%s,%s", d.Format(util.DateFormat), d.Weekday().String()[:2], util.FormatDuration(v))
		for _, p := range projects {
			vp := projectValues[p][i]
			fmt.Fprintf(&sb, ",%s", util.FormatDuration(vp))
		}
		fmt.Fprintf(&sb, "\n")
	}

	return sb.String()
}
