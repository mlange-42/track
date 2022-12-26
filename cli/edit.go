package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/fs"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	// ErrUserAbort is an error for abort by the user
	ErrUserAbort = errors.New("aborted by user")
)

const editComment string = `
# Edit the above YAML definition.
# Then, save the file and close the editor.
# 
# If you remove everything, the operation will be aborted.
`

func editCommand(t *core.Track) *cobra.Command {
	edit := &cobra.Command{
		Use:   "edit",
		Short: "Edit a resource",
		Long: `Edit a resource

Opens the resource as a temporary YAML file for editing in a text editor.
See file .track/config.yml to configure the editor to be used.`,
		Aliases: []string{"e"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	edit.AddCommand(editProjectCommand(t))
	edit.AddCommand(editRecordCommand(t))
	edit.AddCommand(editConfigCommand(t))

	edit.Long += "\n\n" + util.FormatCmdTree(edit)
	return edit
}

func editRecordCommand(t *core.Track) *cobra.Command {
	editProject := &cobra.Command{
		Use:   "record <DATE> <TIME>",
		Short: "Edit a record",
		Long: `Edit a record

Opens the record as a temporary YAML file for editing in a text editor.
See file .track/config.yml to configure the editor to be used.`,
		Aliases: []string{"r"},
		Args:    util.WrappedArgs(cobra.ExactArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			timeString := strings.Join(args, " ")
			tm, err := util.ParseDateTime(timeString)
			if err != nil {
				out.Err("failed to edit record: %s", err)
				return
			}
			err = editRecord(t, tm)
			if err != nil {
				if err == ErrUserAbort {
					out.Warn("failed to edit record: %s", err)
					return
				}
				out.Err("failed to edit record: %s", err)
				return
			}
			out.Success("Saved record '%s'", tm.Format(util.DateTimeFormat))
		},
	}

	return editProject
}

func editProjectCommand(t *core.Track) *cobra.Command {
	editProject := &cobra.Command{
		Use:   "project <NAME>",
		Short: "Edit a project",
		Long: `Edit a project

Opens the project as a temporary YAML file for editing in a text editor.
See file .track/config.yml to configure the editor to be used.`,
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			err := editProject(t, name)
			if err != nil {
				if err == ErrUserAbort {
					out.Warn("failed to edit project: %s", err)
					return
				}
				out.Err("failed to edit project: %s", err)
				return
			}
			out.Success("Saved project '%s'", name)
		},
	}

	return editProject
}

func editConfigCommand(t *core.Track) *cobra.Command {
	editProject := &cobra.Command{
		Use:     "config",
		Short:   "Edit track's config",
		Aliases: []string{"c"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			err := editConfig(t)
			if err != nil {
				if err == ErrUserAbort {
					out.Warn("failed to edit config: %s", err)
					return
				}
				out.Err("failed to edit config: %s", err)
				return
			}
			out.Success("Saved config to %s", fs.ConfigPath())
		},
	}

	return editProject
}

func editRecord(t *core.Track, tm time.Time) error {
	record, err := t.LoadRecordByTime(tm)
	if err != nil {
		return err
	}

	return edit(t, &record,
		fmt.Sprintf("# Record %s\n\n", record.Start.Format(util.DateTimeFormat)),
		func(b []byte) error {
			var newRecord core.Record
			if err := yaml.Unmarshal(b, &newRecord); err != nil {
				return err
			}

			// TODO could change, but requires deleting original file
			if newRecord.Start != record.Start {
				return fmt.Errorf("can't change start time")
			}

			if !newRecord.End.IsZero() && newRecord.End.Before(newRecord.Start) {
				return fmt.Errorf("end time is before start time")
			}

			if err = t.SaveRecord(newRecord, true); err != nil {
				return err
			}
			return nil
		})
}

func editProject(t *core.Track, name string) error {
	project, err := t.LoadProjectByName(name)
	if err != nil {
		return err
	}

	return edit(t, &project,
		fmt.Sprintf("# Project %s\n\n", project.Name),
		func(b []byte) error {
			if err := yaml.Unmarshal(b, &project); err != nil {
				return err
			}

			if project.Name != name {
				return fmt.Errorf("can't change project name")
			}

			if err = t.SaveProject(project, true); err != nil {
				return err
			}
			return nil
		})
}

func editConfig(t *core.Track) error {
	conf, err := core.LoadConfig()
	if err != nil {
		return err
	}

	return edit(t, &conf,
		fmt.Sprintf("# Track config\n\n"),
		func(b []byte) error {
			var newConfig core.Config
			if err := yaml.Unmarshal(b, &newConfig); err != nil {
				return err
			}

			if err = core.SaveConfig(newConfig); err != nil {
				return err
			}
			return nil
		})
}

func edit(t *core.Track, obj any, prefix string, fn func(b []byte) error) error {
	file, err := os.CreateTemp("", "track-*.yml")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	bytes, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = file.WriteString(prefix)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	_, err = file.WriteString(editComment)
	if err != nil {
		return err
	}

	file.Close()

	err = fs.EditFile(file.Name(), t.Config.TextEditor)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return ErrUserAbort
	}

	if err := fn(content); err != nil {
		return err
	}

	return nil
}
