# Configuration

*Track* can be configured for the user's needs.
Most of this configuration resides in a file `config.yaml` in *Track*'s data directory.

[[_TOC_]]

## Data directory

The default data directory is `%USER%/.track`. On Windows, this resolves to `C:\Users\<USER>\.track\`.

The data directory can be changed by setting the environmental variable `TRACK_PATH`.

## Config file

*Track*'s configuration is stored in a file `config.yaml` in the data directory.
See chapter [Manipulating data](./manipulating.md) for editing the config file.

A config file has the following content:

```yaml
# Track config
workspace: default
textEditor: nano
maxBreakDuration: 2h0m0s
emptyCell: .
pauseCell: '-'
```

* `workspace` - *Track*'s current workspace.
* `textEditor` - The text editor to use for editing records etc. Default value system-dependent.
* `maxBreakDuration` - Maximum duration of interruptions of a project to count as ongoing with a break.
* `emptyCell` - Character for empty cells in schedule-like reports (`report week` and `report day`).
* `pauseCell` - Character for pause cells in schedule-like reports (`report week` and `report day`).
