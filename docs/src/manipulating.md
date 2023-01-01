# Manipulating data

Besides the normal tracking workflow with `start`, `stop`, `pause` etc.,
*Track* provides a simple but efficient way to edit all kinds of underlying data:

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

This should open the config file and wait for the user to edit, save and close it.
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

[TODO]

## Editing projects

[TODO]

## Archiving projects

[TODO]

## Deleting records and projects

[TODO]
