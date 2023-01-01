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
	var dryRun bool

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

	edit.PersistentFlags().BoolVar(&dryRun, "dry", false, "Dry run: do not actually change any files")

	edit.AddCommand(editProjectCommand(t, &dryRun))
	edit.AddCommand(editRecordCommand(t, &dryRun))
	edit.AddCommand(editDayCommand(t, &dryRun))
	edit.AddCommand(editConfigCommand(t, &dryRun))

	edit.Long += "\n\n" + formatCmdTree(edit)
	return edit
}

func editRecordCommand(t *core.Track, dryRun *bool) *cobra.Command {

	editRecord := &cobra.Command{
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
				date, err := util.ParseDate(args[0])
				if err != nil {
					out.Err("failed to edit record: %s", err)
					return
				}
				tm, err = time.Parse(util.TimeFormat, args[1])
				if err != nil {
					out.Err("failed to edit record: %s", err)
					return
				}
				tm = util.DateAndTime(date, tm)
			}

			err = editRecord(t, tm, *dryRun)
			if err != nil {
				if err == ErrUserAbort {
					out.Warn("failed to edit record %s: %s", tm.Format(util.DateTimeFormat), err)
					return
				}
				out.Err("failed to edit record %s: %s", tm.Format(util.DateTimeFormat), err)
				return
			}
			if *dryRun {
				out.Success("Saved record %s - dry-run", tm.Format(util.DateTimeFormat))
			} else {
				out.Success("Saved record %s", tm.Format(util.DateTimeFormat))
			}
		},
	}

	return editRecord
}

func editProjectCommand(t *core.Track, dryRun *bool) *cobra.Command {
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
				err = editProject(t, project, *dryRun)
				if err != nil {
					if err == ErrUserAbort {
						out.Warn("failed to edit project: %s", err)
						return
					}
					out.Err("failed to edit project: %s", err)
					return
				}
			}
			if *dryRun {
				out.Success("Saved project %s - dry-run", name)
			} else {
				out.Success("Saved project %s", name)
			}
		},
	}
	editProject.Flags().BoolVarP(&archive, "archive", "a", false, "Archive or un-archive a project. Use like '-a=false'")

	return editProject
}

func editConfigCommand(t *core.Track, dryRun *bool) *cobra.Command {

	editConfig := &cobra.Command{
		Use:   "config",
		Short: "Edit track's config",
		Long: `Edit track's config

Opens the config as a temporary YAML file for editing.
See file .track/config.yml to configure the editor to be used.`,
		Aliases: []string{"c"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			err := editConfig(t, *dryRun)
			if err != nil {
				if err == ErrUserAbort {
					out.Warn("failed to edit config: %s", err)
					return
				}
				out.Err("failed to edit config: %s", err)
				return
			}

			if *dryRun {
				out.Success("Saved config to %s - dry-run", fs.ConfigPath())
			} else {
				out.Success("Saved config to %s", fs.ConfigPath())
			}
		},
	}

	return editConfig
}

func editDayCommand(t *core.Track, dryRun *bool) *cobra.Command {

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
			err = editDay(t, date, *dryRun)
			if err != nil {
				if err == ErrUserAbort {
					out.Warn("failed to edit day: %s", err)
					return
				}
				out.Err("failed to edit day: %s", err)
				return
			}
			if *dryRun {
				out.Success("Saved day records - dry-run")
			} else {
				out.Success("Saved day records")
			}
		},
	}

	return editDay
}
func editRecord(t *core.Track, tm time.Time, dryRun bool) error {
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
			if record.End.IsZero() {
				if !newRecord.End.IsZero() && newRecord.End.After(time.Now()) {
					return fmt.Errorf("can't set end time to the future. Try command 'track edit day' instead")
				}
			} else {
				if newRecord.End.IsZero() {
					return fmt.Errorf("can't open a finished record. Try command 'track edit day' instead")
				}
				if newRecord.End.After(record.End) {
					return fmt.Errorf("can't extend record end time. Try command 'track edit day' instead")
				}
			}

			if err = newRecord.Check(); err != nil {
				return err
			}

			if !dryRun {
				if err = t.SaveRecord(&newRecord, true); err != nil {
					return err
				}
			}
			return nil
		})
}

func editDay(t *core.Track, date time.Time, dryRun bool) error {
	date = util.ToDate(date)
	dateBefore := date.Add(-24 * time.Hour)
	dateAfter := date.Add(24 * time.Hour)

	records, err := t.LoadDateRecordsExact(date)
	if err != nil {
		if errors.Is(err, core.ErrNoRecords) {
			return fmt.Errorf("no records for %s", date.Format(util.DateFormat))
		}
		return err
	}

	projects, err := t.LoadAllProjects()
	if err != nil {
		return err
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
							"can't extend a start time on the day before (%s / %s). Try 'track edit day %s'",
							newFirst.Start.Format(util.TimeFormat),
							oldFirst.Start.Format(util.TimeFormat),
							dateBefore.Format(util.DateFormat),
						)
					}
				} else {
					if newFirst.Start.Before(date) {
						return fmt.Errorf(
							"can't move a start time to the day before (%s). Try 'track edit day %s'",
							newFirst.Start.Format(util.TimeFormat),
							dateBefore.Format(util.DateFormat),
						)
					}
				}
				oldLast, newLast := records[len(records)-1], newRecords[len(newRecords)-1]
				if newLast.Start.After(now) || newLast.End.After(now) {
					return fmt.Errorf("can't date into the future (%s)", newLast.Start.Format(util.TimeFormat))
				}
				if oldLast.End.After(dateAfter) {
					if !newLast.End.IsZero() && newLast.End.After(oldLast.End) {
						return fmt.Errorf(
							"can't extend an end time on the day after (%s). Try 'track edit day %s'",
							newLast.Start.Format(util.TimeFormat),
							dateAfter.Format(util.DateFormat),
						)
					}
				}

				prevStart := time.Time{}
				prevEnd := time.Time{}

				for i, rec := range newRecords {
					if _, ok := projects[rec.Project]; !ok {
						return fmt.Errorf("project '%s' does not exist (%s)", rec.Project, rec.Start.Format(util.TimeFormat))
					}
					if rec.Start.Before(prevStart) {
						return fmt.Errorf(
							"records are not in chronological order (%s / %s)",
							prevStart.Format(util.TimeFormat),
							rec.Start.Format(util.TimeFormat),
						)
					}
					if rec.Start.Before(prevEnd) {
						return fmt.Errorf("records overlap (%s / %s)", prevStart.Format(util.TimeFormat), rec.Start.Format(util.TimeFormat))
					}
					if rec.End.IsZero() {
						if i != len(newRecords)-1 {
							return fmt.Errorf("only the last record can have an open end time (%s)", rec.Start.Format(util.TimeFormat))
						}
						if !oldLast.End.IsZero() && date != today {
							return fmt.Errorf(
								"can't set open end for record starting on another day (%s). Try 'track edit day today'",
								rec.Start.Format(util.TimeFormat),
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

func editProject(t *core.Track, project core.Project, dryRun bool) error {
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

			if !dryRun {
				if err := t.SaveProject(newProject, true); err != nil {
					return err
				}
			}
			return nil
		})
}

func editConfig(t *core.Track, dryRun bool) error {
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

			if !dryRun {
				if err = core.SaveConfig(newConfig); err != nil {
					return err
				}
			}
			return nil
		})
}

func edit[T any](t *core.Track, obj T, comment string, commentPrefix string, marshal func(T) ([]byte, error), unmarshal func(b []byte) error) error {
	content, err := marshal(obj)
	if err != nil {
		return err
	}

	firstTrial := true
	errorComment := ""
	for {
		file, err := os.CreateTemp("", "track-*.yml")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		if firstTrial {
			_, err = file.WriteString(comment)
			if err != nil {
				return err
			}
		} else {
			_, err = file.WriteString(fmt.Sprintf("%s ERROR: %s\n", commentPrefix, errorComment))
			if err != nil {
				return err
			}
		}

		_, err = file.Write(content)
		if err != nil {
			return err
		}

		if firstTrial {
			_, err = file.WriteString(fmt.Sprintf(editComment, commentPrefix))
			if err != nil {
				return err
			}
		}

		file.Close()

		err = fs.EditFile(file.Name(), t.Config.TextEditor)
		if err != nil {
			return err
		}

		content, err = ioutil.ReadFile(file.Name())
		if err != nil {
			return err
		}

		if len(content) == 0 {
			return ErrUserAbort
		}

		firstTrial = false
		if err := unmarshal(content); err != nil {
			out.Err("%s\n", err.Error())
			errorComment = err.Error()
			continue
		}

		break
	}

	return nil
}
