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
		Use:   "export",
		Short: "Export resources",
		Long: `Export resources

Currently, only export of (potentially filtered) records to CSV is supported.`,
		Aliases: []string{"ex"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	export.PersistentFlags().StringSliceVarP(&options.projects, "projects", "p", []string{}, "Projects to include (comma-separated). All projects if not specified")
	export.PersistentFlags().StringSliceVarP(&options.tags, "tags", "t", []string{}, "Tags to include (comma-separated). Includes records with any of the given tags")
	export.PersistentFlags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	export.PersistentFlags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")
	export.PersistentFlags().BoolVarP(&options.includeArchived, "archived", "a", false, "Include records from archived projects")

	export.AddCommand(exportRecordsCommand(t, &options))

	export.Long += "\n\n" + formatCmdTree(export)
	return export
}

func exportRecordsCommand(t *core.Track, options *filterOptions) *cobra.Command {
	records := &cobra.Command{
		Use:   "records",
		Short: "Export records",
		Long: `Export records

Currently, only export to CSV is supported.`,
		Aliases: []string{"r"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to export records: %s", err)
				return
			}

			filters, err := createFilters(options, projects, true)
			if err != nil {
				out.Err("failed to export records: %s", err)
				return
			}

			io := os.Stdout
			writer := csvWriter{
				Separtor: ",",
			}
			writer.WriteHeader(io)

			fn, results, _ := t.AllRecordsFiltered(filters, false)

			go fn()

			for res := range results {
				if res.Err != nil {
					out.Err("failed to export records: %s", res.Err)
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
