package cli

import (
	"fmt"
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

func weekReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var blocksPerHour int
	var exact bool

	week := &cobra.Command{
		Use:   "week [DATE]",
		Short: "Report of activities over a week in the form of a schedule",
		Long: `Report of activities over a week in the form of a schedule

Reports for the current week if no date is given, or for the past 7 days with flag --7days.

If called with a date, reports for the week containing the date, or for the 7 days starting with the date with flag --7days.`,
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
				if !exact {
					start = util.Monday(start)
				}
			} else {
				if exact {
					start = start.Add(-6 * 24 * time.Hour)
				} else {
					start = util.Monday(start)
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

	week.Flags().IntVarP(&blocksPerHour, "width", "w", 12, "Width of the graph, in characters per hour. Auto-scale if not specified")
	week.Flags().BoolVarP(&exact, "7days", "7", false, "Show the report for 7 days instead of the current/given calendar week")

	return week
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
					blocksPerHour = (w - 8)
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
	filters = core.NewFilter(filters.Functions, filterStart, filterEnd)

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

func renderWeekSchedule(t *core.Track, reporter *core.Reporter, active string, startDate time.Time, week bool, blocksPerHour int) (string, error) {
	bph := blocksPerHour

	spaceSym := []rune(t.Config.EmptyCell)[0]
	pauseSym := []rune(t.Config.PauseCell)[0]
	recordSym := []rune(t.Config.RecordCell)[0]

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
	record := make([]int, 24*numDays*blocksPerHour, 24*numDays*blocksPerHour)

	now := time.Now()

	for recIdx, rec := range reporter.Records {
		startIdx, endIdx, ok := toIndexRange(rec.Start, rec.End, startDate, bph, numDays)
		if !ok {
			continue
		}
		index := indices[rec.Project]
		for i := startIdx; i <= endIdx; i++ {
			timeline[i] = index
			record[i] = recIdx
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

	lastRecord := -1
	idxRecord := 0
	currNote := []rune{}
	currName := []rune{}
	for hour := 0; hour < 24; hour++ {
		fmt.Fprintf(&sb, "%02d:00 ", hour)
		for weekday := 0; weekday < numDays; weekday++ {
			s := (weekday*24 + hour) * bph
			fmt.Fprint(&sb, "|")
			for i := s; i < s+bph; i++ {
				rec := record[i]
				pr := timeline[i]
				pause := paused[i]
				if rec != lastRecord {
					lastRecord = rec
					idxRecord = 0
					if pr == 0 {
						currNote = []rune{}
						currName = []rune{}
					} else {
						currNote = []rune(reporter.Records[rec].Note)
						currName = []rune(reporter.Records[rec].Project)
					}
				} else {
					if !pause {
						idxRecord++
					}
				}

				sym := symbols[pr]
				col := colors[pr]
				if pause {
					sym = pauseSym
				}
				if !week && !pause && pr > 0 {
					nameLen := len(currName)
					noteLen := len(currNote)
					if idxRecord == 0 {
						sym = ' '
					} else if idxRecord-1 < nameLen {
						sym = currName[idxRecord-1]
					} else if idxRecord-1 == nameLen {
						sym = ':'
					} else if idxRecord-1 == nameLen+1 {
						sym = ' '
					} else if idxRecord-3-nameLen < noteLen {
						sym = currNote[idxRecord-3-nameLen]
						if sym == '\n' || sym == '\r' {
							sym = ' '
						}
					} else if idxRecord-3-nameLen == noteLen {
						sym = ' '
					} else {
						sym = recordSym
					}
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
	endIdx := int(end.Sub(startDate).Hours()*float64(bph)) - 1
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx >= bph*24*days {
		endIdx = bph*24*days - 1
	}
	return startIdx, endIdx, true
}
