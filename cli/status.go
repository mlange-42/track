package cli

import (
	"fmt"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

type statusInfo struct {
	Project   string
	IsActive  bool
	Stopped   time.Duration
	CurrTime  time.Duration
	CumTime   time.Duration
	BreakTime time.Duration
	TotalTime time.Duration
}

func statusCommand(t *core.Track) *cobra.Command {
	var maxBreakStr string

	status := &cobra.Command{
		Use:   "status [PROJECT]",
		Short: "Reports the status of the running or given project",
		Long: `Reports the status of the running or given project

Columns of the status are:

* curr  - Time of the current record
* total - Recorded time today since the last break longer than --max-break
* break - Break time today since the last break longer than --max-break
* today - Total recorded time since midnight
`,
		Aliases: []string{"s"},
		Args:    util.WrappedArgs(cobra.MaximumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			maxBreak, err := time.ParseDuration(maxBreakStr)
			if err != nil {
				out.Err("failed to show status: %s", err)
				return
			}

			project := ""
			if len(args) > 0 {
				project = args[0]
			}
			info, err := getStatus(t, project, maxBreak)
			if err != nil {
				out.Err("failed to show status: %s", err)
				return
			}

			if project == "" && !info.IsActive {
				out.Warn("No running record. Start tracking or specify a project.\n")
				out.Warn(
					"Stopped project '%s' %s ago\n",
					info.Project,
					util.FormatDuration(info.Stopped),
				)
			}

			out.Print("+------------------+-------+-------+-------+-------+\n")
			out.Print("|          project |  curr | total | break | today |\n")
			out.Print(
				"| %16s | %s | %s | %s | %s |\n",
				info.Project,
				util.FormatDuration(info.CurrTime),
				util.FormatDuration(info.CumTime),
				util.FormatDuration(info.BreakTime),
				util.FormatDuration(info.TotalTime),
			)
			out.Print("+------------------+-------+-------+-------+-------+")
		},
	}
	status.Flags().StringVar(
		&maxBreakStr,
		"max-break",
		t.Config.MaxBreakDuration.String(),
		"Maximum length of breaks to consider them in daily break time.\nThe default can be set in the config file",
	)

	return status
}

func getStatus(t *core.Track, proj string, maxBreak time.Duration) (statusInfo, error) {
	var project string
	open, hasOpenRecord := t.OpenRecord()
	stopped := 0 * time.Second
	if proj != "" {
		project = proj
		if hasOpenRecord && open.Project != project {
			hasOpenRecord = false
		}
	} else {
		if !hasOpenRecord {
			last, err := t.LatestRecord()
			if err != nil {
				return statusInfo{}, fmt.Errorf(("No running record. Start tracking or specify a project."))
			}
			stopped = time.Now().Sub(last.End)
			project = last.Project
		} else {
			project = open.Project
		}
	}

	now := time.Now()
	start := util.ToDate(now)
	filterStart := start.Add(-time.Hour * 24)

	filters := core.FilterFunctions{
		core.FilterByTime(filterStart, time.Time{}),
	}

	reporter, err := core.NewReporter(t, []string{project}, filters, false)
	if err != nil {
		return statusInfo{}, err
	}

	prevEnd := time.Time{}
	currTime := time.Second * 0
	cumTime := time.Second * 0
	breakTime := time.Second * 0
	totalTime := time.Second * 0

	for _, rec := range reporter.Records {
		endTime := rec.End
		if endTime.IsZero() {
			endTime = now
		} else {
			if endTime.Before(start) {
				continue
			}
		}
		startTime := rec.Start
		if startTime.Before(start) {
			startTime = start
		}

		worked := endTime.Sub(startTime)

		if !prevEnd.IsZero() {
			bt := startTime.Sub(prevEnd)
			if bt < maxBreak {
				breakTime += bt
			} else {
				cumTime = time.Second * 0
				breakTime = time.Second * 0
			}
		}

		totalTime += worked
		cumTime += worked
		if rec.End.IsZero() {
			currTime += endTime.Sub(rec.Start)
		}

		prevEnd = endTime
	}

	return statusInfo{
		Project:   project,
		IsActive:  hasOpenRecord,
		Stopped:   stopped,
		CurrTime:  currTime,
		CumTime:   cumTime,
		BreakTime: breakTime,
		TotalTime: totalTime,
	}, nil
}
