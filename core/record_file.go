package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mlange-42/track/util"
)

// FilterResult contains a Record or an error from async filtering
type FilterResult struct {
	Record Record
	Err    error
}

type workerResult struct {
	Index  int
	Record Record
	Err    error
}

type listFilterResult struct {
	Time time.Time
	Err  error
}

// NewRecord creates a new record
func (t *Track) NewRecord(project *Project, note string, tags map[string]string, start time.Time, end time.Time) (Record, error) {
	record := Record{
		Project: project.Name,
		Note:    note,
		Tags:    tags,
		Start:   start,
		End:     end,
		Pause:   []Pause{},
	}

	if err := record.Check(project); err != nil {
		return record, err
	}

	return record, t.SaveRecord(&record, false)
}

// StartRecord starts a new record for the given project at the given time.
func (t *Track) StartRecord(project *Project, note string, tags map[string]string, start time.Time) (Record, error) {
	return t.NewRecord(project, note, tags, start, util.NoTime)
}

// StopRecord stops the currently running record at the given time, and saves it to disk.
func (t *Track) StopRecord(end time.Time) (*Record, error) {
	record, err := t.OpenRecord()
	if err != nil {
		return record, err
	}
	if record == nil {
		return record, fmt.Errorf("no running record")
	}

	record.End = end
	for len(record.Pause) > 0 {
		idx := len(record.Pause) - 1
		if record.Pause[idx].End.IsZero() || record.Pause[idx].End.After(end) {
			record.End = record.Pause[idx].Start
			record.Pause = record.Pause[:idx]
		} else {
			break
		}
	}

	err = t.SaveRecord(record, true)
	if err != nil {
		return record, err
	}
	return record, nil
}

// LoadRecord loads a record by the given start time
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

// OpenRecord returns the open/running record if any.
// Returns a nil reference if no open record is found.
func (t *Track) OpenRecord() (*Record, error) {
	latest, err := t.LatestRecord()
	if err != nil {
		if err == util.ErrNoFiles {
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

// LatestRecord loads the latest record, open/running or not.
// Returns a nil reference if no record is found.
func (t *Track) LatestRecord() (*Record, error) {
	records := t.RecordsDir()
	yearPath, year, err := util.FindLatests(records, true)
	if err != nil {
		if errors.Is(err, util.ErrNoFiles) {
			return nil, nil
		}
		return nil, err
	}
	monthPath, month, err := util.FindLatests(yearPath, true)
	if err != nil {
		if errors.Is(err, util.ErrNoFiles) {
			return nil, nil
		}
		return nil, err
	}
	dayPath, day, err := util.FindLatests(monthPath, true)
	if err != nil {
		if errors.Is(err, util.ErrNoFiles) {
			return nil, nil
		}
		return nil, err
	}
	_, record, err := util.FindLatests(dayPath, false)
	if err != nil {
		if errors.Is(err, util.ErrNoFiles) {
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

// FindLatestRecord loads the latest record that matches the given FilterFunction.
// Returns a nil reference if no record is found.
func (t *Track) FindLatestRecord(cond FilterFunction) (*Record, error) {
	fn, results, stop := t.AllRecordsFiltered(
		FilterFunctions{[]FilterFunction{cond}, util.NoTime, util.NoTime},
		true, // reversed order to find latest record of project
	)
	go fn()

	res := <-results
	if res.Err != nil {
		return nil, res.Err
	}
	close(stop)
	return &res.Record, nil
}

// LoadAllRecords loads all records.
func (t *Track) LoadAllRecords() ([]Record, error) {
	return t.LoadAllRecordsFiltered(NewFilter([]func(*Record) bool{}, util.NoTime, util.NoTime))
}

// LoadAllRecordsFiltered loads all records, filtered by FilterFunctions.
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

// AllRecords is an async version of LoadAllRecords.
//
// Returns a function to be run as goroutine,
// a channel for results, and a channel that can be closed
// to signal end of the search.
func (t *Track) AllRecords() (func(), chan FilterResult, chan struct{}) {
	return t.AllRecordsFiltered(NewFilter([]func(*Record) bool{}, util.NoTime, util.NoTime), false)
}

// AllRecordsFiltered is an async version of LoadAllRecordsFiltered.
//
// Returns a function to be run as goroutine,
// a channel for results, and a channel that can be closed
// to signal end of the search.
func (t *Track) AllRecordsFiltered(filters FilterFunctions, reversed bool) (func(), chan FilterResult, chan struct{}) {
	numWorkers := 32
	results := make(chan FilterResult, 64)

	fn, listResults, stop := t.listAllRecordsFiltered(filters, reversed)

	return func() {
		defer close(results)

		go fn()

		worker := func(index int, tasks chan time.Time, ch chan workerResult) {
			for tm := range tasks {
				record, err := t.LoadRecord(tm)
				ch <- workerResult{index, record, err}
			}
		}

		process := func(index int, times []time.Time, taskChannels []chan time.Time, resChannels []chan workerResult) {
			for i := 0; i < index; i++ {
				taskChannels[i] <- times[i]
			}
			for i := 0; i < index; i++ {
				select {
				case <-stop:
					return
				default:
				}

				res := <-resChannels[i]

				fr := FilterResult{res.Record, res.Err}
				if res.Err != nil {
					results <- fr
					return
				}
				if Filter(&res.Record, filters) {
					results <- fr
				}
			}
		}

		tempTimes := make([]time.Time, numWorkers)

		taskChannels := make([]chan time.Time, numWorkers)
		resChannels := make([]chan workerResult, numWorkers)

		index := 0

		for i := 0; i < numWorkers; i++ {
			taskChannels[i] = make(chan time.Time, 4)
			resChannels[i] = make(chan workerResult, 4)
			defer close(taskChannels[i])
			go worker(i, taskChannels[i], resChannels[i])
		}

		for rec := range listResults {
			if rec.Err != nil {
				results <- FilterResult{Record{}, rec.Err}
				return
			}
			tempTimes[index] = rec.Time

			index++
			if index >= numWorkers {
				process(index, tempTimes, taskChannels, resChannels)
				index = 0
				select {
				case <-stop:
					return
				default:
				}
			}
		}
		if index > 0 {
			process(index, tempTimes, taskChannels, resChannels)
		}
	}, results, stop
}

func (t *Track) listAllRecordsFiltered(filters FilterFunctions, reversed bool) (func(), chan listFilterResult, chan struct{}) {
	results := make(chan listFilterResult, 64)
	stop := make(chan struct{})

	return func() {
		defer close(results)

		path := t.RecordsDir()

		yearDirs, err := os.ReadDir(path)
		if err != nil {
			results <- listFilterResult{util.NoTime, err}
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
				results <- listFilterResult{util.NoTime, err}
				return
			}
			if !filters.Start.IsZero() && year < filters.Start.Year() {
				continue
			}
			if !filters.End.IsZero() && year > filters.End.Year() {
				continue
			}

			monthDirs, err := os.ReadDir(filepath.Join(path, yearDir.Name()))

			if reversed {
				util.Reverse(monthDirs)
			}
			if err != nil {
				results <- listFilterResult{util.NoTime, err}
				return
			}

			for _, monthDir := range monthDirs {
				if !monthDir.IsDir() {
					continue
				}
				month, err := strconv.Atoi(monthDir.Name())
				if err != nil {
					results <- listFilterResult{util.NoTime, err}
					return
				}

				dayDirs, err := os.ReadDir(filepath.Join(path, yearDir.Name(), monthDir.Name()))
				if err != nil {
					results <- listFilterResult{util.NoTime, err}
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
						results <- listFilterResult{util.NoTime, err}
						return
					}

					date := util.Date(year, time.Month(month), day)
					if !filters.Start.IsZero() && date.Before(util.ToDate(filters.Start)) {
						continue
					}
					if !filters.End.IsZero() && date.After(filters.End) {
						continue
					}

					recs, err := t.listDateRecords(date)
					if err != nil {
						results <- listFilterResult{util.NoTime, err}
						return
					}

					if reversed {
						util.Reverse(recs)
					}
					for _, rec := range recs {
						select {
						case <-stop:
							return
						case results <- listFilterResult{rec, nil}:
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

// LoadDateRecordsExact loads all records for the given date,
// including those starting the day before but ending at the given date.
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

// LoadDateRecordsFiltered loads all records for the given date,
// filtered by FilterFunctions.
func (t *Track) LoadDateRecordsFiltered(date time.Time, filters FilterFunctions) ([]Record, error) {
	recs, err := t.listDateRecords(date)
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0, len(recs))

	for _, tm := range recs {
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

func (t *Track) listDateRecords(date time.Time) ([]time.Time, error) {
	subPath := t.RecordDir(date)

	info, err := os.Stat(subPath)
	if err != nil {
		return nil, ErrNoRecords
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", info.Name())
	}

	var records []time.Time

	files, err := os.ReadDir(subPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		tm, err := fileToTime(date, file.Name())
		if err != nil {
			return nil, err
		}
		records = append(records, tm)
	}

	return records, nil
}

// SaveRecord saves the given record to disk.
// Argument `force` allows to overwrite an existing file.
func (t *Track) SaveRecord(record *Record, force bool) error {
	path := t.RecordPath(record.Start)
	if !force && util.FileExists(path) {
		return fmt.Errorf("record already exists")
	}
	dir := t.RecordDir(record.Start)
	err := util.CreateDir(dir)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes := SerializeRecord(record, util.NoTime)

	_, err = fmt.Fprintf(file, "%s Record %s\n", CommentPrefix, record.Start.Format(util.DateTimeFormat))
	if err != nil {
		return err
	}

	_, err = file.WriteString(bytes)

	return err
}

// DeleteRecord deletes a record
func (t *Track) DeleteRecord(record *Record) error {
	path := t.RecordPath(record.Start)
	if !util.FileExists(path) {
		return fmt.Errorf("record does not exist")
	}
	err := os.Remove(path)
	if err != nil {
		return err
	}
	dayDir := filepath.Dir(path)
	empty, err := util.DirIsEmpty(dayDir)
	if err != nil {
		return err
	}
	if empty {
		os.Remove(dayDir)
		monthDir := filepath.Dir(dayDir)
		empty, err := util.DirIsEmpty(monthDir)
		if err != nil {
			return err
		}
		if empty {
			os.Remove(monthDir)
			yearDir := filepath.Dir(monthDir)
			empty, err := util.DirIsEmpty(yearDir)
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
