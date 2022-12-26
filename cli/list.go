package cli

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func listCommand(t *core.Track) *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Short:   "List resources",
		Long:    "List resources",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	list.AddCommand(listProjectsCommand(t))
	list.AddCommand(listRecordsCommand(t))

	list.Long += "\n\n" + formatCmdTree(list)
	return list
}

func listProjectsCommand(t *core.Track) *cobra.Command {
	listProjects := &cobra.Command{
		Use:     "projects",
		Short:   "List all projects",
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.NoArgs),
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

			tree, err := core.ToProjectTree(projects)
			if err != nil {
				out.Err("failed to load projects: %s", err)
				return
			}
			formatter := util.NewTreeFormatter(
				func(t *core.ProjectNode, indent int) string {
					if t.Value.Name == active {
						return color.BgBlue.Sprintf("%s", t.Value.Name)
					}
					return fmt.Sprintf("%s", t.Value.Name)
				},
				2,
			)
			fmt.Print(formatter.FormatTree(tree))
		},
	}

	return listProjects
}

func listRecordsCommand(t *core.Track) *cobra.Command {
	listProjects := &cobra.Command{
		Use:   "records <DATE>",
		Short: "List all records for a date",
		Long: `List all records for a date

The date can either be a date in default formatting, like "2022-12-31",
or a word like "today" or "tomorrow".`,
		Aliases:    []string{"r"},
		Args:       util.WrappedArgs(cobra.ExactArgs(1)),
		ArgAliases: []string{"date"},
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

	out.Print(
		"%-15s %s %s - %s (%s)  %s\n", r.Project,
		date, start, end, util.FormatDuration(dur),
		r.Note,
	)
}
