package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func deleteCommand(t *core.Track) *cobra.Command {
	var dryRun bool

	delete := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a resource",
		Long:    `Delete a resource`,
		Aliases: []string{"D"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	delete.PersistentFlags().BoolVar(&dryRun, "dry", false, "Dry run: do not actually change any files")

	delete.AddCommand(deleteRecordCommand(t, &dryRun))
	delete.AddCommand(deleteProjectCommand(t, &dryRun))

	delete.Long += "\n\n" + formatCmdTree(delete)
	return delete
}

func deleteRecordCommand(t *core.Track, dryRun *bool) *cobra.Command {
	var force bool

	delete := &cobra.Command{
		Use:     "record DATE TIME",
		Short:   "Delete a record",
		Long:    "Delete a record",
		Aliases: []string{"r"},
		Args:    util.WrappedArgs(cobra.ExactArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			timeString := strings.Join(args, " ")
			tm, err := util.ParseDateTime(timeString)
			if err != nil {
				out.Err("failed to delete record: %s", err)
				return
			}
			record, err := t.LoadRecord(tm)
			if err != nil {
				out.Err("failed to delete record: %s", err)
				return
			}

			if !force && !confirm(
				fmt.Sprintf(
					"Really delete record %s (%s) from project '%s' (y/n): ",
					record.Start.Format(util.DateTimeFormat),
					util.FormatDuration(record.Duration(time.Time{}, time.Time{})),
					record.Project,
				),
				"y",
			) {
				out.Warn("failed to delete record: aborted by user")
				return
			}

			if *dryRun {
				out.Success("Deleted record %s from '%s' - dry-run", record.Start.Format(util.DateTimeFormat), record.Project)
			} else {
				err = t.DeleteRecord(&record)
				if err != nil {
					out.Err("failed to delete record: %s", err)
					return
				}
				out.Success("Deleted record %s from '%s'", record.Start.Format(util.DateTimeFormat), record.Project)
			}
		},
	}

	delete.Flags().BoolVarP(&force, "force", "F", false, "Don't prompt for confirmation.")

	return delete
}

func deleteProjectCommand(t *core.Track, dryRun *bool) *cobra.Command {
	var force bool

	delete := &cobra.Command{
		Use:     "project PROJECT",
		Short:   "Delete a project and all associated records",
		Long:    "Delete a project and all associated records",
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to delete project: %s", err)
				return
			}
			pTree, err := t.ToProjectTree(projects)
			if err != nil {
				out.Err("failed to delete project: %s", err)
				return
			}

			pNode, ok := pTree.Nodes[name]
			if !ok {
				out.Err("failed to delete project: no project named '%s'", name)
				return
			}
			if len(pNode.Children) > 0 {
				out.Err("failed to delete project: '%s' has %d child project(s)", name, len(pNode.Children))
				return
			}

			if !force && !confirm(
				fmt.Sprintf(
					"Really delete project '%s' and all associated records? (yes!/n): ",
					pNode.Value.Name,
				),
				"yes!",
			) {
				out.Warn("failed to delete project: aborted by user")
				return
			}

			cnt, err := t.DeleteProject(&pNode.Value, true, *dryRun)
			if err != nil {
				out.Err("failed to delete project: %s", err)
				out.Err("deleted %d records", cnt)
				return
			}
			if *dryRun {
				out.Success("Deleted project '%s' (%d records) - dry-run", pNode.Value.Name, cnt)
			} else {
				out.Success("Deleted project '%s' (%d records)", pNode.Value.Name, cnt)
			}
		},
	}

	delete.Flags().BoolVarP(&force, "force", "F", false, "Don't prompt for confirmation.")

	return delete
}
