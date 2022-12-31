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
	report.PersistentFlags().BoolVarP(&options.includeArchived, "archived", "a", false, "Include records from archived projects")

	report.AddCommand(timelineReportCommand(t, &options))
	report.AddCommand(projectsReportCommand(t, &options))
	report.AddCommand(chartReportCommand(t, &options))
	report.AddCommand(weekReportCommand(t, &options))
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

func projectsReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	projects := &cobra.Command{
		Use:     "projects",
		Short:   "Timeline reports of time tracking",
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
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

			tree, err := t.ToProjectTree(reporter.Projects)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			var active string
			rec, err := t.OpenRecord()
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			if rec != nil {
				active = rec.Project
			}
			formatter := util.NewTreeFormatter(
				func(t *core.ProjectNode, indent int) string {
					fillLen := 16 - (indent + utf8.RuneCountInString(t.Value.Name))
					name := t.Value.Name
					if fillLen < 0 {
						nameRunes := []rune(name)
						name = string(nameRunes[:len(nameRunes)+fillLen-1]) + "."
					}
					var str string
					if t.Value.Name == active {
						str = color.BgBlue.Sprintf("%s", name)
					} else {
						str = fmt.Sprintf("%s", name)
					}
					if fillLen > 0 {
						str += strings.Repeat(" ", fillLen)
					}
					str += " "
					str += t.Value.Render.Sprintf(" %s ", t.Value.Symbol)

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

func chartReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var blocksPerHour int

	day := &cobra.Command{
		Use:     "chart [DATE]",
		Short:   "Report of activities over the day as a bar chart per project",
		Aliases: []string{"c"},
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
			if !cmd.Flags().Changed("width") {
				if w, _, err := util.TerminalSize(); err == nil && w > 0 {
					blocksPerHour = (w - 29) / 24
				}
			}

			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			filterStart := start.Add(-time.Hour * 24)
			filterEnd := start.Add(time.Hour * 24)

			filters, err := createFilters(options, projects, false)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			filters = append(filters, core.FilterByTime(filterStart, filterEnd))

			reporter, err := core.NewReporter(t, options.projects, filters, options.includeArchived, start, filterEnd)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			var active string
			rec, err := t.OpenRecord()
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			if rec != nil {
				active = rec.Project
			}

			str, err := renderDayChart(t, reporter, active, start, blocksPerHour, &[]rune(t.Config.EmptyCell)[0])
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
			fmt.Print(str)
		},
	}

	day.Flags().IntVarP(&blocksPerHour, "width", "w", 3, "Width of the graph, in characters per hour. Auto-scale if not specified")

	return day
}

func weekReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var blocksPerHour int

	day := &cobra.Command{
		Use:     "week [DATE]",
		Short:   "Report of activities over a week in the form of a schedule",
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

			if blocksPerHour <= 0 {
				out.Err("failed to generate report: argument --width must be > 0")
				return
			}
			if !cmd.Flags().Changed("width") {
				if w, _, err := util.TerminalSize(); err == nil && w > 0 {
					blocksPerHour = (w - 14) / 7
				}
			}

			err = schedule(t, start, options, true, blocksPerHour)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
		},
	}

	day.Flags().IntVarP(&blocksPerHour, "width", "w", 12, "Width of the graph, in characters per hour. Auto-scale if not specified")

	return day
}

func dayReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var blocksPerHour int

	day := &cobra.Command{
		Use:     "day [DATE]",
		Short:   "Report of activities over a day in the form of a schedule",
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
			if !cmd.Flags().Changed("width") {
				if w, _, err := util.TerminalSize(); err == nil && w > 0 {
					blocksPerHour = (w - 14)
				}
			}

			err = schedule(t, start, options, false, blocksPerHour)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
		},
	}

	day.Flags().IntVarP(&blocksPerHour, "width", "w", 60, "Width of the graph, in characters per hour. Auto-scale if not specified")

	return day
}

func schedule(t *core.Track, start time.Time, options *filterOptions, week bool, bph int) error {
	var filterStart, filterEnd time.Time

	if week {
		weekDay := (int(start.Weekday()) + 6) % 7
		start = start.Add(time.Duration(-weekDay * 24 * int(time.Hour)))

		filterStart = start.Add(-time.Hour * 24)
		filterEnd = start.Add(time.Hour * 24 * 7)
	} else {
		filterStart = start.Add(-time.Hour * 24)
		filterEnd = start.Add(time.Hour * 24)
	}

	projects, err := t.LoadAllProjects()
	if err != nil {
		return err
	}

	filters, err := createFilters(options, projects, false)
	if err != nil {
		return err
	}
	filters = append(filters, core.FilterByTime(filterStart, filterEnd))

	reporter, err := core.NewReporter(t, options.projects, filters, options.includeArchived, start, filterEnd)
	if err != nil {
		return err
	}
	var active string
	rec, err := t.OpenRecord()
	if err != nil {
		return err
	}
	if rec != nil {
		active = rec.Project
	}

	str, err := renderWeekSchedule(t, reporter, active, start, week, bph)
	if err != nil {
		return err
	}
	fmt.Print(str)
	return nil
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

func renderDayChart(t *core.Track, reporter *core.Reporter, active string, startDate time.Time, blocksPerHour int, space *rune) (string, error) {
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

	fmt.Printf("                    |%s : %s/cell\n",
		startDate.Format(util.DateFormat),
		time.Duration(1e9*(int(time.Hour)/(bph*1e9))).String(),
	)

	timelineStr := map[string]string{}
	for pr, values := range timelines {
		runes := make([]rune, bph*24, bph*24)
		for i, v := range values {
			runes[i] = util.FloatToBlock(v, space)
		}
		timelineStr[pr] = toDayChart(runes, interval*bph)
	}
	timelineStr[t.WorkspaceLabel()] = toDayChartAxis(bph, interval*bph)

	formatter := util.NewTreeFormatter(
		func(t *core.ProjectNode, indent int) string {
			fillLen := 16 - (indent + utf8.RuneCountInString(t.Value.Name))
			name := t.Value.Name
			if fillLen < 0 {
				nameRunes := []rune(name)
				name = string(nameRunes[:len(nameRunes)+fillLen-1]) + "."
			}
			var str string
			if t.Value.Name == active {
				str = color.BgBlue.Sprintf("%s", name)
			} else {
				str = fmt.Sprintf("%s", name)
			}
			if fillLen > 0 {
				str += strings.Repeat(" ", fillLen)
			}
			str += " "
			str += t.Value.Render.Sprintf(" %s ", t.Value.Symbol)
			text, _ := timelineStr[t.Value.Name]
			return fmt.Sprintf("%s%s", str, text)
		},
		2,
	)
	return formatter.FormatTree(tree), nil
}

func toDayChart(runes []rune, interval int) string {
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

func toDayChartAxis(blocksPerHour int, interval int) string {
	fill := strings.Repeat(" ", interval-5)

	sb := strings.Builder{}

	for i := 0; i < 24*blocksPerHour; i += interval {
		fmt.Fprintf(&sb, "|%02d:00%s", i/blocksPerHour, fill)
	}
	fmt.Fprint(&sb, "|")

	return sb.String()
}

func renderWeekSchedule(t *core.Track, reporter *core.Reporter, active string, startDate time.Time, week bool, blocksPerHour int) (string, error) {
	bph := blocksPerHour

	spaceSym := []rune(t.Config.EmptyCell)[0]
	pauseSym := []rune(t.Config.PauseCell)[0]

	projects := maps.Keys(reporter.Projects)
	sort.Strings(projects)
	indices := make(map[string]int, len(projects))
	symbols := make([]rune, len(projects)+1, len(projects)+1)
	colors := make([]color.Style256, len(projects)+1, len(projects)+1)
	symbols[0] = spaceSym
	colors[0] = *color.S256(15, 0)
	for i, p := range projects {
		indices[p] = i + 1
		symbols[i+1] = []rune(reporter.Projects[p].Symbol)[0]
		colors[i+1] = reporter.Projects[p].Render
	}

	numDays := 1
	if week {
		numDays = 7
	}

	timeline := make([]int, 24*numDays*blocksPerHour, 24*numDays*blocksPerHour)
	paused := make([]bool, 24*numDays*blocksPerHour, 24*numDays*blocksPerHour)

	now := time.Now()

	for _, rec := range reporter.Records {
		startIdx, endIdx, ok := toIndexRange(rec.Start, rec.End, startDate, bph, numDays)
		if !ok {
			continue
		}
		index := indices[rec.Project]
		for i := startIdx; i <= endIdx; i++ {
			timeline[i] = index
		}
		for _, p := range rec.Pause {
			startIdx, endIdx, ok := toIndexRange(p.Start, p.End, startDate, bph, numDays)
			if !ok {
				continue
			}
			for i := startIdx; i <= endIdx; i++ {
				paused[i] = true
			}
		}
	}

	nowIdx := int(now.Sub(startDate).Hours() * float64(bph))

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "      |Day %s : %s/cell\n",
		startDate.Format(util.DateFormat),
		time.Duration(1e9*(int(time.Hour)/(bph*1e9))).String(),
	)

	fmt.Fprint(&sb, "      ")
	for weekday := 0; weekday < numDays; weekday++ {
		date := startDate.Add(time.Duration(weekday * 24 * int(time.Hour)))
		str := fmt.Sprintf(
			"%s %02d %s",
			date.Weekday().String()[:2],
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
		for weekday := 0; weekday < numDays; weekday++ {
			s := (weekday*24 + hour) * bph
			fmt.Fprint(&sb, "|")
			for i := s; i < s+bph; i++ {
				pr := timeline[i]
				pause := paused[i]
				sym := symbols[pr]
				col := colors[pr]
				if pause {
					sym = pauseSym
				}
				if i == nowIdx {
					sym = '@'
				}
				fmt.Fprint(&sb, col.Sprintf("%c", sym))
			}
		}
		fmt.Fprintln(&sb, "|")
	}

	totalWidth := 7 + numDays*(bph+1)
	lineWidth := 0

	line1 := ""
	line2 := ""
	for i, p := range projects {
		col := colors[i+1]
		width := utf8.RuneCountInString(p)
		if width < 3 {
			width = 3
		}
		if lineWidth > 0 && lineWidth+width+4 > totalWidth {
			lineWidth = 0
			fmt.Fprintln(&sb, line1)
			fmt.Fprintln(&sb, line2)
			line1 = ""
			line2 = ""
		}

		line1 += col.Sprintf(" %c:%3s ", symbols[indices[p]], p)
		line2 += col.Sprintf(" %*s ", width+2, util.FormatDuration(reporter.TotalTime[p]))
		lineWidth += width + 4
	}
	if len(line1) > 0 {
		fmt.Fprintln(&sb, line1)
		fmt.Fprintln(&sb, line2)
	}

	return sb.String(), nil
}

func toIndexRange(start, end, startDate time.Time, bph int, days int) (int, int, bool) {
	if end.IsZero() {
		end = time.Now()
	}
	if end.Before(startDate) {
		return -1, -1, false
	}

	startIdx := int(start.Sub(startDate).Hours() * float64(bph))
	endIdx := int(end.Sub(startDate).Hours() * float64(bph))
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx >= bph*24*days {
		endIdx = bph*24*days - 1
	}
	return startIdx, endIdx, true
}
