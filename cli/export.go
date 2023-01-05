package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type recordWriter interface {
	Write(io.Writer, chan core.FilterResult) error
}

func exportCommand(t *core.Track) *cobra.Command {
	export := &cobra.Command{
		Use:     "export",
		Short:   "Export resources",
		Long:    `Export resources`,
		Aliases: []string{"ex"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	export.AddCommand(exportRecordsCommand(t))

	export.Long += "\n\n" + formatCmdTree(export)
	return export
}

func exportRecordsCommand(t *core.Track) *cobra.Command {
	options := filterOptions{}
	var json bool
	var yaml bool

	records := &cobra.Command{
		Use:   "records",
		Short: "Export records",
		Long: `Export records

Records can be exported in CSV, JSON and YAML format.
The default export format is CSV.`,
		Aliases: []string{"r"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to export records: %s", err)
				return
			}

			filters, err := createFilters(&options, projects, true)
			if err != nil {
				out.Err("failed to export records: %s", err)
				return
			}

			io := os.Stdout
			var writer recordWriter
			if json {
				writer = jsonWriter{}
			} else if yaml {
				writer = yamlWriter{}
			} else {
				writer = csvWriter{
					Separator: ",",
				}
			}

			fn, results, _ := t.AllRecordsFiltered(filters, false)
			go fn()
			writer.Write(io, results)
		},
	}

	records.Flags().StringSliceVarP(&options.projects, "projects", "p", []string{}, "Projects to include (comma-separated). All projects if not specified")
	records.Flags().StringSliceVarP(&options.tags, "tags", "t", []string{}, "Tags to include (comma-separated). Includes records with any of the given tags")
	records.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	records.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	records.Flags().BoolVar(&json, "json", false, "Export in JSON format")
	records.Flags().BoolVar(&yaml, "yaml", false, "Export in YAML format")

	records.MarkFlagsMutuallyExclusive("json", "yaml")

	return records
}

type csvWriter struct {
	Separator string
}

func (wr csvWriter) writeHeader(w io.Writer) error {
	_, err := fmt.Fprintf(
		w, "%s\n",
		strings.Join([]string{"start", "end", "project", "total", "work", "pause", "note", "tags"}, wr.Separator),
	)
	return err
}

func (wr csvWriter) Write(w io.Writer, results chan core.FilterResult) error {
	err := wr.writeHeader(w)
	if err != nil {
		return err
	}

	for res := range results {
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
				strings.Join(r.Tags, " "),
			}, wr.Separator),
		)
	}

	return err
}

type jsonWriter struct{}

func (wr jsonWriter) Write(w io.Writer, results chan core.FilterResult) error {
	records := []core.Record{}

	for res := range results {
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

type yamlWriter struct{}

func (wr yamlWriter) Write(w io.Writer, results chan core.FilterResult) error {
	records := []core.Record{}

	for res := range results {
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
