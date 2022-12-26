package cli

import (
	"strings"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func deleteCommand(t *core.Track) *cobra.Command {
	edit := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a resource",
		Long:    `Delete a resource`,
		Aliases: []string{"e"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	edit.AddCommand(deleteRecordCommand(t))

	edit.Long += "\n\n" + formatCmdTree(edit)
	return edit
}

func deleteRecordCommand(t *core.Track) *cobra.Command {
	var force bool

	delete := &cobra.Command{
		Use:     "record <DATE> <TIME>",
		Short:   "Delete a record",
		Long:    "Delete a record",
		Aliases: []string{"r"},
		Args:    util.WrappedArgs(cobra.ExactArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			open, ok := t.OpenRecord()
			if !ok {
				out.Err("failed to stop record: no record running")
				return
			}

			timeString := strings.Join(args, " ")
			tm, err := util.ParseDateTime(timeString)
			if err != nil {
				out.Err("failed to delete record: %s", err)
				return
			}
			record, err := t.LoadRecordByTime(tm)
			if err != nil {
				out.Err("failed to delete record: %s", err)
				return
			}

			if !force && !confirmDeleteRecord(open) {
				out.Err("failed to delete record: aborted by user")
				return
			}

			err = t.DeleteRecord(record)
			if err != nil {
				out.Err("failed to delete record: %s", err)
				return
			}
			out.Success("Deleted record %s from '%s'", record.Start.Format(util.DateTimeFormat), record.Project)
		},
	}

	delete.Flags().BoolVarP(&force, "force", "F", false, "Don't prompt for confirmation.")

	return delete
}
