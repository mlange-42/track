package records

import (
	"io"

	"github.com/mlange-42/track/core"
	"gopkg.in/yaml.v3"
)

// YAMLRenderer renders records for YAML export
type YAMLRenderer struct {
	Results chan core.FilterResult
}

// Render renders a stream of records
func (wr YAMLRenderer) Render(w io.Writer) error {
	records := []core.Record{}

	for res := range wr.Results {
		if res.Err != nil {
			return res.Err
		}
		records = append(records, res.Record)
	}

	bytes, err := yaml.Marshal(records)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return err
}
