# Track

[![Test status](https://github.com/mlange-42/track/actions/workflows/tests.yml/badge.svg)](https://github.com/mlange-42/track/actions/workflows/tests.yml)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/track)
[![MIT license](https://img.shields.io/github/license/mlange-42/track)](https://github.com/mlange-42/track/blob/main/LICENSE)

Track is a time tracking command line tool

## Installation

**Using Go:**

```shell
go install github.com/mlange-42/track@latest
```

**Without Go:**

Download binaries for your OS from the [Releases](https://github.com/mlange-42/track/releases/).

## Features

* Track your working time from the command line
* Natural language-like syntax
* Supports hierarchical project structure
* Records stored as plain-text files for human readability and editing
* Different types of text-based and graphical reports

## Usage

* [First steps](#first-steps)
* [Getting information](#getting-information)
* [All subcommands](#all-subcommands)

Get Help:

```shell
track -h
track <command> -h
```

### First steps

Any time tracking `track` is associated to a *Project*.
Before you can start tracking, create a project:

```shell
track create project MyProject
```

Now, start tracking time on the project:

```shell
track start MyProject
```

To stop tracking, use:

```shell
track stop
```

### Getting information

`track` provides several commands to display information.

#### `status`

The `status` status informs about the running project:

```shell
track status
```

Prints something like:

```text
+------------------+-------+-------+-------+-------+
|          project |  curr | total | break | today |
|        MyProject | 01:05 | 01:05 | 00:10 | 01:53 |
+------------------+-------+-------+-------+-------+
```

#### `list`

The `list` commands lists projects or records:

```shell
track list projects
```

```shell
track list records today
track list records yesterday
track list records 2022-12-31
```

#### `report`

The `report` command generates several kinds of reports.

Subcommand `projects` shows statistics of time tracked per project:

```shell
track report projects
```

Subcommand `day` shows a timeline over the current or the given day:

```shell
track report day
```

Subcommand `day` shows a calendar-like view of the current or given week:

```shell
track report week
```

Subcommand `timeline` shows statistics of time tracked per day, week or month:

```shell
track report timeline days
```

### All subcommands

```text
track
├─create
│ ├─project PROJECT
│ └─workspace WORKSPACE
├─delete
│ ├─project PROJECT
│ └─record DATE TIME
├─edit
│ ├─config
│ ├─day [DATE]
│ ├─project PROJECT
│ └─record [[DATE] TIME]
├─export
│ └─records
├─list
│ ├─colors
│ ├─projects
│ ├─records [DATE]
│ ├─tags
│ └─workspaces
├─pause [NOTE...]
├─report
│ ├─chart [DATE]
│ ├─day [DATE]
│ ├─projects
│ ├─timeline (days|weeks|months)
│ └─week [DATE]
├─resume [NOTE...]
├─start PROJECT [NOTE...]
├─status [PROJECT]
├─stop
├─switch PROJECT [NOTE...]
└─workspace WORKSPACE
```

## References

* Heavily inspired by [`timetrace`](https://github.com/dominikbraun/timetrace) and [`klog`](https://github.com/jotaen/klog)
