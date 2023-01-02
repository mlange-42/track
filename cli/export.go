package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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

	export.AddCommand(exportRecordsCommand(t))

	export.Long += "\n\n" + formatCmdTree(export)
	return export
}

func exportRecordsCommand(t *core.Track) *cobra.Command {
	options := filterOptions{}

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

			filters, err := createFilters(&options, projects, true)
			if err != nil {
				out.Err("failed to export records: %s", err)
				return
			}

			io := os.Stdout
			writer := csvWriter{
				Separator: ",",
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

	records.Flags().StringSliceVarP(&options.projects, "projects", "p", []string{}, "Projects to include (comma-separated). All projects if not specified")
	records.Flags().StringSliceVarP(&options.tags, "tags", "t", []string{}, "Tags to include (comma-separated). Includes records with any of the given tags")
	records.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	records.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")
	records.Flags().BoolVarP(&options.includeArchived, "archived", "a", false, "Include records from archived projects")

	return records
}

type csvWriter struct {
	Separator string
}

func (wr csvWriter) WriteHeader(w io.Writer) error {
	_, err := fmt.Fprintf(
		w, "%s\n",
		strings.Join([]string{"start", "end", "project", "total", "work", "pause", "note", "tags"}, wr.Separator),
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

	noTime := time.Time{}

	_, err := fmt.Fprintf(
		w, "%s\n",
		strings.Join([]string{
			r.Start.Format(util.DateTimeFormat),
			endTime,
			r.Project,
			util.FormatDuration(r.TotalDuration(noTime, noTime)),
			util.FormatDuration(r.Duration(noTime, noTime)),
			util.FormatDuration(r.PauseDuration(noTime, noTime)),
			fmt.Sprintf("\"%s\"", strings.ReplaceAll(r.Note, "\n", "\\n")),
			strings.Join(r.Tags, " "),
		}, wr.Separator),
	)
	return err
}
