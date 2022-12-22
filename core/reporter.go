package core

import (
	"time"
)

// Reporter for generating reports
type Reporter struct {
	Track       *Track
	Records     []Record
	Projects    []Project
	ProjectTime map[string]time.Duration
}

// NewReporter creates a new Reporter from filters
func NewReporter(t *Track, proj []string, filters FilterFunctions) (*Reporter, error) {
	var err error
	var projects []Project
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
			projects = append(projects, project)
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

	for _, rec := range records {
		totals[rec.Project] = totals[rec.Project] + rec.Duration()
	}

	report := Reporter{
		Track:       t,
		Records:     records,
		Projects:    projects,
		ProjectTime: totals,
	}
	return &report, nil
}
