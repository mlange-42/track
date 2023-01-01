# Manipulating data

Besides the normal tracking workflow with `start`, `stop`, `pause` etc.,
*Track* provides a simple but effective way to edit all kinds of underlying data.

All data is stored in human-readable formats.
Thus, it can be edited by letting *Track* open it as a temporary file in a text editor.
After editing the file, the user closes it to confirm.
*Track* checks the data for consistency, and only then overwrites any files.

[[_TOC_]]

## Editing the config

For *Track*'s editing to work properly, the text editor to be used must be set in the config file.
On Windows, `notepad.exe` is set as the default, while it is `nano` on other systems.
You can try if the setup works for you by editing the config:

```shell
track edit config
```

This should open the config YAML file and wait for the user to edit, save and close it.
If this does not work properly, set the editor manually. Open the config file under

```text
%USER%/.track/config.yml
```

There, set the text editor entry to a program of your choice, e.g.:

```yaml
textEditor: vim
```

Then, save the file and try to edit using *Track* again with `track edit config`.

## Editing records

There are two ways for editing records:

* **Edit a single record** with `track edit record [[DATE] TIME]`  
  Good for changing a record's project, note or pauses, but limits editing of start and end time of the record.

* **Edit all records of a day** in a single file with `track edit day [DATE]`  
  Allows for changing start and end times, in addition to the other properties.
  Also checks consistency between records (no overlap etc.).

The file format/syntax for editing records should be quite obvious.
It is the same format that *Track* uses to store records.

When editing a full day, records are separated by lines starting with 4 dashes: `----`.

For details on the file format, see appendix [File formats](./file-formats.md).

## Editing projects

Projects can be edited just like the config or records:

```shell
track edit project MyProject
```

Projects are stored and edited as YAML files. For details, see chapter [Projects](./projects.md).

## Archiving projects

Projects can be archived.
Archiving a project has no effects except that the project is excluded from lists and reports.

Some commands have a flag `--archived` to include archived projects.

To archive on un-archive a projects, use the `--archive` flag:

```shell
track edit project --archive
track edit project --archive=false
```

## Deleting records and projects

Records and entire projects (including all their records) can be deleted using the CLI.

Delete a record:

```shell
track delete record 2023-01-01 15:05
```

Delete a project, including all records of the project:

```shell
track delete project MyProject
```

The `delete` commands ask for user confirmation before actually deleting anything.
