package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mlange-42/track/fs"
	"github.com/mlange-42/track/util"
	"gopkg.in/yaml.v3"
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
}

type yamlRecord struct {
	Project string
	Start   util.Time
	End     util.Time
	Note    string
	Tags    []string
}

// MarshalYAML converts a Record to YAML bytes
func (r *Record) MarshalYAML() (interface{}, error) {
	return &yamlRecord{
		Project: r.Project,
		Note:    r.Note,
		Tags:    r.Tags,
		Start:   util.Time(r.Start),
		End:     util.Time(r.End),
	}, nil
}

// UnmarshalYAML converts YAML bytes to a Record
func (r *Record) UnmarshalYAML(value *yaml.Node) error {
	rec := yamlRecord{}
	if err := value.Decode(&rec); err != nil {
		return err
	}
	r.Project = rec.Project
	r.Note = rec.Note
	r.Tags = rec.Tags
	r.Start = time.Time(rec.Start)
	r.End = time.Time(rec.End)

	return nil
}

// HasEnded reports whether the record has an end time
func (r Record) HasEnded() bool {
	return !r.End.IsZero()
}

// Duration reports the duration of a record
func (r Record) Duration() time.Duration {
	t := r.End
	if t.IsZero() {
		t = time.Now()
	}
	return t.Sub(r.Start)
}

// RecordPath returns the full path for a record
func (t *Track) RecordPath(tm time.Time) string {
	return filepath.Join(
		t.RecordDir(tm),
		fmt.Sprintf("%s.yml", tm.Format(util.FileTimeFormat)),
	)
}

// RecordDir returns the directory path for a record
func (t *Track) RecordDir(tm time.Time) string {
	return filepath.Join(
		fs.RecordsDir(),
		tm.Format(util.FileDateFormat),
	)
}

// SaveRecord saves a record to disk
func (t *Track) SaveRecord(record Record, force bool) error {
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

	bytes, err := yaml.Marshal(&record)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "# Record %s\n\n", record.Start.Format(util.DateTimeFormat))
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// DeleteRecord deletes a record
func (t *Track) DeleteRecord(record Record) error {
	path := t.RecordPath(record.Start)
	if !fs.FileExists(path) {
		return fmt.Errorf("record does not exist")
	}
	return os.Remove(path)
}

// LoadRecordByTime loads a record
func (t *Track) LoadRecordByTime(tm time.Time) (Record, error) {
	path := t.RecordPath(tm)
	return t.LoadRecord(path)
}

// LoadRecord loads a record
func (t *Track) LoadRecord(path string) (Record, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Record{}, err
	}

	var record Record

	if err := yaml.Unmarshal(file, &record); err != nil {
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
	path := fs.RecordsDir()

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var records []Record

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		recs, err := t.LoadDateRecordsFiltered(dir.Name(), filters)
		if err != nil {
			return nil, err
		}
		records = append(records, recs...)
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
		path := fs.RecordsDir()

		dirs, err := ioutil.ReadDir(path)
		if err != nil {
			results <- FilterResult{Record{}, err}
			return
		}

		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}
			recs, err := t.LoadDateRecordsFiltered(dir.Name(), filters)
			if err != nil {
				results <- FilterResult{Record{}, err}
				return
			}
			for _, rec := range recs {
				results <- FilterResult{rec, nil}
			}
		}
		close(results)
	}, results
}

// AllRecordsFiltered is an async version of LoadAllRecordsFiltered
func (t *Track) allRecordsFiltered(
	filters FilterFunctions,
	results chan struct {
		Record
		error
	}) {
	path := fs.RecordsDir()

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		results <- struct {
			Record
			error
		}{Record{}, err}
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		recs, err := t.LoadDateRecordsFiltered(dir.Name(), filters)
		if err != nil {
			results <- struct {
				Record
				error
			}{Record{}, err}
			return
		}
		for _, rec := range recs {
			results <- struct {
				Record
				error
			}{rec, nil}
		}
	}
	close(results)
}

// LoadDateRecords loads all records for the given date string/directory
func (t *Track) LoadDateRecords(dir string) ([]Record, error) {
	return t.LoadDateRecordsFiltered(dir, []func(*Record) bool{})
}

// LoadDateRecordsFiltered loads all records for the given date string/directory
func (t *Track) LoadDateRecordsFiltered(dir string, filters FilterFunctions) ([]Record, error) {
	path := fs.RecordsDir()
	subPath := filepath.Join(path, dir)

	info, err := os.Stat(subPath)
	if err != nil {
		return nil, ErrNoRecords
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", dir)
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

		record, err := t.LoadRecord(filepath.Join(subPath, file.Name()))
		if err != nil {
			return nil, err
		}
		if Filter(&record, filters) {
			records = append(records, record)
		}
	}

	return records, nil
}

// LatestRecord loads the latest record
func (t *Track) LatestRecord() (Record, error) {
	records := fs.RecordsDir()
	records, err := fs.FindLatests(records, true)
	if err != nil {
		return Record{}, err
	}
	record, err := fs.FindLatests(records, false)
	if err != nil {
		return Record{}, err
	}

	rec, err := t.LoadRecord(record)
	if err != nil {
		return Record{}, err
	}

	return rec, nil
}

// OpenRecord returns the open record if any
func (t *Track) OpenRecord() (rec Record, ok bool) {
	latest, err := t.LatestRecord()
	if err != nil {
		if err == fs.ErrNoFiles {
			return Record{}, false
		}
		return Record{}, false
	}
	if latest.HasEnded() {
		return Record{}, false
	}
	return latest, true
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

	return record, t.SaveRecord(record, false)
}

// StopRecord stops and saves the current record
func (t *Track) StopRecord(end time.Time) (Record, error) {
	record, ok := t.OpenRecord()
	if !ok {
		return record, fmt.Errorf("no running record")
	}

	record.End = end

	err := t.SaveRecord(record, true)
	if err != nil {
		return record, err
	}
	return record, nil
}

// ExtractTags extracts elements with the tag prefix
func (t *Track) ExtractTags(tokens []string) []string {
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
