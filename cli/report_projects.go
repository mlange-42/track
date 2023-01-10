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
		Short:   "Shows the project tree with time statistics",
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}

			filters, err := createFilters(options, projects, false)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}

			startTime, endTime, err := parseStartEnd(options)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}
			reporter, err := core.NewReporter(
				t, options.projects, filters,
				options.includeArchived, startTime, endTime,
			)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}

			tree, err := t.ToProjectTree(reporter.Projects)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}
			var active string
			rec, err := t.OpenRecord()
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
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
						str = name
					}
					if fillLen > 0 {
						str += strings.Repeat(" ", fillLen)
					}
					str += " "
					str += t.Value.Render.Sprintf(" %s ", t.Value.Symbol)

					return fmt.Sprintf(
						"%s %6s (%6s)", str,
						util.FormatDuration(reporter.TotalTime[t.Value.Name], false),
						util.FormatDuration(reporter.ProjectTime[t.Value.Name], false),
					)
				},
				2,
			)
			out.Print(formatter.FormatTree(tree))
			return nil
		},
	}
	projects.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	projects.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	return projects
}
