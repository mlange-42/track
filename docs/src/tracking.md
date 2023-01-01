# Time tracking

*Track* provides several commands to start, stop, pause etc. time tracking.
They are explained in this chapter in detail.

[[_TOC_]]

## Projects

Each time tracking record is associated to a project. Before any tracking, a project needs to be created:

```shell
track create project MyProject
```

For more details on projects, see chapter [Projects](./projects.md).

## Start

To start tracking time on a project, use the `start` command:

```shell
track start MyProject
```

## Note and tags

Records can have a not and tags.
All positional arguments after the project's name are concatenated to the note text.
Words prefixed with '+' are extracted as tags.
Here is an example:

```shell
track start MyProject work on +artwork
```

## Status

To check the tracking status at any time, use:

```shell
track status
```

It will print a summary of the running or the last record:

```text
 SUCCESS  Record 2023-01-01 09:08
09:08 - ?
    - 10:25 - 10m / Short walk
    MyProject

work on +GUI +design
+------------------+-------+-------+-------+-------+
|          project |  curr | total | break | today |
|        MyProject | 02:05 | 02:05 | 00:10 | 02:53 |
+------------------+-------+-------+-------+-------+
```

## Stop

Command `stop` stops tracking:

```shell
track stop
```

## Pause

A record can contain multiple pause entries.
To insert a pause, use command `pause` with a duration:

```shell
track pause --duration 10m
```

This will insert a pause of 10 minutes, ending just now.
After the command, the record is not in paused mode.

To start a pause with an open end, use command `pause` without the duration option:

```shell
track pause
```

## Resume

To resume a paused record, use command `resume`:

```shell
track resume
```

The `resume` commands provides several flags:

* `--skip` to skip the running pause instead of closing it
* `--last` to resume an already finished record. Can be combined with `--skip`

## Switch

To switch to a different project, command `switch` can be used instead of successive `stop` and `start`:

```shell
track switch MyProject
```

Notes and tags apply here just as with `start`.

## Time corrections

For the case that you did not start, stop, pause etc. at the correct time, all commands described in this chapter have flags to correct time:

* `--at` foredates the command to the given time  
  ```shell
  track start --at 14:00
  ```
* `--ago` foredates the command by the given duration  
  ```shell
  track start --ago 10m
  ```

These flags are mutually exclusive.
