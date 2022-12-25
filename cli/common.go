package cli

import (
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
)

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
	}
	if !(startTime.IsZero() && endTime.IsZero()) {
		filters = append(filters, core.FilterByTime(startTime, endTime))
	}

	return filters, nil
}
