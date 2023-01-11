package cli

import (
	"fmt"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/render/schedule"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			start := util.ToDate(time.Now())

			var err error
			if len(args) > 0 {
				start, err = util.ParseDate(args[0])
				if err != nil {
					return fmt.Errorf("failed to generate report: %s", err)
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
				return fmt.Errorf("failed to generate report: argument --width must be > 0")
			}
			if !cmd.Flags().Changed("width") {
				if w, _, err := util.TerminalSize(); err == nil && w > 0 {
					blocksPerHour = (w - 14) / 7
				}
			}

			err = renderSchedule(t, start, options, true, blocksPerHour)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err)
			}
			return nil
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

			err = renderSchedule(t, start, options, false, blocksPerHour)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}
		},
	}

	day.Flags().IntVarP(&blocksPerHour, "width", "w", 60, "Width of the graph, in characters per hour. Auto-scale if not specified")

	return day
}

func renderSchedule(t *core.Track, start time.Time, options *filterOptions, week bool, bph int) error {
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

	renderer := schedule.TextRenderer{Weekly: week, BlocksPerHour: bph}
	str, err := renderer.Render(t, reporter, start)
	if err != nil {
		return err
	}
	out.Print(str)
	return nil
}
