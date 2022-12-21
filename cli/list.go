package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func listCommand(t *core.Track) *cobra.Command {
	create := &cobra.Command{
		Use:   "list",
		Short: "List resources",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	create.AddCommand(listProjectsCommand(t))
	create.AddCommand(listRecordsCommand(t))

	return create
}

func listProjectsCommand(t *core.Track) *cobra.Command {
	listProjects := &cobra.Command{
		Use:   "projects",
		Short: "List all projects",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to load projects: %s", err)
				return
			}

			for _, project := range projects {
				out.Success("%s", project.Name)
			}
		},
	}

	return listProjects
}

func listRecordsCommand(t *core.Track) *cobra.Command {
	listProjects := &cobra.Command{
		Use:   "records",
		Short: "List all records",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			records, err := t.LoadAllRecords()
			if err != nil {
				out.Err("failed to load records: %s", err)
				return
			}

			for _, record := range records {
				printRecord(record)
			}
		},
	}

	return listProjects
}

func printRecord(r core.Record) {
	start := r.Start.Format(util.DateTimeFormat)

	var end string
	if r.HasEnded() {
		end = r.End.Format(util.DateTimeFormat)
	} else {
		end = util.NoDateTime
	}
	dur := r.Duration()

	out.Success(
		"%-15s %s - %s (%.1fhr)  %s\n", r.Project,
		start, end, dur.Hours(),
		r.Note,
	)
}
