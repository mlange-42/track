package cli

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gookit/color"
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
	options := filterOptions{}

	report := &cobra.Command{
		Use:     "report",
		Short:   "Generate reports of time tracking",
		Long:    "Generate reports of time tracking",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	report.PersistentFlags().StringSliceVarP(&options.projects, "projects", "p", []string{}, "Projects to include (comma-separated). All projects if not specified")
	report.PersistentFlags().StringSliceVarP(&options.tags, "tags", "t", []string{}, "Tags to include (comma-separated). Includes records with any of the given tags")

	report.AddCommand(timelineReportCommand(t, &options))
	report.AddCommand(projectsReportCommand(t, &options))
	report.AddCommand(dayReportCommand(t, &options))

	report.Long += "\n\n" + formatCmdTree(report)
	return report
}

func timelineReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	timeline := &cobra.Command{
		Use:     "timeline (days|weeks|months)",
		Short:   "Timeline reports of time tracking",
		Aliases: []string{"t"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			mode := args[0]

			filters, err := createFilters(options, false)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			timelineFunc, ok := timelineModes[mode]
			if !ok {
				out.Err("failed to generate report: invalid timeline argument '%s'", mode)
				return
			}

			reporter, err := core.NewReporter(t, options.projects, filters)
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

func projectsReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	projects := &cobra.Command{
		Use:     "projects",
		Short:   "Timeline reports of time tracking",
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			filters, err := createFilters(options, false)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			reporter, err := core.NewReporter(t, options.projects, filters)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			tree, err := core.ToProjectTree(reporter.Projects)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			var active string
			if rec, ok := t.OpenRecord(); ok {
				active = rec.Project
			}
			formatter := util.NewTreeFormatter(
				func(t *core.ProjectNode, indent int) string {
					fillLen := 16 - (indent + utf8.RuneCountInString(t.Value.Name))
					var str string
					if t.Value.Name == active {
						str = color.BgBlue.Sprintf("%s", t.Value.Name)
					} else {
						str = fmt.Sprintf("%s", t.Value.Name)
					}
					if fillLen > 0 {
						str += strings.Repeat(" ", fillLen)
					}
					return fmt.Sprintf("%s %s", str, util.FormatDuration(reporter.ProjectTime[t.Value.Name]))
				},
				2,
			)
			fmt.Print(formatter.FormatTree(tree))
		},
	}
	projects.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	projects.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	return projects
}

func dayReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var blocksPerHour int

	day := &cobra.Command{
		Use:     "day [DATE]",
		Short:   "Report of activities over the day",
		Aliases: []string{"d"},
		Args:    util.WrappedArgs(cobra.MaximumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			start := util.ToDate(time.Now())
			var err error
			if len(args) > 0 {
				start, err = util.ParseDate(args[0])
				if err != nil {
					out.Err("failed to generate report: %s", err)
					return
				}
			}
			if blocksPerHour <= 0 {
				out.Err("failed to generate report: argument --width must be > 0")
				return
			}
			filterStart := start.Add(-time.Hour * 24)
			filterEnd := start.Add(time.Hour * 24)

			filters, err := createFilters(options, false)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			filters = append(filters, core.FilterByTime(filterStart, filterEnd))

			reporter, err := core.NewReporter(t, options.projects, filters)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			var active string
			if rec, ok := t.OpenRecord(); ok {
				active = rec.Project
			}

			str, err := renderDayTimeline(reporter, active, start, blocksPerHour)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			fmt.Print(str)
		},
	}

	day.Flags().IntVarP(&blocksPerHour, "width", "w", 3, "Width of the graph, in characters per hour")

	return day
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
		values[d] = values[d] + rec.Duration().Hours()
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

func renderDayTimeline(reporter *core.Reporter, active string, startDate time.Time, blocksPerHour int) (string, error) {
	bph := blocksPerHour

	tree, err := core.ToProjectTree(reporter.Projects)
	if err != nil {
		return "", err
	}

	timelines := map[string][]float64{}
	for pr := range reporter.Projects {
		timelines[pr] = make([]float64, bph*24, bph*24)
	}

	now := time.Now()

	for _, rec := range reporter.Records {
		endTime := rec.End
		if rec.End.IsZero() {
			endTime = now
		}
		if endTime.Before(startDate) {
			continue
		}
		start := rec.Start.Sub(startDate).Hours() * float64(bph)
		end := endTime.Sub(startDate).Hours() * float64(bph)

		if start < 0 {
			start = 0
		}
		if end > float64(bph*24) {
			end = float64(bph * 24)
		}
		startIdx := int(start)
		endIdx := int(end)

		timeline := timelines[rec.Project]
		for i := startIdx; i <= endIdx; i++ {
			startProp := start - float64(i)
			if startProp < 0 {
				startProp = 0.0
			}
			endProp := end - float64(i)
			if endProp > 1 {
				endProp = 1.0
			}
			timeline[i] += endProp - startProp
		}
	}

	timelineStr := map[string]string{}
	for pr, values := range timelines {
		runes := make([]rune, bph*24, bph*24)
		for i, v := range values {
			runes[i] = util.FloatToBlock(v)
		}
		timelineStr[pr] = fmt.Sprintf(
			"|%s|%s|%s|%s|",
			string(runes[0:6*bph]),
			string(runes[6*bph:12*bph]),
			string(runes[12*bph:18*bph]),
			string(runes[18*bph:24*bph]),
		)
	}
	fill := strings.Repeat(" ", 6*bph-5)
	timelineStr[core.RootName] = fmt.Sprintf(
		"|%02d:00%s|%02d:00%s|%02d:00%s|%02d:00%s|%02d:00",
		0, fill,
		6, fill,
		12, fill,
		18, fill,
		24,
	)

	formatter := util.NewTreeFormatter(
		func(t *core.ProjectNode, indent int) string {
			fillLen := 16 - (indent + utf8.RuneCountInString(t.Value.Name))
			var str string
			if t.Value.Name == active {
				str = color.BgBlue.Sprintf("%s", t.Value.Name)
			} else {
				str = fmt.Sprintf("%s", t.Value.Name)
			}
			if fillLen > 0 {
				str += strings.Repeat(" ", fillLen)
			}
			text, _ := timelineStr[t.Value.Name]
			return fmt.Sprintf("%s %s", str, text)
		},
		2,
	)
	return formatter.FormatTree(tree), nil
}
