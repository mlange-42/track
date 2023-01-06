package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mlange-42/track/fs"
	"github.com/mlange-42/track/util"
)

// LoadRecord loads a record
func (t *Track) LoadRecord(tm time.Time) (Record, error) {
	path := t.RecordPath(tm)
	file, err := os.ReadFile(path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return Record{}, ErrRecordNotFound
		}
		return Record{}, err
	}

	record, err := DeserializeRecord(string(file), tm)
	if err != nil {
		return Record{}, err
	}

	return record, nil
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
	if latest == nil {
		return nil, nil
	}
	if latest.HasEnded() {
		return nil, nil
	}
	return latest, nil
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

// FindLatestRecord loads the latest record for the given condition. Returns a nil reference if no record is found.
func (t *Track) FindLatestRecord(cond FilterFunction) (*Record, error) {
	fn, results, stop := t.AllRecordsFiltered(
		FilterFunctions{[]FilterFunction{cond}, util.NoTime, util.NoTime},
		true, // reversed order to find latest record of project
	)
	go fn()

	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		close(stop)
		return &res.Record, nil
	}
	return nil, nil
}

// LoadAllRecords loads all records
func (t *Track) LoadAllRecords() ([]Record, error) {
	return t.LoadAllRecordsFiltered(NewFilter([]func(*Record) bool{}, util.NoTime, util.NoTime))
}

// LoadAllRecordsFiltered loads all records
func (t *Track) LoadAllRecordsFiltered(filters FilterFunctions) ([]Record, error) {
	fn, results, _ := t.AllRecordsFiltered(filters, false)
	go fn()

	var records []Record
	for res := range results {
		if res.Err != nil {
			return records, res.Err
		}
		records = append(records, res.Record)
	}

	return records, nil
}

// AllRecordsFiltered is an async version of LoadAllRecordsFiltered
func (t *Track) AllRecordsFiltered(filters FilterFunctions, reversed bool) (func(), chan FilterResult, chan struct{}) {
	results := make(chan FilterResult, 32)
	stop := make(chan struct{})

	return func() {
		defer close(results)

		path := t.RecordsDir()

		yearDirs, err := ioutil.ReadDir(path)
		if err != nil {
			results <- FilterResult{Record{}, err}
			return
		}
		if reversed {
			util.Reverse(yearDirs)
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
			if !filters.Start.IsZero() && year < filters.Start.Year() {
				continue
			}
			if !filters.End.IsZero() && year > filters.End.Year() {
				continue
			}

			monthDirs, err := ioutil.ReadDir(filepath.Join(path, yearDir.Name()))

			if reversed {
				util.Reverse(monthDirs)
			}
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

				if reversed {
					util.Reverse(dayDirs)
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

					date := util.Date(year, time.Month(month), day)
					if !filters.Start.IsZero() && date.Before(util.ToDate(filters.Start)) {
						continue
					}
					if !filters.End.IsZero() && date.After(filters.End) {
						continue
					}

					recs, err := t.LoadDateRecordsFiltered(date, filters)
					if err != nil {
						results <- FilterResult{Record{}, err}
						return
					}

					if reversed {
						util.Reverse(recs)
					}
					for _, rec := range recs {
						select {
						case <-stop:
							return
						case results <- FilterResult{rec, nil}:
						}
					}
				}
			}
		}
	}, results, stop
}

// LoadDateRecords loads all records for the given date
func (t *Track) LoadDateRecords(date time.Time) ([]Record, error) {
	return t.LoadDateRecordsFiltered(date, FilterFunctions{})
}

// LoadDateRecordsExact loads all records for the given date, including those starting the das before
func (t *Track) LoadDateRecordsExact(date time.Time) ([]Record, error) {
	date = util.ToDate(date)
	dateBefore := date.Add(-24 * time.Hour)
	dateAfter := date.Add(24 * time.Hour)

	filters := FilterFunctions{
		[]FilterFunction{FilterByTime(date, dateAfter)},
		util.NoTime,
		util.NoTime,
	}

	records, err := t.LoadDateRecordsFiltered(dateBefore, filters)
	if err != nil && !errors.Is(err, ErrNoRecords) {
		return nil, err
	}
	records2, err := t.LoadDateRecordsFiltered(date, filters)
	if err != nil && !errors.Is(err, ErrNoRecords) {
		return nil, err
	}
	records = append(records, records2...)

	if len(records) == 0 {
		return nil, ErrNoRecords
	}
	return records, nil
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
