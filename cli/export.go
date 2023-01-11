package cli

import (
	"fmt"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/render"
	"github.com/mlange-42/track/render/records"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

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
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to export records: %s", err)
			}

			filters, err := createFilters(&options, projects, true)
			if err != nil {
				return fmt.Errorf("failed to export records: %s", err)
			}

			fn, results, _ := t.AllRecordsFiltered(filters, false)
			go fn()

			io := out.StdOut
			var writer render.Renderer
			if json {
				writer = records.JSONRenderer{Results: results}
			} else if yaml {
				writer = records.YAMLRenderer{Results: results}
			} else {
				writer = records.CsvRenderer{Separator: ",", Results: results}
			}

			writer.Render(io)

			return nil
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
