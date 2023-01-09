package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/mlange-42/track/util"
)

var builder = strings.Builder{}

// SerializeRecord converts a record to a serialization string
func SerializeRecord(r *Record, date time.Time) string {
	builder.Reset()

	reference := date
	if reference.IsZero() {
		reference = r.Start
	}
	startDate := util.FormatTimeWithOffset(r.Start, reference)
	endTime := util.FormatTimeWithOffset(r.End, reference)

	fmt.Fprintf(&builder, "%s - %s", startDate, endTime)
	for _, p := range r.Pause {
		startTime := util.FormatTimeWithOffset(p.Start, reference)
		if p.End.IsZero() {
			fmt.Fprintf(&builder, "\n    - %s - ?", startTime)
		} else {
			fmt.Fprintf(&builder, "\n    - %s - %s", startTime, p.End.Sub(p.Start).Round(time.Second))
		}
		if p.Note != "" {
			fmt.Fprintf(&builder, " / %s", p.Note)
		}
	}
	fmt.Fprintf(&builder, "\n    %s", r.Project)

	if len(r.Note) > 0 {
		fmt.Fprintf(&builder, "\n\n%s", r.Note)
	}
	fmt.Fprint(&builder, "\n")
	return builder.String()
}

// DeserializeRecord converts a serialization string to a record
func DeserializeRecord(str string, date time.Time) (Record, error) {
	str = strings.TrimSpace(str)
	lines := strings.Split(strings.ReplaceAll(str, "\r\n", "\n"), "\n")
	index, ok := skipLines(lines, 0, true)
	if !ok {
		return Record{}, fmt.Errorf("invalid record: missing time range (1st line)")
	}
	start, end, err := util.ParseTimeRange(lines[index], date)
	index++
	if err != nil {
		return Record{}, err
	}

	pause := []Pause{}
	for {
		ln := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(ln, "- ") {
			break
		}
		ln = strings.TrimPrefix(ln, "- ")
		lnParts := strings.SplitN(ln, "/", 2)
		pStart, pEnd, err := util.ParseTimeRange(lnParts[0], date)
		index++
		if err != nil {
			return Record{}, err
		}
		note := ""
		if len(lnParts) > 1 {
			note = strings.TrimSpace(lnParts[1])
		}
		pause = append(pause,
			Pause{
				Start: pStart,
				End:   pEnd,
				Note:  note,
			},
		)
	}

	index, ok = skipLines(lines, index, true)
	if !ok {
		return Record{}, fmt.Errorf("invalid record: missing project (2nd line)")
	}
	projectName := strings.TrimSpace(lines[index])
	index++

	notes := []string{}
	index, ok = skipLines(lines, index, true)
	if ok {
		for ok {
			notes = append(notes, lines[index])
			index++
			index, ok = skipLines(lines, index, false)
		}
	}
	tags, err := ExtractTagsSlice(notes)
	if err != nil {
		return Record{}, err
	}

	return Record{
		Project: projectName,
		Start:   start,
		End:     end,
		Note:    strings.TrimSpace(strings.Join(notes, "\n")),
		Tags:    tags,
		Pause:   pause,
	}, nil
}

func skipLines(lines []string, index int, skipEmpty bool) (int, bool) {
	if index >= len(lines) {
		return index, false
	}
	for (skipEmpty && strings.TrimSpace(lines[index]) == "") || strings.HasPrefix(lines[index], CommentPrefix) {
		index++
		if index >= len(lines) {
			return index, false
		}
	}
	return index, true
}
