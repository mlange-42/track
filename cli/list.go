package cli

import (
	"fmt"
	"time"

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

	list.AddCommand(listWorkspacesCommand(t))
	list.AddCommand(listProjectsCommand(t))
	list.AddCommand(listRecordsCommand(t))
	list.AddCommand(listColorsCommand(t))

	list.Long += "\n\n" + formatCmdTree(list)
	return list
}

func listProjectsCommand(t *core.Track) *cobra.Command {
	var includeArchived bool

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
			if !includeArchived {
				pr := make(map[string]core.Project)
				for n, p := range projects {
					if !p.Archived {
						pr[n] = p
					}
				}
				projects = pr
			}

			var active string
			if rec, ok := t.OpenRecord(); ok {
				active = rec.Project
			}

			tree, err := t.ToProjectTree(projects)
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
	listProjects.Flags().BoolVarP(&includeArchived, "archived", "a", false, "Include records from archived projects")

	return listProjects
}

func listWorkspacesCommand(t *core.Track) *cobra.Command {
	listWorkspaces := &cobra.Command{
		Use:     "workspaces",
		Short:   "List all workspaces",
		Aliases: []string{"w"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			ws, err := t.AllWorkspaces()
			if err != nil {
				out.Err("failed to load workspaces: %s", err)
				return
			}

			for i, w := range ws {
				if w == t.Workspace() {
					color.BgBlue.Printf("%s", w)
				} else {
					fmt.Printf("%s", w)
				}
				if i < len(ws)-1 {
					fmt.Print("\n")
				}
			}
		},
	}

	return listWorkspaces
}

func listRecordsCommand(t *core.Track) *cobra.Command {
	var includeArchived bool

	listProjects := &cobra.Command{
		Use:   "records [DATE]",
		Short: "List all records for a date",
		Long: `List all records for a date

The date can either be a date in default formatting, like "2022-12-31",
or a word like "yesterday" or  "today" (the default).`,
		Aliases:    []string{"r"},
		Args:       util.WrappedArgs(cobra.MaximumNArgs(1)),
		ArgAliases: []string{"date"},
		Run: func(cmd *cobra.Command, args []string) {
			date := util.ToDate(time.Now())
			var err error
			if len(args) > 0 {
				date, err = util.ParseDate(args[0])
				if err != nil {
					out.Err("failed to load records: %s", err)
					return
				}
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

			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to export records: %s", err)
				return
			}
			for _, record := range records {
				if includeArchived || !projects[record.Project].Archived {
					printRecord(record)
				}
			}
		},
	}
	listProjects.Flags().BoolVarP(&includeArchived, "archived", "a", false, "Include records from archived projects")

	return listProjects
}

func listColorsCommand(t *core.Track) *cobra.Command {
	listColors := &cobra.Command{
		Use:     "colors",
		Short:   "Lists the 256 available colors",
		Aliases: []string{"c"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			printColorChart()
		},
	}

	return listColors
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

func printColorChart() {
	var row, block, i uint8

	color.C256(0, true).Printf("%3d", 0)
	fmt.Print(" ")
	for i = 1; i < 16; i++ {
		color.S256(0, i).Printf("%3d", i)
		fmt.Print(" ")
	}
	fmt.Print("\n\n")

	const rowOffset uint8 = 6
	const blockOffset uint8 = 36

	for _, start := range []uint8{16, 124} {
		for row = 0; row < 6; row++ {
			for block = 0; block < 3; block++ {
				idx := start + row*rowOffset + block*blockOffset
				for i = 0; i < 6; i++ {
					if row < 3 {
						color.C256(idx+i, true).Printf("%3d", idx+i)
					} else {
						color.S256(0, idx+i).Printf("%3d", idx+i)
					}
					fmt.Print(" ")
				}
				fmt.Print("  ")
			}
			fmt.Println()
		}
		fmt.Println()
	}

	for i = 232; i < 244; i++ {
		color.C256(i, true).Printf("%3d", i)
		fmt.Print(" ")
	}
	fmt.Println()
	for i := 244; i <= 255; i++ {
		color.S256(0, uint8(i)).Printf("%3d", i)
		fmt.Print(" ")
	}
}
