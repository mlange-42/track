# Changelog

## [unpublished]

### Features

* Day report with horizontal timelines over projects as rows: `track report day [DATE]` (#25)
* Command `track break DURATION` to insert breaks (#28)
* Projects are structured in workspaces (#30)
* Records can be started and stopped with explicit time or offset from now (#31)
* Shows aggregated and self time of projects in report (#33)
* Projects have a color (#35)
* Week report with vertical day columns and projects denoted by color and initial letter: `track report week [DATE]` (#35)

### Bug fixes

* Store time zone in record files (#24)
* Use local time zone when parsing dates and times from string (#24)
* ~~Disable flawed time aggregation over parent projects (#24)~~ (#33)
* Correct time aggregation from child projects (#33)

### Other

* Simplify CLI help usage strings (#26)
* Date is optional in `track list records` (#27)
