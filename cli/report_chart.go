package cli

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gookit/color"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

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
