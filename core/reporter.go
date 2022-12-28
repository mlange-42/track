package core

import (
	"fmt"
	"time"

	"golang.org/x/exp/maps"
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
	Track        *Track
	Records      []Record
	Projects     map[string]Project
	ProjectTime  map[string]time.Duration
	TotalTime    map[string]time.Duration
	AllProjects  map[string]Project
	ProjectsTree *ProjectTree
	TimeRange    TimeRange
}

// NewReporter creates a new Reporter from filters
func NewReporter(
	t *Track, proj []string,
	filters FilterFunctions, includeArchived bool,
	start, end time.Time,
) (*Reporter, error) {

	allProjects, err := t.LoadAllProjects()
	if err != nil {
		return nil, err
	}
	for _, p := range proj {
		if _, ok := allProjects[p]; !ok {
			return nil, fmt.Errorf("no project named '%s'", p)
		}
	}

	projectsTree, err := t.ToProjectTree(allProjects)
	if err != nil {
		return nil, fmt.Errorf("duplicate project name: %s", err)
	}

	projects := make(map[string]Project)
	if len(proj) == 0 {
		if includeArchived {
			projects = allProjects
		} else {
			for _, p := range allProjects {
				if !p.Archived {
					projects[p.Name] = p
				}
			}
		}
	} else {
		for _, p := range proj {
			project := allProjects[p]
			projects[project.Name] = project

			desc, ok := projectsTree.Descendants(project.Name)
			if !ok {
				return nil, fmt.Errorf("BUG! Project '%s' not in project tree", project.Name)
			}
			for _, p2 := range desc {
				if includeArchived || !p2.Value.Archived {
					if _, ok = projects[p2.Value.Name]; !ok {
						projects[p2.Value.Name] = p2.Value
					}
				}
			}
		}
	}

	filters = append(filters, FilterByProjects(maps.Keys(projects)))
	records, err := t.LoadAllRecordsFiltered(filters)
	if err != nil {
		return nil, err
	}

	totals := make(map[string]time.Duration, len(projects)+1)
	totals[projectsTree.Root.Value.Name] = time.Second * 0.0
	for _, p := range projects {
		totals[p.Name] = time.Second * 0.0
	}

	now := time.Now()
	tRange := TimeRange{}
	for _, rec := range records {
		rStart := rec.Start
		rEnd := rec.End
		if rEnd.IsZero() {
			rEnd = now
		}
		if !start.IsZero() {
			if rEnd.Before(start) {
				continue
			}
			if rStart.Before(start) {
				rStart = start
			}
		}
		if !end.IsZero() {
			if rStart.After(end) {
				continue
			}
			if rEnd.After(end) {
				rEnd = end
			}
		}

		dur := rEnd.Sub(rStart)
		if dur > 0 {
			totals[rec.Project] = totals[rec.Project] + dur
		}

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

	projectTotals := make(map[string]time.Duration, len(totals))
	for k, v := range totals {
		projectTotals[k] = v
	}

	aggregate(projectsTree.Root, totals, func(a, b time.Duration) time.Duration { return a + b })

	report := Reporter{
		Track:        t,
		Records:      records,
		Projects:     projects,
		ProjectTime:  projectTotals,
		TotalTime:    totals,
		AllProjects:  allProjects,
		ProjectsTree: projectsTree,
		TimeRange:    tRange,
	}
	return &report, nil
}

func aggregate[T any](t *ProjectNode, values map[string]T, fn func(a, b T) T) (T, bool) {
	agg, ok := values[t.Value.Name]
	if !ok {
		return agg, false
	}
	for _, child := range t.Children {
		if v, ok := aggregate(child, values, fn); ok {
			agg = fn(agg, v)
		}
	}
	values[t.Value.Name] = agg
	return agg, true
}
