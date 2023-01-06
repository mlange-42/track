# Command tree

*Track*'s full command tree. See also `track --help`.

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
├─move
│ └─project PROJECT WORKSPACE
├─pause [NOTE...]
├─report
│ ├─chart [DATE]
│ ├─day [DATE]
│ ├─projects
│ ├─tags
│ ├─timeline (days|weeks|months)
│ ├─treemap
│ └─week [DATE]
├─resume [NOTE...]
├─start PROJECT [NOTE...]
├─status [PROJECT]
├─stop
├─switch PROJECT [NOTE...]
└─workspace WORKSPACE
```
