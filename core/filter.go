package core

import "time"

// Filter checks a record using multiple filters
func Filter(record *Record, filters []func(r *Record) bool) bool {
	for _, f := range filters {
		if !f(record) {
			return false
		}
	}
	return true
}

// FilterByProjects returns a function for filtering by project names
func FilterByProjects(projects []string) func(r *Record) bool {
	prj := make(map[string]bool)
	for _, p := range projects {
		prj[p] = true
	}
	return func(r *Record) bool {
		_, ok := prj[r.Project]
		return ok
	}
}

// FilterByTime returns a function for filtering by time
func FilterByTime(start, end time.Time) func(r *Record) bool {
	return func(r *Record) bool {
		return (start.IsZero() || r.End.After(start)) && (end.IsZero() || r.Start.Before(end))
	}
}

// FilterByTagsAny returns a function for filtering by tags
func FilterByTagsAny(tags []string) func(r *Record) bool {
	tg := make(map[string]bool)
	for _, t := range tags {
		tg[t] = true
	}
	return func(r *Record) bool {
		for _, t := range r.Tags {
			if _, ok := tg[t]; ok {
				return true
			}
		}
		return false
	}
}

// FilterByTagsAll returns a function for filtering by tags
func FilterByTagsAll(tags []string) func(r *Record) bool {
	return func(r *Record) bool {
		for _, t := range tags {
			found := false
			for _, t2 := range r.Tags {
				if t == t2 {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}
