package records

import (
	"fmt"
	"io"
	"strings"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
)

// CsvRenderer renders records for CSV export
type CsvRenderer struct {
	Separator string
	Results   chan core.FilterResult
}

func (wr CsvRenderer) writeHeader(w io.Writer) error {
	_, err := fmt.Fprintf(
		w, "%s\n",
		strings.Join([]string{"start", "end", "project", "total", "work", "pause", "note", "tags"}, wr.Separator),
	)
	return err
}

// Render renders a stream of records
func (wr CsvRenderer) Render(w io.Writer) error {
	err := wr.writeHeader(w)
	if err != nil {
		return err
	}

	for res := range wr.Results {
		if res.Err != nil {
			return res.Err
		}
		r := res.Record

		var endTime string
		if r.End.IsZero() {
			endTime = ""
		} else {
			endTime = r.End.Format(util.DateTimeFormat)
		}

		tags := make([]string, len(r.Tags), len(r.Tags))
		i := 0
		for k, v := range r.Tags {
			tags[i] = fmt.Sprintf("%s=%s", k, v)
			i++
		}

		_, err = fmt.Fprintf(
			w, "%s\n",
			strings.Join([]string{
				r.Start.Format(util.DateTimeFormat),
				endTime,
				r.Project,
				util.FormatDuration(r.TotalDuration(util.NoTime, util.NoTime)),
				util.FormatDuration(r.Duration(util.NoTime, util.NoTime)),
				util.FormatDuration(r.PauseDuration(util.NoTime, util.NoTime)),
				fmt.Sprintf("\"%s\"", strings.ReplaceAll(r.Note, "\n", "\\n")),
				strings.Join(tags, " "),
			}, wr.Separator),
		)
	}

	return err
}
