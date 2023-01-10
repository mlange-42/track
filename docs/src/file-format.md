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
records
└─2023
  └─01
    └─10
      └─08-15.trk
```

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

[TODO]

## Pauses

Pauses are optional.

Pauses must be listed in chronological order, and must not overlap.
Pauses must not exceed the record's time span.

[TODO]

## Note

The note is optional.

[TODO]

## Tags

Tags are derived from the note, and are optional.

[TODO]

## Temporary multi-record files

[TODO]
