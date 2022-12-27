package cli

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gookit/color"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
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
	report.AddCommand(weekReportCommand(t, &options))

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

			tree, err := t.ToProjectTree(reporter.Projects)
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
					return fmt.Sprintf(
						"%s %s (%s)", str,
						util.FormatDuration(reporter.TotalTime[t.Value.Name]),
						util.FormatDuration(reporter.ProjectTime[t.Value.Name]),
					)
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

			str, err := renderDayTimeline(t, reporter, active, start, blocksPerHour)
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

func weekReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var blocksPerHour int

	day := &cobra.Command{
		Use:     "week [DATE]",
		Short:   "Report of activities over a week",
		Aliases: []string{"w"},
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
			weekDay := (int(start.Weekday()) + 6) % 7
			start = start.Add(time.Duration(-weekDay * 24 * int(time.Hour)))

			if blocksPerHour <= 0 {
				out.Err("failed to generate report: argument --width must be > 0")
				return
			}
			filterStart := start.Add(-time.Hour * 24)
			filterEnd := start.Add(time.Hour * 24 * 7)

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

			str, err := renderWeekTimeline(t, reporter, active, start, blocksPerHour)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			fmt.Print(str)
		},
	}

	day.Flags().IntVarP(&blocksPerHour, "width", "w", 12, "Width of the graph, in characters per hour")

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

func renderDayTimeline(t *core.Track, reporter *core.Reporter, active string, startDate time.Time, blocksPerHour int) (string, error) {
	bph := blocksPerHour

	tree, err := t.ToProjectTree(reporter.Projects)
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
		if endIdx >= bph*24 {
			endIdx = bph*24 - 1
		}

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

	interval := 6
	if bph > 1 {
		interval = 3
	}

	timelineStr := map[string]string{}
	for pr, values := range timelines {
		runes := make([]rune, bph*24, bph*24)
		for i, v := range values {
			runes[i] = util.FloatToBlock(v)
		}
		timelineStr[pr] = toDayTimeline(runes, interval*bph)
	}
	timelineStr[t.WorkspaceLabel()] = toDayAxis(bph, interval*bph)

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

func toDayTimeline(runes []rune, interval int) string {
	sb := strings.Builder{}

	for i, r := range runes {
		if i%interval == 0 {
			fmt.Fprint(&sb, "|")
		}
		fmt.Fprint(&sb, string(r))
	}

	fmt.Fprint(&sb, "|")

	return sb.String()
}

func toDayAxis(blocksPerHour int, interval int) string {
	fill := strings.Repeat(" ", interval-5)

	sb := strings.Builder{}

	for i := 0; i < 24*blocksPerHour; i += interval {
		fmt.Fprintf(&sb, "|%02d:00%s", i/blocksPerHour, fill)
	}
	fmt.Fprint(&sb, "|")

	return sb.String()
}

func renderWeekTimeline(t *core.Track, reporter *core.Reporter, active string, startDate time.Time, blocksPerHour int) (string, error) {
	bph := blocksPerHour

	projects := maps.Keys(reporter.Projects)
	sort.Strings(projects)
	indices := make(map[string]int, len(projects))
	symbols := make([]rune, len(projects)+1, len(projects)+1)
	colors := make([]uint8, len(projects)+1, len(projects)+1)
	symbols[0] = '·'
	colors[0] = 0
	for i, p := range projects {
		indices[p] = i + 1
		symbols[i+1] = []rune(p)[0]
		colors[i+1] = reporter.Projects[p].Color
	}

	timeline := make([]int, 24*7*blocksPerHour, 24*7*blocksPerHour)

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

		startIdx := int(start)
		endIdx := int(end)
		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx >= bph*24*7 {
			endIdx = bph*24*7 - 1
		}

		index := indices[rec.Project]

		for i := startIdx; i <= endIdx; i++ {
			timeline[i] = index
		}
	}

	timeStr := make([]rune, len(timeline), len(timeline))
	for i, idx := range timeline {
		timeStr[i] = symbols[idx]
	}

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "      |Week %s - %s\n", startDate.Format(util.DateFormat), startDate.Add(6*24*time.Hour).Format(util.DateFormat))

	fmt.Fprint(&sb, "      ")
	for weekday := 0; weekday < 7; weekday++ {
		date := startDate.Add(time.Duration(weekday * 24 * int(time.Hour)))
		str := fmt.Sprintf(
			"%s %02d %s",
			time.Weekday((weekday + 1) % 7).String()[:2],
			date.Day(),
			date.Month().String()[:3],
		)
		if len(str) > bph {
			fmt.Fprintf(&sb, "|%s", str[:bph])
		} else {
			fmt.Fprintf(&sb, "|%s%s", str, strings.Repeat(" ", bph-len(str)))
		}
	}
	fmt.Fprintln(&sb, "|")

	for hour := 0; hour < 24; hour++ {
		fmt.Fprintf(&sb, "%02d:00 ", hour)
		for weekday := 0; weekday < 7; weekday++ {
			s := (weekday*24 + hour) * bph
			fmt.Fprint(&sb, "|")
			for i := s; i < s+bph; i++ {
				pr := timeline[i]
				sym := symbols[pr]
				col := colors[pr]
				if col == 0 {
					fmt.Fprintf(&sb, "%c", sym)
				} else {
					fmt.Fprint(&sb, color.C256(col, true).Sprintf("%c", sym))
				}
			}
		}
		fmt.Fprintln(&sb, "|")
	}

	totalWidth := 7 + 7*(bph+1)
	lineWidth := 0
	for i, p := range projects {
		col := colors[i+1]
		width := utf8.RuneCountInString(p)
		if lineWidth > 0 && lineWidth+width+2 > totalWidth {
			lineWidth = 0
			fmt.Fprintln(&sb)
		}

		fmt.Fprint(&sb, color.C256(col, true).Sprintf(" %s ", p))
		lineWidth += width + 2
	}

	return sb.String(), nil
}
