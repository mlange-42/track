# File format

This chapter describes the file format that is used for storing *Track* records.
The format is also used for editing records using the `edit record` and `edit day` commands.

[[_TOC_]]

## Overview

Records are stored in a simple, human-readable text format.

Records are stored in a directory structure representing the date,
with file names representing the starting time.
In the following example, the location of a record starting at `2023-01-10 8:15` is shown:

```text
records/
└─2023/
  └─01/
    └─10/
      └─08-15.trk
```

These files should not be edited directly.
*Track* provides an `edit` command that opens the entries to be edited in a temporary file,
and performs checks before replacing the original data.
See chapter [Manipulating data](./manipulating.md) for details.

The content of the file could look like this:

```text
# Record 2023-01-10 8:15
8:15 - 17:00
    - 10:15 - 15m / Breakfast
    - 13:00 - 30m / Lunch
    
    ProjectA

Work on +GUI +design
```

* The first line, starting with `#`, is a comment; date and time in it are just informative
* The next line represents the time span of the record.
* Subsequent lines that start with `-` (dash, plus optional indentation) are pauses
* The first non-empty line (rather, non-only-whitespace) after pauses is the project name
* Everything after the next non-empty line is the record's note, including tags

## Comments and empty lines

Lines that start with `#` (exactly, no indent/whitespace allowed) are comments.
Comments are ignored.

Lines that are completely empty, or that contain only whitespace characters (SPACE, TAB) are considered empty.
Lines considered empty are ignored, except within the note. Lines considered empty before and after any non-empty note lines are ignored.

## Structure

* The first line that is not ignored (i.e. not comment or "empty") represents the time span of the record.
* Subsequent lines that start with `-` (dash, plus optional indentation) are pauses
* The first line after pauses that is not ignored is the project name (excluding optional indentation)
* Everything after any subsequent ignored lines it the record's note; notes can comprise multiple lines

## Time ranges

There are three ways to define time ranges:

* A starting time and an end time, separated by `-` (dash, surrounded by optional spaces):  
  `08:15 - 17:00`
* A starting time and a duration, separated by `-` (as above):  
  `08:15 - 8h45m`
* An open, still running time span is defined by `?` as the second element:  
  `08:15 - ?`

### Time and duration format

Times can be specified in the format `hh:mm` or `h:mm`, like `08:15` or `8:15`.
All times are in 24h format. 12h `am`/`pm` format is not supported.

Durations are in the usual Go format: `10h15m23s`. Zero-valued entries can be left out. So e.g. `15m` is also valid.

### Day shifts

In some cases, particularly when tracking time over midnight, the starting time of a record may be on another day than the end time. This is denoted by the shift markers `<` and `>`.

The most common case is a record that goes over midnight and ends the day after its start.
The time range of such a record would look like this:

```
22:00 - 00:30>
```

Another case is full-day editing using `edit day` (see [Temporary multi-record files](#temporary-multi-record-files)).
Here, a record that starts the day before but ends on the day to edit would start like this:

```
<22:00 - 00:30
```

## Pauses

A record can contain an arbitrary number of pauses.

All lines immediately after the record's time range that start with `-` (dash, plus optional indentation) are considered pause entries.

A pause has the form:

```
- start - duration / Note
- start - duration
- start - end / Note
- start - end
```

The following are valid pause entries:

```
- 10:00 - 20m / Breakfast
- 10:00 - 20m
- 10:00 - 10:20 / Breakfast
- 10:00 - 10:20
```

*Track* uses the duration version for saving records. When editing records, both forms are valid.

For parsing the pause's time range, the rules of [Time ranges](#time-ranges) apply.

Pauses are optional.

Pauses must be listed in chronological order, and must not overlap.
Pauses must not exceed the record's time span.

## Project

The first line after any (optional) pause entries that is not ignored (i.e. not comment or "empty")
is considered the project name. Any whitespace characters at the start and the end of the line are removed. I.e. indentation can be used.

The project name is obligatory.

## Note

The note is optional.

All lines after the project name are considered the note.
Any "empty" lines at the start and the end of the note are removed.
Empty lines between non-empty lines of a note are preserved.

A note can contain tags.

## Tags

Tags are derived from the note, and are optional.

Tags are identified by the prefix `+` and must be surrounded/are delimited by spaces.
Tag can have an optional value, which is separated from the tag's name by `=` (without any whitespace characters).

Here is an example of a note that contains a tag `tag` without a value, and a tag `key` with a value:

```
A note featuring a +tag and a +key=value pair for a tag with a value
```

## Temporary multi-record files

When using the `edit day` command, *Track* assembles the respective records in a single temporary file for the user to edit.
In this file, lines that start with `----` (4 dashes) delimit individual records:

```
8:15 - 13:00
    ProjectA

Work on +GUI +design

-------------

14:00 - 13:30
    ProjectA

Working group +meeting

-------------

13:30 - 15:00
    ProjectB

Draft +paper
```
