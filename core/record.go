package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mlange-42/track/fs"
)

// Record holds and manipulates data for a record
type Record struct {
	Project string
	Note    string
	Start   time.Time
	End     time.Time
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
func (t *Track) RecordPath(record Record) string {
	return filepath.Join(
		t.RecordDir(record),
		fmt.Sprintf("%02d-%02d-%02d.json", record.Start.Hour(), record.Start.Minute(), record.Start.Second()),
	)
}

// RecordDir returns the directory path for a record
func (t *Track) RecordDir(record Record) string {
	return filepath.Join(
		fs.RecordsDir(),
		fmt.Sprintf("%04d-%02d-%02d", record.Start.Year(), record.Start.Month(), record.Start.Day()),
	)
}

// SaveRecord saves a record to disk
func (t *Track) SaveRecord(record Record, force bool) error {
	path := t.RecordPath(record)
	if !force && fs.FileExists(path) {
		return fmt.Errorf("Record already exists")
	}
	dir := t.RecordDir(record)
	err := fs.CreateDir(dir)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&record, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// LoadRecord loads a record
func (t *Track) LoadRecord(path string) (Record, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Record{}, err
	}

	var record Record

	if err := json.Unmarshal(file, &record); err != nil {
		return Record{}, err
	}

	return record, nil
}

// LoadAllRecords loads all records
func (t *Track) LoadAllRecords() ([]Record, error) {
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
		subPath := filepath.Join(path, dir.Name())

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
