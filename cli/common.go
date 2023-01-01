package cli

import (
	"fmt"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func formatCmdTree(command *cobra.Command) string {
	str, err := util.FormatCmdTree(command)
	if err != nil {
		panic(err)
	}
	return str
}

type filterOptions struct {
	projects        []string
	tags            []string
	start           string
	end             string
	includeArchived bool
}

func createFilters(options *filterOptions, projects map[string]core.Project, filterProjects bool) (core.FilterFunctions, error) {
	var filters core.FilterFunctions

	if filterProjects && len(options.projects) > 0 {
		filters = append(filters, core.FilterByProjects(options.projects))
	}

	if !options.includeArchived {
		filters = append(filters, core.FilterByArchived(false, projects))
	}

	if len(options.tags) > 0 {
		filters = append(filters, core.FilterByTagsAny(options.tags))
	}

	startTime, endTime, err := parseStartEnd(options)
	if err != nil {
		return nil, err
	}
	if !(startTime.IsZero() && endTime.IsZero()) {
		filters = append(filters, core.FilterByTime(startTime, endTime))
	}

	return filters, nil
}

func parseStartEnd(options *filterOptions) (time.Time, time.Time, error) {
	var err error
	var startTime time.Time
	var endTime time.Time
	if len(options.start) > 0 {
		startTime, err = util.ParseDate(options.start)
		if err != nil {
			return startTime, endTime, err
		}
	}
	if len(options.end) > 0 {
		endTime, err = util.ParseDate(options.end)
		if err != nil {
			return startTime, endTime, err
		}
		endTime = endTime.Add(time.Hour * 24)
	}
	return startTime, endTime, nil
}

func confirm(question string) bool {
	answer, err := out.Scan(question)
	if err != nil {
		return false
	}
	if answer != "yes!" {
		return false
	}
	return true
}

func getStopTime(open *core.Record, ago time.Duration, at string) (time.Time, error) {
	now := time.Now()
	stopTime := now
	if ago != 0 {
		stopTime = stopTime.Add(-ago)
	}
	if at != "" {
		var err error
		stopTime, err = util.ParseDateTime(fmt.Sprintf("%s %s", stopTime.Format(util.DateFormat), at))
		if err != nil {
			return time.Time{}, err
		}
		if stopTime.After(now) {
			altTime := stopTime.Add(-24 * time.Hour)
			if altTime.Before(now) && altTime.After(open.Start) {
				stopTime = altTime
			}
		}
	}
	if stopTime.After(now) {
		return stopTime, fmt.Errorf("can't stop at a time in the future")
	}
	if stopTime.Before(open.Start) {
		return stopTime, fmt.Errorf("can't stop at a time before the start of the record")
	}
	return stopTime, nil
}

func getStartTime(lastEnd time.Time, ago time.Duration, at string) (time.Time, error) {
	now := time.Now()
	startTime := now
	if ago != 0 {
		startTime = startTime.Add(-ago)
	}
	if at != "" {
		var err error
		startTime, err = util.ParseDateTime(fmt.Sprintf("%s %s", startTime.Format(util.DateFormat), at))
		if err != nil {
			return time.Time{}, err
		}
		if !lastEnd.IsZero() && startTime.After(now) {
			altTime := startTime.Add(-24 * time.Hour)
			if altTime.Before(now) && altTime.After(lastEnd) {
				startTime = altTime
			}
		}
	}
	if startTime.After(now) {
		return startTime, fmt.Errorf("can't start at a time in the future")
	}
	if !lastEnd.IsZero() && startTime.Before(lastEnd) {
		return startTime, fmt.Errorf("can't start at a time before the last stop/pause")
	}
	return startTime, nil
}
