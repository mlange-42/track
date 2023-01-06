# Reports

*Track*'s `report` command provides different textual and graphical reports for (potentially filtered) tracking data.

[[_TOC_]]

## Filters

All `report` sub-commands support filtering via flags, for:
* Projects with `--projects`
* Tags with `--tags`

Lists for these flags should be comma-separated, like `--projects ProjectA,ProjectB`.

Further, most sub-commands support restricting the time range using the flags `--start` and `--end`. Both flags accept a date, like `2023-01-01` or `yesterday`. The end date is inclusive.

## Projects report

Command `report projects` prints a tree-like list of projects, with total time (incl. child projects) and time spent per project:

```
track report projects
```

Prints something like this:

```text
<default>
└─Private         P  08:25 (00:00)
  └─Coding        C  08:25 (00:00)
    └─MyApp       M  08:25 (08:25)
```

Here is an example using filters:

```
track report projects --start 2023-01-01 --end 2023-01-07 --projects MyApp --tags GUI,design
```

## Tags report

Command `report tags` prints a list of tags, with work time and pause time per tag.
Usage is the same as for `report projects`.

## Week report

Command `report week` prints a time-table of the current or given week:

```
track report week
track report week 2023-01-01
```

## Day report

Command `report day` prints a time-table of the current or given day, similar to the [Week report](#week-report). In addition, record bars are labelled with the record's note

```
track report day
track report day yesterday
track report day 2023-01-01
```

## Chart report

Command `report chart` shows the time spent per project, as a bar chart time series over the current or given day:

```shell
track report chart
track report chart yesterday
track report chart 2023-01-01
```

Prints something like this:

```text
                    |2023-01-01 : 20m0s/cell
<default>           |00:00    |03:00    |06:00    |09:00    |12:00    |15:00    |18:00    |21:00    |
└─Private         P |.........|.........|.........|.........|.........|.........|.........|.........|
  └─Coding        C |.........|.........|.........|.........|.........|.........|.........|.........|
    └─MyApp       M |███▂.....|.........|.........|.........|.....▂█▄.|.██.▂████|█▄.▅.....|.........|
```

## Treemap report

Command `report treemap` generates an SVG treemap visualization of time spent per project.
Here, we pipe the SVG to a file:

```shell
track report treemap > test.svg
```

You can also open the file with the default program for SVG (ideally a web browser) immediately:

```shell
track report treemap > test.svg && test.svg
```

## Timeline reports

Command `report timeline` shows total time spent per day, week or month as a bar chart time series:

```
track report timeline days
track report timeline weeks
track report timeline months
```

Prints something like this:

```text
Th 2022-12-29  06:39  |||||||||||||.
Fr 2022-12-30  09:20  ||||||||||||||||||:
Sa 2022-12-31  03:51  |||||||:
Su 2023-01-01  11:07  ||||||||||||||||||||||.
Mo 2023-01-02  09:44  |||||||||||||||||||.
Tu 2023-01-03  08:30  |||||||||||||||||
We 2023-01-04  09:51  |||||||||||||||||||:
Th 2023-01-05  09:01  ||||||||||||||||||
Fr 2023-01-06  03:30  |||||||
```

Timeline reports can be exported in CSV format using the flag `--csv`.
With flag `--table`, a separate column for each project is included in the report.
