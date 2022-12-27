# track

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
* Everything stored as YAML for human readability
* Different types of text-based and graphical reports (in progress)

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
+------------------+-------+-------+-------+
|          project |  curr | today | break |
|        MyProject | 00:30 | 01:45 | 00:13 |
+------------------+-------+-------+-------+
```

#### `list`

The `list` commands lists projects or records.

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

Subcommand `timeline` shows statistics of time tracked per day, week or month:

```shell
track report timeline days
```

Subcommand `day` shows a timeline over the current or the given day:

```shell
track report day
```

### All subcommands

```text
track
├─break DURATION
├─create
│ ├─project PROJECT
│ └─workspace WORKSPACE
├─delete
│ └─record DATE TIME
├─edit
│ ├─config
│ ├─project PROJECT
│ └─record DATE TIME
├─export
│ └─records
├─list
│ ├─projects
│ ├─records [DATE]
│ └─workspaces
├─report
│ ├─day [DATE]
│ ├─projects
│ └─timeline (days|weeks|months)
├─resume [NOTE...]
├─start PROJECT [NOTE...]
├─status [PROJECT]
├─stop
├─switch PROJECT [NOTE...]
└─workspace WORKSPACE
```

## References

* Heavily inspired by [`github.com/dominikbraun/timetrace`](https://github.com/dominikbraun/timetrace)
