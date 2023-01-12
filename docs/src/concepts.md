# Concepts

This chapter explains *Track*'s primary concepts.

To start using *Track* immediately, you can skip to chapter [Time tracking](./tracking.md) and come back here if something seems unclear.

[[_TOC_]]

## Records

Time tracking entries are organized in records.

Each record is associated to a [project](#projects).
A record is defined as a  contiguous time span spent on a project.
It is characterized by it's start and end time.

Each records can contain an arbitrary number of pauses. Each pause is characterized by a start time and a duration.

Further, each record can have a note and tags.

For details, see chapter [Time tracking](./tracking.md).

## Projects

Projects are the primary way of structuring time tracking in *Track*.
Projects can be organized in a hierarchical tree-like structure.

For details, see chapter [Projects](./projects.md).

## Command line interface

The command line interface heavily relies on nested subcommands.
This allows for a natural language-like syntax, like

```shell
track list records yesterday
```

For normal usage, flags (i.e. options prefixed with `--` or `-`) are rarely required.

Further, all subcommands can be abbreviated by their first letter (with a few exceptions that use a different character or two letters). E.g. this is equivalent to the above list command:

```shell
track l r yesterday
```

For the full subcommand tree, see appendix [Command tree](./command-tree.md).

## File format

*Track* uses a human-readable plain-text format to store records.
This allows for easy editing, simply using a text editor.

A time tracking record looks like this:

```text
8:15 - 17:00
    - 10:15 - 15m / Breakfast
    - 13:00 - 30m / Lunch
    
    ProjectA

Work on +GUI +design
```

*Track* provides an `edit` command that opens the entries to be edited in a temporary file,
and performs checks before replacing the original data.
See chapter [Manipulating data](./manipulating.md) for details.

For more details on the record format, see appendix [File format](./file-format.md).
