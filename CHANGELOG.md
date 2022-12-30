# Changelog

## [unpublished]

### Features

* Switched to storing records in plain-text files with simple custom format (#52)
* Pause times can now be stored inside records instead of stop/start for a pause (#52)
* New command `pause` to pause a record, or insert a finished pause (#52)
* The tag prefix is now `#` instead of `+` (record file comments are `//`) (#53)
* Projects have a foreground color, in addition to the background color (#54)
* The `resume` command has flags `--at` and `--ago` to correct times (#55)
* The `start` command has a flag `--copy` to copy note and tags from the previous record in the same project (#58)
* The command `track edit record` can be used without arguments, or with only time as argument (#59)
* Command `track list tags` to list all tags, with number of occurrences (#60)
* Command `track edit day` to edit all records of a day in a single file (#61)

### Other

* Removed command `break` (#52)
* Print only the first line of notes in the records list (#57)

## [v0.2.1]

### Bug fixes

* Limit time aggregation in reports to exact requested time range (#45)

## [v0.2.0]

### Features

* Day report with horizontal timelines over projects as rows: `track report day [DATE]` (#25)
* Command `track break DURATION` to insert breaks (#28)
* Projects are structured in workspaces (#30)
* Records can be started and stopped with explicit time or offset from now (#31)
* Shows aggregated and self time of projects in report (#33)
* Projects have a color (#35)
* List available colors with `track list colors` (#35)
* Week report with vertical day columns and projects denoted by color and initial letter: `track report week [DATE]` (#35)
* Auto-scale day and week report width based on terminal width (#35)
* Projects can be archived (and un-archived) to exclude them from reports etc. (#39)
* Project colors and symbols are shown in reports and lists (#40)
* Weekly sum of project time shown in week report (#41)
* Delete a project with `track delete project PROJECT` (#43)

### Bug fixes

* Store time zone in record files (#24)
* Use local time zone when parsing dates and times from string (#24)
* ~~Disable flawed time aggregation over parent projects (#24)~~ (#33)
* Correct time aggregation from child projects (#33)
* Check parents to prevent circular relations (#38)

### Other

* Simplify CLI help usage strings (#26)
* Date is optional in `track list records` (#27)
* Records are stored in hierarchical folder structure like `2022/12/31` instead of folders like `2022-12-31` (#44)
