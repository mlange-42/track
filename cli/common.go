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
	projects []string
	tags     []string
	start    string
	end      string
}

func createFilters(options *filterOptions, projects bool) (core.FilterFunctions, error) {
	var filters core.FilterFunctions

	if projects && len(options.projects) > 0 {
		filters = append(filters, core.FilterByProjects(options.projects))
	}

	if len(options.tags) > 0 {
		filters = append(filters, core.FilterByTagsAny(options.tags))
	}
	var err error
	var startTime time.Time
	var endTime time.Time
	if len(options.start) > 0 {
		startTime, err = util.ParseDate(options.start)
		if err != nil {
			return nil, err
		}
	}
	if len(options.end) > 0 {
		endTime, err = util.ParseDate(options.end)
		if err != nil {
			return nil, err
		}
		endTime = endTime.Add(time.Hour * 24)
	}
	if !(startTime.IsZero() && endTime.IsZero()) {
		filters = append(filters, core.FilterByTime(startTime, endTime))
	}

	return filters, nil
}

func confirmDeleteRecord(rec core.Record) bool {
	question := fmt.Sprintf(
		"Really delete record %s (%s) from project '%s' (y/n): ",
		rec.Start.Format(util.DateTimeFormat),
		util.FormatDuration(rec.Duration()),
		rec.Project,
	)

	answer, err := out.Scan(question)
	if err != nil {
		return false
	}
	if answer != "y" {
		return false
	}
	return true
}
