package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func moveCommand(t *core.Track) *cobra.Command {
	var dryRun bool

	move := &cobra.Command{
		Use:     "move",
		Short:   "Move resources",
		Long:    `Move resources.`,
		Aliases: []string{"m"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	move.PersistentFlags().BoolVar(&dryRun, "dry", false, "Dry run: do not actually change any files")

	move.AddCommand(moveProjectCommand(t, &dryRun))

	move.Long += "\n\n" + formatCmdTree(move)
	return move
}

func moveProjectCommand(t *core.Track, dryRun *bool) *cobra.Command {

	moveProject := &cobra.Command{
		Use:     "project PROJECT WORKSPACE",
		Short:   "Move a project to another workspace",
		Long:    `Move a project to another workspace.`,
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ExactArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			workspace := args[1]

			project, err := t.LoadProjectByName(name)
			if err != nil {
				out.Err("failed to move project: %s", err)
				return
			}

			projects, err := t.LoadAllProjects()
			if err != nil {
				out.Err("failed to move project: %s", err)
				return
			}
			pTree, err := t.ToProjectTree(projects)
			if err != nil {
				out.Err("failed to move project: %s", err)
				return
			}
			pNode, ok := pTree.Nodes[name]
			if !ok {
				out.Err("failed to move project: no project named '%s'", name)
				return
			}
			if len(pNode.Children) > 0 {
				out.Err("failed to move project: '%s' has %d child project(s)", name, len(pNode.Children))
				return
			}

			if !t.WorkspaceExists(workspace) {
				out.Err("failed to move project: workspace '%s' does not exist", workspace)
				return
			}
			if t.Workspace() == workspace {
				out.Err("failed to move project: project is already in workspace '%s'", workspace)
				return
			}

			prevWorkspace := t.Config.Workspace
			t.Config.Workspace = workspace

			if t.ProjectExists(name) {
				out.Err("failed to move project: a project '%s' already exists in workspace '%s'", name, workspace)
				return
			}
			if t.ProjectExists(project.Parent) {
				project.Parent = ""
			}

			t.Config.Workspace = prevWorkspace

			filters := core.NewFilter(
				[]core.FilterFunction{
					core.FilterByProjects([]string{project.Name}),
				}, util.NoTime, util.NoTime,
			)

			records, err := t.LoadAllRecordsFiltered(filters)
			if err != nil {
				out.Err("failed to move project: %s", err)
				return
			}

			t.Config.Workspace = workspace

			if !*dryRun {
				err = t.SaveProject(project, false)
				if err != nil {
					out.Err("failed to move project: %s", err)
					return
				}

				for _, rec := range records {
					err = t.SaveRecord(&rec, false)
					if err != nil {
						out.Err("failed to move project: %s", err)
						return
					}
				}
			}

			t.Config.Workspace = prevWorkspace

			_, err = t.DeleteProject(&project, true, *dryRun)
			if err != nil {
				out.Err("failed to move project: %s", err)
				return
			}

			if *dryRun {
				out.Success("Moved project '%s' to workspace '%s' (%d records) - dry-run", name, workspace, len(records))
			} else {
				out.Success("Moved project '%s' to workspace '%s' (%d records)", name, workspace, len(records))
			}
		},
	}

	return moveProject
}
