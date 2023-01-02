package cli

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gookit/color"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

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
