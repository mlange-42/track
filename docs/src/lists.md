# Lists

*Track*'s `list` command provides lists of different resources.

[[_TOC_]]

## Records

The `list records` command lists all records for the given day, or for the current day if no date is given:

```shell
track list records
track list records yesterday
track list records 2023-01-01
```

## Projects

The `list projects` command lists all projects as a tree showing the project hierarchy:

```shell
track list projects
```

Gives something like this:

```text
<default>
└─Private         P
  └─Coding        C
    └─MyApp       M
```

## Workspaces

The `list workspaces` command lists all available workspaces:

```shell
track list projects
```

## Tags

The `list tags` command lists all tags with their number of occurrences:

```shell
track list tags
```

## Colors

The `list colors` command shows all available colors for project configuration, with their indices:

```shell
track list colors
```

Shows this:

![Colors](./images/colors.png)
*Available colors with indices*
