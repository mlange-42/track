package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

type recordWriter interface {
	WriteHeader(io.Writer) error
	Write(io.Writer, *core.Record) error
}

func exportCommand(t *core.Track) *cobra.Command {
	options := filterOptions{}

	export := &cobra.Command{
		Use:     "export",
		Short:   "Export resources",
		Aliases: []string{"ex"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	export.PersistentFlags().StringSliceVarP(&options.projects, "projects", "p", []string{}, "Projects to include (comma-separated). All projects if not specified")
	export.PersistentFlags().StringSliceVarP(&options.tags, "tags", "t", []string{}, "Tags to include (comma-separated). Includes records with any of the given tags")
	export.PersistentFlags().StringVarP(&options.start, "start", "s", "", "Start date")
	export.PersistentFlags().StringVarP(&options.end, "end", "e", "", "End date")

	export.AddCommand(exportRecordsCommand(t, &options))

	return export
}

func exportRecordsCommand(t *core.Track, options *filterOptions) *cobra.Command {
	records := &cobra.Command{
		Use:     "records",
		Short:   "Export records",
		Aliases: []string{"r"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			filters, err := createFilters(options, true)
			if err != nil {
				out.Err("failed to generate report: %s", err)
				return
			}

			io := os.Stdout
			writer := csvWriter{
				Separtor: ",",
			}
			writer.WriteHeader(io)

			fn, results := t.AllRecordsFiltered(filters)

			go fn()

			for res := range results {
				if res.Err != nil {
					out.Err("failed to generate report: %s", res.Err)
					return
				}
				writer.Write(io, &res.Record)
			}
		},
	}
	return records
}

type csvWriter struct {
	Separtor string
}

func (wr csvWriter) WriteHeader(w io.Writer) error {
	_, err := fmt.Fprintf(
		w, "%s\n",
		strings.Join([]string{"start", "end", "project", "note", "tags"}, wr.Separtor),
	)
	return err
}

func (wr csvWriter) Write(w io.Writer, r *core.Record) error {
	var endTime string
	if r.End.IsZero() {
		endTime = ""
	} else {
		endTime = r.End.Format(util.DateTimeFormat)
	}

	_, err := fmt.Fprintf(
		w, "%s\n",
		strings.Join([]string{
			r.Start.Format(util.DateTimeFormat),
			endTime,
			r.Project,
			r.Note,
			strings.Join(r.Tags, " "),
		}, wr.Separtor),
	)
	return err
}
