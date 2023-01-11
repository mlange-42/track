package records

import (
	"encoding/json"
	"io"

	"github.com/mlange-42/track/core"
)

// JSONRenderer renders records for JSON export
type JSONRenderer struct {
	Results chan core.FilterResult
}

// Render renders a stream of records
func (wr JSONRenderer) Render(w io.Writer) error {
	records := []core.Record{}

	for res := range wr.Results {
		if res.Err != nil {
			return res.Err
		}
		records = append(records, res.Record)
	}

	bytes, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return err
}
