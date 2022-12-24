package core

import (
	"time"
)

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Duration calculates the duration of a time range
func (r TimeRange) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

// Reporter for generating reports
type Reporter struct {
	Track       *Track
	Records     []Record
	Projects    map[string]Project
	ProjectTime map[string]time.Duration
	TimeRange   TimeRange
}

// NewReporter creates a new Reporter from filters
func NewReporter(t *Track, proj []string, filters FilterFunctions) (*Reporter, error) {
	var err error
	projects := make(map[string]Project)
	if len(proj) == 0 {
		projects, err = t.LoadAllProjects()
		if err != nil {
			return nil, err
		}
	} else {
		for _, p := range proj {
			project, err := t.LoadProjectByName(p)
			if err != nil {
				return nil, err
			}
			projects[project.Name] = project
		}
	}

	records, err := t.LoadAllRecordsFiltered(filters)
	if err != nil {
		return nil, err
	}

	totals := make(map[string]time.Duration, len(projects))
	for _, p := range projects {
		totals[p.Name] = time.Second * 0.0
	}

	tRange := TimeRange{}
	for _, rec := range records {
		totals[rec.Project] = totals[rec.Project] + rec.Duration()
		if tRange.Start.IsZero() || rec.Start.Before(tRange.Start) {
			tRange.Start = rec.Start
		}
		if rec.End.IsZero() {
			if tRange.End.IsZero() || rec.Start.After(tRange.End) {
				tRange.End = rec.Start
			}
		} else {
			if tRange.End.IsZero() || rec.End.After(tRange.End) {
				tRange.End = rec.End
			}
		}
	}

	report := Reporter{
		Track:       t,
		Records:     records,
		Projects:    projects,
		ProjectTime: totals,
		TimeRange:   tRange,
	}
	return &report, nil
}
