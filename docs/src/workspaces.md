# Workspaces

Workspaces are a high-level way to group [projects](./projects.md).

*Track* is always in one particular workspace, and the content of all other workspaces is invisible.
Workspaces are completely independent and separated from one another.

The default workspace is `default`.

To create another workspace, use:

```shell
track create MyWorkspace
```

Switch workspaces with command `workspace`:

```shell
track workspace MyWorkspace
```

To list all available workspaces, use

```shell
track list workspaces
```
