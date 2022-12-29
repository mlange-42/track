package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mlange-42/track/fs"
	"github.com/mlange-42/track/util"
)

// TagPrefix denotes tags in record notes
const TagPrefix = "+"

var (
	// ErrNoRecords is an error for no records found for a date
	ErrNoRecords = errors.New("no records for date")
)

// Record holds and manipulates data for a record
type Record struct {
	Project string
	Start   time.Time
	End     time.Time
	Note    string
	Tags    []string
	Pause   []Pause
}

// Pause holds information about a pause in a record
type Pause struct {
	Start time.Time
	End   time.Time
}

// HasEnded reports whether the record has an end time
func (r Record) HasEnded() bool {
	return !r.End.IsZero()
}

// Duration reports the duration of a record
func (r Record) Duration(min, max time.Time) time.Duration {
	dur := util.DurationClip(r.Start, r.End, min, max)
	dur -= r.PauseDuration(min, max)
	return dur
}

// PauseDuration reports the duration of all pauses of the record
func (r Record) PauseDuration(min, max time.Time) time.Duration {
	dur := time.Second * 0
	for _, p := range r.Pause {
		dur += util.DurationClip(p.Start, p.End, min, max)
	}
	return dur
}

// Serialize converts a record to a serialization string
func (r Record) Serialize() string {
	endTime := "?"
	if !r.End.IsZero() {
		endTime = r.End.Format(util.TimeFormat)
		if util.ToDate(r.End).After(r.Start) {
			endTime = "+" + endTime
		}
	}
	res := fmt.Sprintf("%s - %s", r.Start.Format(util.TimeFormat), endTime)
	for _, p := range r.Pause {
		endTime := "?"
		if !p.End.IsZero() {
			endTime = p.End.Sub(p.Start).String()
		}
		res += fmt.Sprintf("\n    - %s - %s", p.Start.Format(util.TimeFormat), endTime)
	}
	res += fmt.Sprintf("\n\n    %s", r.Project)

	if len(r.Note) > 0 {
		res += fmt.Sprintf("\n\n%s", r.Note)
	}
	return res
}

// DeserializeRecord converts a serialization string to a record
func DeserializeRecord(str string, date time.Time) (Record, error) {
	str = strings.TrimSpace(str)
	lines := strings.Split(strings.ReplaceAll(str, "\r\n", "\n"), "\n")
	index, ok := skipLines(lines, 0, true)
	if !ok {
		return Record{}, fmt.Errorf("invalid record: missing time range (1st line)")
	}
	start, end, err := util.ParseTimeRange(lines[index], date)
	index++
	if err != nil {
		return Record{}, err
	}

	pause := []Pause{}
	for {
		ln := strings.TrimSpace(lines[index])
		if strings.HasPrefix(ln, "- ") {
			pStart, pEnd, err := util.ParseTimeRange(strings.TrimPrefix(ln, "- "), date)
			index++
			if err != nil {
				return Record{}, err
			}
			pause = append(pause, Pause{Start: pStart, End: pEnd})
		} else {
			break
		}
	}

	index, ok = skipLines(lines, index, true)
	if !ok {
		return Record{}, fmt.Errorf("invalid record: missing project (2nd line)")
	}
	project := strings.TrimSpace(lines[index])
	index++

	note := ""
	tags := []string{}
	index, ok = skipLines(lines, index, true)
	if ok {
		note = strings.TrimSpace(strings.Join(lines[index:], "\n"))
		tags = ExtractTagsSlice(lines[index:])
	}
	return Record{
		Project: project,
		Start:   start,
		End:     end,
		Note:    note,
		Tags:    tags,
		Pause:   pause,
	}, nil
}

func skipLines(lines []string, index int, skipEmpty bool) (int, bool) {
	if index >= len(lines) {
		return index, false
	}
	for (skipEmpty && strings.TrimSpace(lines[index]) == "") || strings.HasPrefix(lines[index], "#") {
		index++
		if index >= len(lines) {
			return index, false
		}
	}
	return index, true
}

// RecordsDir returns the records storage directory
func (t *Track) RecordsDir() string {
	return filepath.Join(fs.RootDir(), t.Workspace(), fs.RecordsDirName())
}

// WorkspaceRecordsDir returns the records storage directory for the given workspace
func (t *Track) WorkspaceRecordsDir(ws string) string {
	return filepath.Join(fs.RootDir(), ws, fs.RecordsDirName())
}

// RecordPath returns the full path for a record
func (t *Track) RecordPath(tm time.Time) string {
	return filepath.Join(
		t.RecordDir(tm),
		fmt.Sprintf("%s.trk", tm.Format(util.FileTimeFormat)),
	)
}

// RecordDir returns the directory path for a record
func (t *Track) RecordDir(tm time.Time) string {
	return filepath.Join(
		t.RecordsDir(),
		fmt.Sprintf("%04d", tm.Year()),
		fmt.Sprintf("%02d", int(tm.Month())),
		fmt.Sprintf("%02d", tm.Day()),
	)
}

// SaveRecord saves a record to disk
func (t *Track) SaveRecord(record *Record, force bool) error {
	path := t.RecordPath(record.Start)
	if !force && fs.FileExists(path) {
		return fmt.Errorf("record already exists")
	}
	dir := t.RecordDir(record.Start)
	err := fs.CreateDir(dir)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer file.Close()

	if err != nil {
		return err
	}

	bytes := record.Serialize()

	_, err = fmt.Fprintf(file, "# Record %s\n", record.Start.Format(util.DateTimeFormat))
	if err != nil {
		return err
	}

	_, err = file.WriteString(bytes)

	return err
}

// DeleteRecord deletes a record
func (t *Track) DeleteRecord(record *Record) error {
	path := t.RecordPath(record.Start)
	if !fs.FileExists(path) {
		return fmt.Errorf("record does not exist")
	}
	err := os.Remove(path)
	if err != nil {
		return err
	}
	dayDir := filepath.Dir(path)
	empty, err := fs.DirIsEmpty(dayDir)
	if err != nil {
		return err
	}
	if empty {
		os.Remove(dayDir)
		monthDir := filepath.Dir(dayDir)
		empty, err := fs.DirIsEmpty(monthDir)
		if err != nil {
			return err
		}
		if empty {
			os.Remove(monthDir)
			yearDir := filepath.Dir(monthDir)
			empty, err := fs.DirIsEmpty(yearDir)
			if err != nil {
				return err
			}
			if empty {
				os.Remove(yearDir)

			}
		}
	}
	return nil
}

// LoadRecord loads a record
func (t *Track) LoadRecord(tm time.Time) (Record, error) {
	path := t.RecordPath(tm)
	file, err := os.ReadFile(path)
	if err != nil {
		return Record{}, err
	}

	record, err := DeserializeRecord(string(file), tm)
	if err != nil {
		return Record{}, err
	}

	return record, nil
}

// LoadAllRecords loads all records
func (t *Track) LoadAllRecords() ([]Record, error) {
	return t.LoadAllRecordsFiltered([]func(*Record) bool{})
}

// LoadAllRecordsFiltered loads all records
func (t *Track) LoadAllRecordsFiltered(filters FilterFunctions) ([]Record, error) {
	path := t.RecordsDir()

	var records []Record

	yearDirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, yearDir := range yearDirs {
		if !yearDir.IsDir() {
			continue
		}
		year, err := strconv.Atoi(yearDir.Name())
		if err != nil {
			return nil, err
		}
		monthDirs, err := ioutil.ReadDir(filepath.Join(path, yearDir.Name()))
		if err != nil {
			return nil, err
		}

		for _, monthDir := range monthDirs {
			if !monthDir.IsDir() {
				continue
			}
			month, err := strconv.Atoi(monthDir.Name())
			if err != nil {
				return nil, err
			}

			dayDirs, err := ioutil.ReadDir(filepath.Join(path, yearDir.Name(), monthDir.Name()))
			if err != nil {
				return nil, err
			}

			for _, dayDir := range dayDirs {
				if !dayDir.IsDir() {
					continue
				}
				day, err := strconv.Atoi(dayDir.Name())
				if err != nil {
					return nil, err
				}

				recs, err := t.LoadDateRecordsFiltered(util.Date(year, time.Month(month), day), filters)
				if err != nil {
					return nil, err
				}
				records = append(records, recs...)
			}
		}
	}

	return records, nil
}

// FilterResult contains a Report or an error from async filtering
type FilterResult struct {
	Record Record
	Err    error
}

// AllRecordsFiltered is an async version of LoadAllRecordsFiltered
func (t *Track) AllRecordsFiltered(filters FilterFunctions) (func(), chan FilterResult) {
	results := make(chan FilterResult, 32)

	return func() {
		path := t.RecordsDir()

		yearDirs, err := ioutil.ReadDir(path)
		if err != nil {
			results <- FilterResult{Record{}, err}
			return
		}

		for _, yearDir := range yearDirs {
			if !yearDir.IsDir() {
				continue
			}
			year, err := strconv.Atoi(yearDir.Name())
			if err != nil {
				results <- FilterResult{Record{}, err}
				return
			}
			monthDirs, err := ioutil.ReadDir(filepath.Join(path, yearDir.Name()))
			if err != nil {
				results <- FilterResult{Record{}, err}
				return
			}

			for _, monthDir := range monthDirs {
				if !monthDir.IsDir() {
					continue
				}
				month, err := strconv.Atoi(monthDir.Name())
				if err != nil {
					results <- FilterResult{Record{}, err}
					return
				}

				dayDirs, err := ioutil.ReadDir(filepath.Join(path, yearDir.Name(), monthDir.Name()))
				if err != nil {
					results <- FilterResult{Record{}, err}
					return
				}

				for _, dayDir := range dayDirs {
					if !dayDir.IsDir() {
						continue
					}
					day, err := strconv.Atoi(dayDir.Name())
					if err != nil {
						results <- FilterResult{Record{}, err}
						return
					}

					recs, err := t.LoadDateRecordsFiltered(util.Date(year, time.Month(month), day), filters)
					if err != nil {
						results <- FilterResult{Record{}, err}
						return
					}
					for _, rec := range recs {
						results <- FilterResult{rec, nil}
					}
				}
			}
		}
		close(results)
	}, results
}

// LoadDateRecords loads all records for the given date string/directory
func (t *Track) LoadDateRecords(date time.Time) ([]Record, error) {
	return t.LoadDateRecordsFiltered(date, []func(*Record) bool{})
}

// LoadDateRecordsFiltered loads all records for the given date string/directory
func (t *Track) LoadDateRecordsFiltered(date time.Time, filters FilterFunctions) ([]Record, error) {
	subPath := t.RecordDir(date)

	info, err := os.Stat(subPath)
	if err != nil {
		return nil, ErrNoRecords
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", info.Name())
	}

	var records []Record

	files, err := ioutil.ReadDir(subPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		tm, err := fileToTime(date, file.Name())
		record, err := t.LoadRecord(tm)
		if err != nil {
			return nil, err
		}
		if Filter(&record, filters) {
			records = append(records, record)
		}
	}

	return records, nil
}

// LatestRecord loads the latest record. Returns a nil reference if no record is found.
func (t *Track) LatestRecord() (*Record, error) {
	records := t.RecordsDir()
	yearPath, year, err := fs.FindLatests(records, true)
	if err != nil {
		if errors.Is(err, fs.ErrNoFiles) {
			return nil, nil
		}
		return nil, err
	}
	monthPath, month, err := fs.FindLatests(yearPath, true)
	if err != nil {
		if errors.Is(err, fs.ErrNoFiles) {
			return nil, nil
		}
		return nil, err
	}
	dayPath, day, err := fs.FindLatests(monthPath, true)
	if err != nil {
		if errors.Is(err, fs.ErrNoFiles) {
			return nil, nil
		}
		return nil, err
	}
	_, record, err := fs.FindLatests(dayPath, false)
	if err != nil {
		if errors.Is(err, fs.ErrNoFiles) {
			return nil, nil
		}
		return nil, err
	}

	tm, err := pathToTime(year, month, day, record)
	if err != nil {
		return nil, err
	}
	rec, err := t.LoadRecord(tm)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}

func pathToTime(y, m, d, file string) (time.Time, error) {
	return time.ParseInLocation(
		util.FileDateTimeFormat,
		fmt.Sprintf("%s-%s-%s %s", y, m, d, strings.Split(file, ".")[0]),
		time.Local,
	)
}

func fileToTime(date time.Time, file string) (time.Time, error) {
	t, err := time.ParseInLocation(util.FileTimeFormat, strings.Split(file, ".")[0], time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return util.DateAndTime(date, t), nil
}

// OpenRecord returns the open record if any. Returns a nil reference if no open record is found.
func (t *Track) OpenRecord() (*Record, error) {
	latest, err := t.LatestRecord()
	if err != nil {
		if err == fs.ErrNoFiles {
			return nil, nil
		}
		return nil, err
	}
	if latest.HasEnded() {
		return nil, nil
	}
	return latest, nil
}

// StartRecord starts and saves a record
func (t *Track) StartRecord(project, note string, tags []string, start time.Time) (Record, error) {
	record := Record{
		Project: project,
		Note:    note,
		Tags:    tags,
		Start:   start,
		End:     time.Time{},
	}

	return record, t.SaveRecord(&record, false)
}

// StopRecord stops and saves the current record
func (t *Track) StopRecord(end time.Time) (*Record, error) {
	record, err := t.OpenRecord()
	if err != nil {
		return record, err
	}
	if record == nil {
		return record, fmt.Errorf("no running record")
	}

	record.End = end

	err = t.SaveRecord(record, true)
	if err != nil {
		return record, err
	}
	return record, nil
}

// ExtractTagsSlice extracts elements with the tag prefix
func ExtractTagsSlice(tokens []string) []string {
	var result []string
	mapped := make(map[string]bool)
	for _, token := range tokens {
		subTokens := strings.Split(token, " ")
		for _, subToken := range subTokens {
			if strings.HasPrefix(subToken, TagPrefix) {
				if _, ok := mapped[subToken]; !ok {
					mapped[subToken] = true
					result = append(result, strings.TrimPrefix(subToken, TagPrefix))
				}
			}
		}
	}
	return result
}

// ExtractTags extracts elements with the tag prefix
func ExtractTags(text string) []string {
	var result []string
	mapped := make(map[string]bool)
	subTokens := strings.Split(text, " ")
	for _, subToken := range subTokens {
		if strings.HasPrefix(subToken, TagPrefix) {
			if _, ok := mapped[subToken]; !ok {
				mapped[subToken] = true
				result = append(result, strings.TrimPrefix(subToken, TagPrefix))
			}
		}
	}
	return result
}
