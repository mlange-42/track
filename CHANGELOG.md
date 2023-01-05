# Changelog

## [[unpublished]](https://github.com/mlange-42/track/compare/v0.3.2...main)

### Features

* New command `report tag` to list tags with time statistics (#93)
* Add flag `--dry` to `delete` command for dry-run (#96)

### Bugfixes

* Fix #94 broken child project time aggregation (#95)

## [[v0.3.2]](https://github.com/mlange-42/track/compare/v0.3.1...v0.3.2)

### Features

* Export records to JSON and YAML in addition to CSV (#88)
* Add total, work and pause durations to CSV export (#88)
* New command `report treemap` to generate an SVG treemap of time per project (#89)
* Flag `--rename` for command `edit project` for renaming projects (#90)
* Additional alias `?` for command `status` (#91)

## [[v0.3.1]](https://github.com/mlange-42/track/compare/v0.3.0...v0.3.1)

:warning: In case of an error after upgrading to `v0.3.1`, delete the config file `%USER%/.track/config.yml`.

### Features

* The data path (normally `%USER%/.track`) can be set via env var `TRACK_PATH` (#81)
* Configurable fill character for records in day report (#82)
* Adds flag `--7days` to week reports to show 7 days instead of a calendar week (#85)
* Adds flag `--copy` to the switch command to copy note and tags like in start command (#86)

### Bugfixes

* Fix crash when requesting the open record and there are no records at all (#80)
* Fix for different `yes` answers in confirm prompt (#84)

## [[v0.3.0]](https://github.com/mlange-42/track/compare/v0.2.1...v0.3.0)

### Features

* Switched to storing records in plain-text files with simple custom format (#52)
* Pause times can now be stored inside records instead of stop/start for a pause (#52)
* New command `pause` to pause a record, or insert a finished pause (#52)
* ~~The tag prefix is now `#` instead of `+` (record file comments are `//`) (#53)~~ (#62)
* Projects have a foreground color, in addition to the background color (#54)
* The `resume` command has flags `--at` and `--ago` to correct times (#55)
* The `start` command has a flag `--copy` to copy note and tags from the previous record in the same project (#58)
* The command `track edit record` can be used without arguments, or with only time as argument (#59)
* Command `track list tags` to list all tags, with number of occurrences (#60)
* Command `track edit day` to edit all records of a day in a single file (#61)
* Command `track report day` now prints a schedule-like report similar to `week` (#63)
* All `edit` subcommands have a flag `--dry` for dry-run (#68)
* All `edit` subcommands re-open the edited file if parsing/checks fail; improved error messages (#69)

### Bug fixes

* Apply `--at`/`--ago` also to start time in `switch`, not only to end time of previous record (#74)

### Other

* Removed command `break` (#52)
* Print only the first line of notes in the records list (#57)
* Record that starts the day before is included in `list records [DATE]` command (#62)
* Previous command `report day` (showing bar charts) renamed to `report chart` (#63)
* Print the latest record in `status`, in the default plain-text format (#75)
* Limit changes of end time when editing a single record (#76)
* Allow date words like "yesterday" in record selection for editing (#76)
* Extensive documentation under [mlange-42.github.io/track/](https://mlange-42.github.io/track/) (#71)

## [[v0.2.1]](https://github.com/mlange-42/track/compare/v0.2.0...v0.2.1)

### Bug fixes

* Limit time aggregation in reports to exact requested time range (#45)

## [[v0.2.0]](https://github.com/mlange-42/track/compare/v0.1.0...v0.2.0)

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
* Flags `--at` and `--ago` for all commands like `start`, `stop`, `switch` and `pause` (#72, #73)

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
