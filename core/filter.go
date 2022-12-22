package core

import "time"

// FilterFunction is an alias for func(r *Record) bool
type FilterFunction = func(r *Record) bool

// FilterFunctions is an alias for []func(r *Record) bool
type FilterFunctions = []FilterFunction

// Filter checks a record using multiple filters
func Filter(record *Record, filters FilterFunctions) bool {
	for _, f := range filters {
		if !f(record) {
			return false
		}
	}
	return true
}

// FilterByProjects returns a function for filtering by project names
func FilterByProjects(projects []string) FilterFunction {
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
func FilterByTime(start, end time.Time) FilterFunction {
	return func(r *Record) bool {
		return (start.IsZero() || r.End.After(start)) && (end.IsZero() || r.Start.Before(end))
	}
}

// FilterByTagsAny returns a function for filtering by tags
func FilterByTagsAny(tags []string) FilterFunction {
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
func FilterByTagsAll(tags []string) FilterFunction {
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
