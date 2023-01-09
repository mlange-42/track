package core

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mlange-42/track/util"
)

// TagPrefix denotes tags in record notes
const TagPrefix = "+"

// CommentPrefix denotes comments in record files
const CommentPrefix = "#"

// YamlCommentPrefix denotes comments in YAML files
const YamlCommentPrefix = "#"

var (
	// ErrNoRecords is an error for no records found for a date
	ErrNoRecords = errors.New("no records for date")
	// ErrRecordNotFound is an error for a particular record not found
	ErrRecordNotFound = errors.New("record not found")
)

// Record holds and manipulates data for a record
type Record struct {
	Project string            `json:"project"`
	Start   time.Time         `json:"start"`
	End     time.Time         `json:"end"`
	Note    string            `json:"note"`
	Tags    map[string]string `json:"tags"`
	Pause   []Pause           `json:"pause"`
}

// Pause holds information about a pause in a record
type Pause struct {
	Start time.Time
	End   time.Time
	Note  string
}

// Duration reports the duration of a pause
func (p *Pause) Duration(min, max time.Time) time.Duration {
	return util.DurationClip(p.Start, p.End, min, max)
}

// HasEnded reports whether the record has an end time
func (r *Record) HasEnded() bool {
	return !r.End.IsZero()
}

// IsPaused reports whether the record is paused
func (r *Record) IsPaused() bool {
	return len(r.Pause) > 0 && r.Pause[len(r.Pause)-1].End.IsZero()
}

// CurrentPause returns the current pause
func (r *Record) CurrentPause() (Pause, bool) {
	if len(r.Pause) > 0 && r.Pause[len(r.Pause)-1].End.IsZero() {
		return r.Pause[len(r.Pause)-1], true
	}
	return Pause{}, false
}

// LastPause returns the last pause. Returns false if there is an open pause
func (r *Record) LastPause() (Pause, bool) {
	if len(r.Pause) > 0 && !r.Pause[len(r.Pause)-1].End.IsZero() {
		return r.Pause[len(r.Pause)-1], true
	}
	return Pause{}, false
}

// Duration reports the duration of a record, excluding pause time
func (r *Record) Duration(min, max time.Time) time.Duration {
	dur := util.DurationClip(r.Start, r.End, min, max)
	dur -= r.PauseDuration(min, max)
	return dur
}

// TotalDuration reports the duration of a record, including pause time
func (r *Record) TotalDuration(min, max time.Time) time.Duration {
	return util.DurationClip(r.Start, r.End, min, max)
}

// PauseDuration reports the duration of all pauses of the record
func (r *Record) PauseDuration(min, max time.Time) time.Duration {
	dur := time.Second * 0
	for _, p := range r.Pause {
		dur += p.Duration(min, max)
	}
	return dur
}

// CurrentPauseDuration reports the duration of an open pause
func (r *Record) CurrentPauseDuration(min, max time.Time) time.Duration {
	dur := time.Second * 0
	if len(r.Pause) == 0 {
		return dur
	}
	if !r.IsPaused() {
		return dur
	}
	last := r.Pause[len(r.Pause)-1]

	return last.Duration(min, max)
}

// Check checks consistency of a record
func (r *Record) Check() error {
	if !r.End.IsZero() && r.End.Before(r.Start) {
		return fmt.Errorf("end time is before start time")
	}
	prevStart := util.NoTime
	prevEnd := util.NoTime
	for _, p := range r.Pause {
		if p.Start.Before(r.Start) {
			return fmt.Errorf("pause starts before record")
		}
		if !p.End.IsZero() && p.End.Before(p.Start) {
			return fmt.Errorf("pause ends before its start")
		}
		if !r.End.IsZero() {
			if p.End.IsZero() {
				return fmt.Errorf("pause is ongoing but record is finished")
			}
			if p.End.After(r.End) {
				return fmt.Errorf("pause ends after record")
			}
		}
		if prevStart.After(p.Start) {
			return fmt.Errorf("pause starts not in chronological order")
		}
		if prevEnd.After(p.Start) {
			return fmt.Errorf("pauses overlap")
		}
		prevStart = p.Start
		prevEnd = p.End
	}
	return nil
}

// InsertPause inserts a pause into a record
func (r *Record) InsertPause(start time.Time, end time.Time, note string) (Pause, error) {
	if len(r.Pause) == 0 {
		if start.Before(r.Start) {
			return Pause{}, fmt.Errorf("start of pause before start of current record")
		}
	} else {
		if start.Before(r.Pause[len(r.Pause)-1].End) {
			return Pause{}, fmt.Errorf("start of pause before end of previous pause")
		}
	}
	r.Pause = append(r.Pause, Pause{Start: start, End: end, Note: note})
	return r.Pause[len(r.Pause)-1], nil
}

// PopPause pops the last pause
func (r *Record) PopPause() (Pause, bool) {
	if len(r.Pause) == 0 {
		return Pause{}, false
	}
	p := r.Pause[len(r.Pause)-1]
	r.Pause = r.Pause[:len(r.Pause)-1]
	return p, true
}

// EndPause closes the last, open pause
func (r *Record) EndPause(t time.Time) (Pause, error) {
	if len(r.Pause) == 0 {
		return Pause{}, fmt.Errorf("no pause to end")
	}
	if !r.Pause[len(r.Pause)-1].End.IsZero() {
		return Pause{}, fmt.Errorf("last pause is already ended")
	}
	r.Pause[len(r.Pause)-1].End = t
	return r.Pause[len(r.Pause)-1], nil
}

// ExtractTagsSlice extracts elements with the tag prefix
func ExtractTagsSlice(tokens []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, token := range tokens {
		subTokens := strings.Split(token, " ")
		for _, subToken := range subTokens {
			if strings.HasPrefix(subToken, TagPrefix) {
				parts := strings.SplitN(strings.TrimPrefix(subToken, TagPrefix), "=", 2)
				value := ""
				if len(parts) > 1 {
					value = parts[1]
				}
				if old, ok := result[parts[0]]; ok && value != old {
					return nil, fmt.Errorf("tag '%s' already has value '%s'", parts[0], value)
				}
				result[parts[0]] = value
			}
		}
	}
	return result, nil
}

// ExtractTags extracts elements with the tag prefix
func ExtractTags(text string) (map[string]string, error) {
	result := make(map[string]string)
	subTokens := strings.Split(text, " ")
	for _, subToken := range subTokens {
		parts := strings.SplitN(strings.TrimPrefix(subToken, TagPrefix), "=", 2)
		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}
		if old, ok := result[parts[0]]; ok && value != old {
			return nil, fmt.Errorf("tag '%s' already has value '%s'", parts[0], value)
		}
		result[parts[0]] = value
	}
	return result, nil
}

func pathToTime(y, m, d, file string) (time.Time, error) {
	return time.ParseInLocation(
		util.FileDateTimeFormat,
		fmt.Sprintf("%s-%s-%s %s", y, m, d, strings.Split(file, ".")[0]),
		time.Local,
	)
}

func fileToTime(date time.Time, file string) (time.Time, error) {
	t, err := time.ParseInLocation(util.FileTimeFormat, strings.Split(file, ".")[0], time.Local)
	if err != nil {
		return util.NoTime, err
	}
	return util.DateAndTime(date, t), nil
}
