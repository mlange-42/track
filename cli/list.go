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

			var active string
			if rec, ok := t.OpenRecord(); ok {
				active = rec.Project
			}

			for _, project := range projects {
				if project.Name == active {
					out.Success("*%s\n", project.Name)
				} else {
					out.Success(" %s\n", project.Name)
				}
			}
		},
	}

	return listProjects
}

func listRecordsCommand(t *core.Track) *cobra.Command {
	listProjects := &cobra.Command{
		Use:   "records <date>",
		Short: "List all records",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			date, err := util.ParseDate(args[0])
			if err != nil {
				out.Err("failed to load records: %s", err)
				return
			}
			dir := date.Format(util.FileDateFormat)
			records, err := t.LoadDateRecords(dir)
			if err != nil {
				if err == core.ErrNoRecords {
					out.Err("no records for date %s", dir)
					return
				}
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
	date := r.Start.Format(util.DateFormat)
	start := r.Start.Format(util.TimeFormat)

	var end string
	if r.HasEnded() {
		end = r.End.Format(util.TimeFormat)
	} else {
		end = util.NoTime
	}
	dur := r.Duration()

	out.Success(
		"%-15s %s %s - %s (%.1fhr)  %s\n", r.Project,
		date, start, end, dur.Hours(),
		r.Note,
	)
}
