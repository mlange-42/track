package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type tagStats struct {
	Count int
	Work  time.Duration
	Pause time.Duration
}

func tagsReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	tagsReport := &cobra.Command{
		Use:     "tags",
		Short:   "Shows tags with time statistics",
		Aliases: []string{"t"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to generate report: %s", err.Error())
				return
			}

			filters, err := createFilters(options, projects, false)
			if err != nil {
				out.Err("failed to generate report: %s", err.Error())
				return
			}

			startTime, endTime, err := parseStartEnd(options)
			if err != nil {
				out.Err("failed to generate report: %s", err.Error())
				return
			}
			reporter, err := core.NewReporter(
				t, options.projects, filters,
				options.includeArchived, startTime, endTime,
			)
			if err != nil {
				out.Err("failed to generate report: %s", err.Error())
				return
			}

			tags := map[string]bool{}
			for _, tag := range options.tags {
				tags[tag] = true
			}

			allTags := map[string]*tagStats{}
			for _, rec := range reporter.Records {
				dur := rec.Duration(util.NoTime, util.NoTime)
				pause := rec.PauseDuration(util.NoTime, util.NoTime)
				for _, tag := range rec.Tags {
					if _, ok := tags[tag]; ok || len(tags) == 0 {
						if entry, ok := allTags[tag]; ok {
							entry.Work += dur
							entry.Pause += pause
							entry.Count++
						} else {
							allTags[tag] = &tagStats{
								Count: 1,
								Work:  dur,
								Pause: pause,
							}
						}
					}
				}
			}

			keys := maps.Keys(allTags)
			sort.Strings(keys)

			for _, tag := range keys {
				stats := allTags[tag]
				fillLen := 15 - utf8.RuneCountInString(tag)
				name := tag
				if fillLen < 0 {
					nameRunes := []rune(tag)
					tag = string(nameRunes[:len(nameRunes)+fillLen-1]) + "."
				}

				str := name
				if fillLen > 0 {
					str += strings.Repeat(" ", fillLen)
				}
				fmt.Printf(
					"%s %3d  %s (%s)\n", str,
					stats.Count,
					util.FormatDuration(stats.Work),
					util.FormatDuration(stats.Pause),
				)
			}
		},
	}
	tagsReport.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	tagsReport.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	return tagsReport
}
