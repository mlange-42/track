package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/spf13/cobra"
)

func reportCommand(t *core.Track) *cobra.Command {
	options := filterOptions{}

	report := &cobra.Command{
		Use:     "report",
		Short:   "Generate reports of time tracking",
		Long:    "Generate reports of time tracking",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	report.PersistentFlags().StringSliceVarP(&options.projects, "projects", "p", []string{}, "Projects to include (comma-separated). All projects if not specified")
	report.PersistentFlags().StringSliceVarP(&options.tags, "tags", "t", []string{}, "Tags to include (comma-separated). Includes records with any of the given tags")
	report.PersistentFlags().BoolVarP(&options.includeArchived, "archived", "a", false, "Include records from archived projects")

	report.AddCommand(timelineReportCommand(t, &options))
	report.AddCommand(projectsReportCommand(t, &options))
	report.AddCommand(chartReportCommand(t, &options))
	report.AddCommand(weekReportCommand(t, &options))
	report.AddCommand(dayReportCommand(t, &options))
	report.AddCommand(treemapReportCommand(t, &options))

	report.Long += "\n\n" + formatCmdTree(report)
	return report
}
