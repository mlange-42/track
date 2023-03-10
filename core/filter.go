package core

import (
	"time"

	"github.com/mlange-42/track/util"
)

// FilterFunction is an alias for func(r *Record) bool
type FilterFunction = func(r *Record) bool

// FilterFunctions are used to filter records.
// FilterFunctions contains []func(r *Record) bool and a time range.
type FilterFunctions struct {
	Functions []FilterFunction
	Start     time.Time
	End       time.Time
}

// NewFilter creates a FilterFunctions struct
func NewFilter(fn []FilterFunction, start, end time.Time) FilterFunctions {
	if !(start.IsZero() && end.IsZero()) {
		fn = append(fn, FilterByTime(start, end))
	}
	return FilterFunctions{
		Functions: fn,
		Start:     start,
		End:       end,
	}
}

// Filter checks a record using multiple filters from FilterFunctions
func Filter(record *Record, filters FilterFunctions) bool {
	for _, f := range filters.Functions {
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
//
// Keeps all records that are partially included in the given time span.
// Zero times in the given time span are ignored, resulting in an open time span.
//
// For records with a zero end, only the start time is compared
func FilterByTime(start, end time.Time) FilterFunction {
	now := time.Now()
	return func(r *Record) bool {
		endTime := r.End
		if endTime.IsZero() {
			endTime = now
		}
		return (start.IsZero() || endTime.After(start)) && (end.IsZero() || r.Start.Before(end))
	}
}

// FilterByArchived returns a function for filtering by archived/not archived
func FilterByArchived(archived bool, projects map[string]Project) FilterFunction {
	return func(r *Record) bool {
		return projects[r.Project].Archived == archived
	}
}

// FilterByTagsAny returns a function for filtering by tags
func FilterByTagsAny(tags []util.Pair[string, string]) FilterFunction {
	tg := map[string]map[string]bool{}
	for _, kv := range tags {
		if m, ok := tg[kv.Key]; ok {
			m[kv.Value] = true
		} else {
			tg[kv.Key] = map[string]bool{kv.Value: true}
		}
	}

	return func(r *Record) bool {
		for t, v := range r.Tags {
			if values, ok := tg[t]; ok {
				if _, ok := values[""]; ok {
					return true
				}
				if _, ok := values[v]; ok {
					return true
				}
			}
		}
		return false
	}
}

// FilterByTagsAll returns a function for filtering by tags
func FilterByTagsAll(tags []util.Pair[string, string]) FilterFunction {
	return func(r *Record) bool {
		for _, kv := range tags {
			found := false
			for t2, v2 := range r.Tags {
				if kv.Key == t2 && (kv.Value == "" || kv.Value == v2) {
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
