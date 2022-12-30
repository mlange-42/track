package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode/utf8"

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
%[1]s Edit the above definition.
%[1]s Then, save the file and close the editor.
%[1]s 
%[1]s If you remove everything, the operation will be aborted.
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
	edit.AddCommand(editDayCommand(t))
	edit.AddCommand(editConfigCommand(t))

	edit.Long += "\n\n" + formatCmdTree(edit)
	return edit
}

func editRecordCommand(t *core.Track) *cobra.Command {
	editProject := &cobra.Command{
		Use:   "record [[DATE] TIME]",
		Short: "Edit a record",
		Long: `Edit a record

Opens the record as a temporary file for editing.
See file .track/config.yml to configure the editor to be used.

Edits the last or open record if no date and time are given.

Uses the current date if only a time is given.`,
		Aliases: []string{"r"},
		Args:    util.WrappedArgs(cobra.MaximumNArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			tm := time.Time{}
			switch len(args) {
			case 0:
				last, err := t.LatestRecord()
				if err != nil {
					out.Err("failed to edit record: %s", err)
					return
				}
				tm = last.Start
			case 1:
				tm, err = time.ParseInLocation(util.TimeFormat, args[0], time.Local)
				if err != nil {
					out.Err("failed to edit record: %s", err)
					return
				}
				tm = util.DateAndTime(time.Now(), tm)
			case 2:
				timeString := strings.Join(args, " ")
				tm, err = util.ParseDateTime(timeString)
				if err != nil {
					out.Err("failed to edit record: %s", err)
					return
				}
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
	var archive bool

	editProject := &cobra.Command{
		Use:   "project PROJECT",
		Short: "Edit a project",
		Long: `Edit a project

Opens the project as a temporary YAML file for editing if no flags are given.
See file .track/config.yml to configure the editor to be used.`,
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			project, err := t.LoadProjectByName(name)
			if err != nil {
				out.Err("failed to edit project: %s", err)
				return
			}

			changed := false
			if cmd.Flags().Changed("archive") {
				if project.Archived == archive {
					out.Warn("New value for 'archive' equals old value\n")
				} else {
					project.Archived = archive
					if archive {
						out.Success("Archived project '%s'\n", project.Name)
					} else {
						out.Success("Un-archived project '%s'\n", project.Name)
					}
				}
				changed = true
			}

			if changed {
				if err := t.SaveProject(project, true); err != nil {
					out.Err("failed to edit project: %s", err)
					return
				}
			} else {
				err = editProject(t, project)
				if err != nil {
					if err == ErrUserAbort {
						out.Warn("failed to edit project: %s", err)
						return
					}
					out.Err("failed to edit project: %s", err)
					return
				}
			}
			out.Success("Saved project '%s'", name)
		},
	}
	editProject.Flags().BoolVarP(&archive, "archive", "a", false, "Archive or un-archive a project. Use like '-a=false'")

	return editProject
}

func editConfigCommand(t *core.Track) *cobra.Command {
	editProject := &cobra.Command{
		Use:   "config",
		Short: "Edit track's config",
		Long: `Edit track's config

Opens the config as a temporary YAML file for editing.
See file .track/config.yml to configure the editor to be used.`,
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

func editDayCommand(t *core.Track) *cobra.Command {
	var dryRun bool

	editDay := &cobra.Command{
		Use:   "day [DATE]",
		Short: "Edit all records of one day",
		Long: `Edit all records of one day

Opens the records in a single temporary file for editing.
See file .track/config.yml to configure the editor to be used.`,
		Aliases: []string{"d"},
		Args:    util.WrappedArgs(cobra.MaximumNArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			date := util.ToDate(time.Now())
			if len(args) > 0 {
				date, err = util.ParseDate(args[0])
				if err != nil {
					out.Err("failed to edit day: %s", err)
					return
				}
			}
			err = editDay(t, date, dryRun)
			if err != nil {
				if err == ErrUserAbort {
					out.Warn("failed to edit day: %s", err)
					return
				}
				out.Err("failed to edit day: %s", err)
				return
			}
			if dryRun {
				out.Success("Saved day records - dry-run")
			} else {
				out.Success("Saved day records")
			}
		},
	}
	editDay.Flags().BoolVar(&dryRun, "dry", false, "Dry run: do not actually change any files")

	return editDay
}
func editRecord(t *core.Track, tm time.Time) error {
	record, err := t.LoadRecord(tm)
	if err != nil {
		return err
	}

	return edit(t, &record,
		fmt.Sprintf("%s Record %s\n\n", core.CommentPrefix, record.Start.Format(util.DateTimeFormat)),
		core.CommentPrefix,
		func(r *core.Record) ([]byte, error) {
			str := r.Serialize(time.Time{})
			return []byte(str), nil
		},
		func(b []byte) error {
			newRecord, err := core.DeserializeRecord(string(b), record.Start)
			if err != nil {
				return err
			}

			if newRecord.Start != record.Start {
				return fmt.Errorf("can't change start time. Try command 'track edit day' instead")
			}

			if err = newRecord.Check(); err != nil {
				return err
			}

			if err = t.SaveRecord(&newRecord, true); err != nil {
				return err
			}
			return nil
		})
}

func editDay(t *core.Track, date time.Time, dryRun bool) error {
	date = util.ToDate(date)
	dateBefore := date.Add(-24 * time.Hour)
	dateAfter := date.Add(24 * time.Hour)

	filters := core.FilterFunctions{
		core.FilterByTime(date, dateAfter),
	}

	records, err := t.LoadDateRecordsFiltered(dateBefore, filters)
	if err != nil && !errors.Is(err, core.ErrNoRecords) {
		return err
	}
	records2, err := t.LoadDateRecordsFiltered(date, filters)
	if err != nil && !errors.Is(err, core.ErrNoRecords) {
		return err
	}
	records = append(records, records2...)

	if len(records) == 0 {
		return fmt.Errorf("no records for %s", date.Format(util.DateFormat))
	}

	return edit(t, records,
		fmt.Sprintf("%[1]s Records for %s\n%[1]s Clear file to abort\n\n", core.CommentPrefix, date.Format(util.DateFormat)),
		core.CommentPrefix,
		func(records []core.Record) ([]byte, error) {
			str := ""
			for i, rec := range records {
				str += rec.Serialize(date)
				if i < len(records)-1 {
					str += "\n--------------------\n\n"
				}
			}
			return []byte(str), nil
		},
		func(b []byte) error {
			lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
			prevIdx := 0

			newRecords := []core.Record{}

			for i, line := range lines {
				if strings.HasPrefix(line, "----") || i == len(lines)-1 {
					endIdx := i
					if i == len(lines)-1 {
						endIdx = len(lines)
					}
					str := strings.Join(lines[prevIdx:endIdx], "\n")
					rec, err := core.DeserializeRecord(str, date)
					if err != nil {
						return err
					}
					newRecords = append(newRecords, rec)
					prevIdx = i + 1
				}
			}

			now := time.Now()
			today := util.ToDate(now)

			if len(newRecords) > 0 {
				oldFirst, newFirst := records[0], newRecords[0]
				if oldFirst.Start.Before(date) {
					if newFirst.Start.Before(oldFirst.Start) {
						return fmt.Errorf(
							"can't extend a start time on the day before. try 'track edit day %s'",
							dateBefore.Format(util.DateFormat),
						)
					}
				} else {
					if newFirst.Start.Before(date) {
						return fmt.Errorf(
							"can't move a start time to the day before. try 'track edit day %s'",
							dateBefore.Format(util.DateFormat),
						)
					}
				}
				oldLast, newLast := records[len(records)-1], newRecords[len(newRecords)-1]
				if newLast.Start.After(now) || newLast.End.After(now) {
					return fmt.Errorf("can't date into the future")
				}
				if oldLast.End.After(dateAfter) {
					if !newLast.End.IsZero() && newLast.End.After(oldLast.End) {
						return fmt.Errorf(
							"can't extend an end time on the day after. try 'track edit day %s'",
							dateAfter.Format(util.DateFormat),
						)
					}
				}

				prevStart := time.Time{}
				prevEnd := time.Time{}

				for i, rec := range newRecords {
					if rec.Start.Before(prevStart) {
						return fmt.Errorf("records are not in chronological order")
					}
					if rec.Start.Before(prevEnd) {
						return fmt.Errorf("records overlap (%s / %s)", prevStart.Format(util.TimeFormat), rec.Start.Format(util.TimeFormat))
					}
					if rec.End.IsZero() {
						if i != len(newRecords)-1 {
							return fmt.Errorf("only the last record can have an open end time")
						}
						if !oldLast.End.IsZero() && date != today {
							return fmt.Errorf(
								"can't set open end for record starting on another day. Try 'track edit day today'",
							)
						}
					}
					prevStart = rec.Start
					prevEnd = rec.End
				}
			}

			if !dryRun {
				for _, rec := range records {
					t.DeleteRecord(&rec)
				}

				for _, rec := range newRecords {
					t.SaveRecord(&rec, false)
				}
			}

			return nil
		})
}

func editProject(t *core.Track, project core.Project) error {
	return edit(t, &project,
		fmt.Sprintf("%s Project %s\n\n", core.YamlCommentPrefix, project.Name),
		core.YamlCommentPrefix,
		func(r *core.Project) ([]byte, error) {
			return yaml.Marshal(r)
		},
		func(b []byte) error {
			var newProject core.Project
			if err := yaml.Unmarshal(b, &newProject); err != nil {
				return err
			}

			if newProject.Name != project.Name {
				return fmt.Errorf("can't change project name")
			}
			if utf8.RuneCountInString(newProject.Symbol) != 1 {
				return fmt.Errorf("symbol must be a single character")
			}
			if err := t.CheckParents(newProject); err != nil {
				return err
			}
			if err := t.SaveProject(newProject, true); err != nil {
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
		fmt.Sprintf("%s Track config\n\n", core.YamlCommentPrefix),
		core.YamlCommentPrefix,
		func(r *core.Config) ([]byte, error) {
			return yaml.Marshal(r)
		},
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

func edit[T any](t *core.Track, obj T, comment string, commentPrefix string, marshal func(T) ([]byte, error), unmarshal func(b []byte) error) error {
	file, err := os.CreateTemp("", "track-*.yml")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	bytes, err := marshal(obj)
	if err != nil {
		return err
	}

	_, err = file.WriteString(comment)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf(editComment, commentPrefix))
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

	if err := unmarshal(content); err != nil {
		return err
	}

	return nil
}
