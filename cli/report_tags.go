package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type tagStats struct {
	Count  int
	Work   time.Duration
	Pause  time.Duration
	Values map[string]*tagValueStats
}

type tagValueStats struct {
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

			tags := map[string]bool{}
			for _, tag := range options.tags {
				k, _ := core.ParseTag(tag)
				tags[k] = true
			}

			valueStats := false
			if len(tags) == 1 {
				valueStats = true
			}

			allTags := map[string]*tagStats{}
			for _, rec := range reporter.Records {
				dur := rec.Duration(util.NoTime, util.NoTime)
				pause := rec.PauseDuration(util.NoTime, util.NoTime)
				for tag, value := range rec.Tags {
					if _, ok := tags[tag]; ok || len(tags) == 0 {
						if _, ok := allTags[tag]; !ok {
							allTags[tag] = &tagStats{
								Count: 0,
								Work:  0,
								Pause: 0,
							}
						}
						entry, _ := allTags[tag]
						entry.Work += dur
						entry.Pause += pause
						entry.Count++

						if entry.Values == nil {
							entry.Values = make(map[string]*tagValueStats)
						}
						if _, ok := entry.Values[value]; !ok {
							entry.Values[value] = &tagValueStats{
								Count: 0,
								Work:  0,
								Pause: 0,
							}
						}
						values, _ := entry.Values[value]
						values.Work += dur
						values.Pause += pause
						values.Count++
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
					nameRunes := []rune(name)
					name = string(nameRunes[:len(nameRunes)+fillLen-1]) + "."
				}

				str := name
				if fillLen > 0 {
					str += strings.Repeat(" ", fillLen)
				}
				fmt.Printf(
					"%s %3d  %6s (%5s)", str,
					stats.Count,
					util.FormatDuration(stats.Work, false),
					util.FormatDuration(stats.Pause, false),
				)
				if !valueStats {
					values := stats.Values
					if values != nil && len(values) > 0 {
						vKeys := maps.Keys(values)
						sort.Strings(vKeys)
						if _, ok := values[""]; !ok || len(vKeys) > 1 {
							fmt.Printf(" [%s]", strings.Join(vKeys, " "))
						}
					}
					fmt.Printf("\n")
				} else {
					fmt.Printf("\n")
					values := stats.Values
					if values == nil {
						continue
					}
					vKeys := maps.Keys(values)
					sort.Strings(vKeys)
					for _, v := range vKeys {
						vStats := values[v]
						fillLen := 13 - utf8.RuneCountInString(v)
						name := v
						if fillLen < 0 {
							nameRunes := []rune(name)
							name = string(nameRunes[:len(nameRunes)+fillLen-1]) + "."
						}

						str := name
						if fillLen > 0 {
							str += strings.Repeat(" ", fillLen)
						}
						fmt.Printf(
							"  %s %3d  %6s (%5s)\n", str,
							vStats.Count,
							util.FormatDuration(vStats.Work, false),
							util.FormatDuration(vStats.Pause, false),
						)
					}
				}
			}
			return nil
		},
	}
	tagsReport.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	tagsReport.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	return tagsReport
}
