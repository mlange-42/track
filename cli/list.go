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
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to list projects: %s", err)
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
				return fmt.Errorf("failed to list projects: %s", err)
			}
			if rec != nil {
				active = rec.Project
			}

			tree, err := t.ToProjectTree(projects)
			if err != nil {
				return fmt.Errorf("failed to list projects: %s", err)
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
			out.Print(formatter.FormatTree(tree))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, err := t.AllWorkspaces()
			if err != nil {
				return fmt.Errorf("failed to load workspaces: %s", err)
			}

			for i, w := range ws {
				if w == t.Workspace() {
					color.BgBlue.Printf("%s", w)
				} else {
					out.Print("%s", w)
				}
				if i < len(ws)-1 {
					out.Print("\n")
				}
			}
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			date := util.ToDate(time.Now())
			var err error
			if len(args) > 0 {
				date, err = util.ParseDate(args[0])
				if err != nil {
					return fmt.Errorf("failed to load records: %s", err)
				}
			}

			records, err := t.LoadDateRecordsExact(date)
			if err != nil {
				if err == core.ErrNoRecords {
					out.Warn("no records for %s", date.Format(util.DateFormat))
					return nil
				}
				return fmt.Errorf("failed to load records: %s", err)
			}

			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to export records: %s", err)
			}
			for _, record := range records {
				project := projects[record.Project]
				if includeArchived || !project.Archived {
					printRecord(record, project)
				}
			}
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			printColorChart()
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			err := printTags(t, includeArchived)
			if err != nil {
				return fmt.Errorf("failed to list tags: %s", err.Error())
			}
			return nil
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
		end = util.NoTimeString
	}
	dur := r.Duration(util.NoTime, util.NoTime)
	pause := r.PauseDuration(util.NoTime, util.NoTime)

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
		"%s%s %s %s %s - %s (%5s + %5s)  %s\n", name, fill,
		project.Render.Sprintf(" %s ", project.Symbol),
		date, start, end, util.FormatDuration(dur, false), util.FormatDuration(pause, false),
		note,
	)
}

func printColorChart() {
	var row, block, i uint8

	color.C256(0, true).Printf("%3d", 0)
	out.Print(" ")
	for i = 1; i < 16; i++ {
		color.S256(0, i).Printf("%3d", i)
		out.Print(" ")
	}
	out.Print("\n\n")

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
					out.Print(" ")
				}
				out.Print("  ")
			}
			out.Print("\n")
		}
		out.Print("\n")
	}

	for i = 232; i < 244; i++ {
		color.C256(i, true).Printf("%3d", i)
		out.Print(" ")
	}
	out.Print("\n")
	for i := 244; i <= 255; i++ {
		color.S256(0, uint8(i)).Printf("%3d", i)
		out.Print(" ")
	}
}

func printTags(t *core.Track, includeArchived bool) error {
	projects, err := t.LoadAllProjects()
	if err != nil {
		return err
	}

	tags := map[string]int{}
	values := map[string]map[string]bool{}

	filters := []core.FilterFunction{}
	if !includeArchived {
		filters = append(filters, core.FilterByArchived(false, projects))
	}

	fn, results, _ := t.AllRecordsFiltered(core.NewFilter(filters, util.NoTime, util.NoTime), false)

	go fn()
	for res := range results {
		if res.Err != nil {
			return res.Err
		}
		for tag, value := range res.Record.Tags {
			if v, ok := tags[tag]; ok {
				tags[tag] = v + 1
				values[tag][value] = true
			} else {
				tags[tag] = 1
				values[tag] = map[string]bool{value: true}
			}
		}
	}

	keys := maps.Keys(tags)
	sort.Strings(keys)

	for _, tag := range keys {
		v := maps.Keys(values[tag])
		sort.Strings(v)
		out.Print("%16s %4d", tag, tags[tag])
		if len(v) > 1 || (len(v) > 0 && v[0] != "") {
			out.Print(" [%s]", strings.Join(v, " "))
		}
		out.Print("\n")
	}

	return nil
}
