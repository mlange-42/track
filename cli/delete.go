package cli

import (
	"fmt"
	"strings"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			timeString := strings.Join(args, " ")
			tm, err := util.ParseDateTime(timeString)
			if err != nil {
				return fmt.Errorf("failed to delete record: %s", err)
			}
			record, err := t.LoadRecord(tm)
			if err != nil {
				return fmt.Errorf("failed to delete record: %s", err)
			}

			if !force && !confirm(
				fmt.Sprintf(
					"Really delete record %s (%s) from project '%s' (y/n): ",
					record.Start.Format(util.DateTimeFormat),
					util.FormatDuration(record.Duration(util.NoTime, util.NoTime)),
					record.Project,
				),
				"y",
			) {
				return fmt.Errorf("failed to delete record: aborted by user")
			}

			if *dryRun {
				out.Success("Deleted record %s from '%s' - dry-run", record.Start.Format(util.DateTimeFormat), record.Project)
			} else {
				err = t.DeleteRecord(&record)
				if err != nil {
					return fmt.Errorf("failed to delete record: %s", err)
				}
				out.Success("Deleted record %s from '%s'", record.Start.Format(util.DateTimeFormat), record.Project)
			}
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to delete project: %s", err)
			}
			pTree, err := t.ToProjectTree(projects)
			if err != nil {
				return fmt.Errorf("failed to delete project: %s", err)
			}

			pNode, ok := pTree.Nodes[name]
			if !ok {
				return fmt.Errorf("failed to delete project: no project named '%s'", name)
			}
			if len(pNode.Children) > 0 {
				return fmt.Errorf("failed to delete project: '%s' has %d child project(s)", name, len(pNode.Children))
			}

			if !force && !confirm(
				fmt.Sprintf(
					"Really delete project '%s' and all associated records? (yes!/n): ",
					pNode.Value.Name,
				),
				"yes!",
			) {
				return fmt.Errorf("failed to delete project: aborted by user")
			}

			cnt, err := t.DeleteProject(&pNode.Value, true, *dryRun)
			if err != nil {
				return fmt.Errorf("failed to delete project: %s (deleted %d records)", err, cnt)
			}
			if *dryRun {
				out.Success("Deleted project '%s' (%d records) - dry-run", pNode.Value.Name, cnt)
			} else {
				out.Success("Deleted project '%s' (%d records)", pNode.Value.Name, cnt)
			}
			return nil
		},
	}

	delete.Flags().BoolVarP(&force, "force", "F", false, "Don't prompt for confirmation.")

	return delete
}
