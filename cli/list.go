package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gookit/color"
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
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
	list.AddCommand(listTagsCommand(t))

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
			rec, err := t.OpenRecord()
			if err != nil {
				out.Err("failed to load projects: %s", err)
				return
			}
			if rec != nil {
				active = rec.Project
			}

			tree, err := t.ToProjectTree(projects)
			if err != nil {
				out.Err("failed to load projects: %s", err)
				return
			}
			formatter := util.NewTreeFormatter(
				func(t *core.ProjectNode, indent int) string {
					fillLen := 16 - (indent + utf8.RuneCountInString(t.Value.Name))
					name := t.Value.Name
					if fillLen < 0 {
						nameRunes := []rune(name)
						name = string(nameRunes[:len(nameRunes)+fillLen-1]) + "."
					}
					var str string
					if t.Value.Name == active {
						str = color.BgBlue.Sprintf("%s", name)
					} else {
						str = fmt.Sprintf("%s", name)
					}
					if fillLen > 0 {
						str += strings.Repeat(" ", fillLen)
					}
					str += " "
					str += t.Value.Render.Sprintf(" %s ", t.Value.Symbol)
					return str
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

			records, err := t.LoadDateRecords(date)
			if err != nil {
				if err == core.ErrNoRecords {
					out.Err("no records for date %s", date.Format(util.DateFormat))
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
				project := projects[record.Project]
				if includeArchived || !project.Archived {
					printRecord(record, project)
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

func listTagsCommand(t *core.Track) *cobra.Command {
	var includeArchived bool

	listTags := &cobra.Command{
		Use:     "tags",
		Short:   "Lists all tags",
		Aliases: []string{"t"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			err := printTags(t, includeArchived)
			if err != nil {
				out.Err("failed to list tags: %s", err.Error())
			}
		},
	}
	listTags.Flags().BoolVarP(&includeArchived, "archived", "a", false, "Include records from archived projects")

	return listTags
}

func printRecord(r core.Record, project core.Project) {
	date := r.Start.Format(util.DateFormat)
	start := r.Start.Format(util.TimeFormat)

	var end string
	if r.HasEnded() {
		end = r.End.Format(util.TimeFormat)
	} else {
		end = util.NoTime
	}
	dur := r.Duration(time.Time{}, time.Time{})
	pause := r.PauseDuration(time.Time{}, time.Time{})

	fillLen := 16 - utf8.RuneCountInString(r.Project)
	name := r.Project
	if fillLen < 0 {
		nameRunes := []rune(name)
		name = string(nameRunes[:len(nameRunes)+fillLen-1]) + "."
	}

	fill := ""
	if fillLen > 0 {
		fill = strings.Repeat(" ", fillLen)
	}
	note := ""
	if r.Note != "" {
		parts := strings.SplitN(strings.ReplaceAll(r.Note, "\r\n", "\n"), "\n", 2)
		note = parts[0]
		if len(parts) > 1 {
			note += " ..."
		}
	}
	out.Print(
		"%s%s %s %s %s - %s (%s + %s)  %s\n", name, fill,
		project.Render.Sprintf(" %s ", project.Symbol),
		date, start, end, util.FormatDuration(dur), util.FormatDuration(pause),
		note,
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

func printTags(t *core.Track, includeArchived bool) error {
	projects, err := t.LoadAllProjects()
	if err != nil {
		return err
	}

	tags := map[string]int{}

	filters := core.FilterFunctions{}
	if !includeArchived {
		filters = append(filters, core.FilterByArchived(false, projects))
	}

	fn, results, _ := t.AllRecordsFiltered(filters, false)

	go fn()
	for res := range results {
		if res.Err != nil {
			return res.Err
		}
		for _, tag := range res.Record.Tags {
			if v, ok := tags[tag]; ok {
				tags[tag] = v + 1
			} else {
				tags[tag] = 1
			}
		}
	}

	keys := maps.Keys(tags)
	sort.Strings(keys)

	for _, tag := range keys {
		out.Print("%16s %4d\n", tag, tags[tag])
	}

	return nil
}
