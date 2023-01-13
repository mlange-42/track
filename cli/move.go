package cli

import (
	"fmt"

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
		Use:   "project PROJECT WORKSPACE",
		Short: "Move a project to another workspace",
		Long: `Move a project to another workspace.

Moves the project and all associated records to the given workspace.
If there is no project with the same name as the parent of the project, the parent is set to none.`,
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ExactArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			workspace := args[1]

			project, err := t.LoadProject(name)
			if err != nil {
				return fmt.Errorf("failed to move project: %s", err)
			}

			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to move project: %s", err)
			}
			pTree, err := t.ToProjectTree(projects)
			if err != nil {
				return fmt.Errorf("failed to move project: %s", err)
			}
			pNode, ok := pTree.Nodes[name]
			if !ok {
				return fmt.Errorf("failed to move project: no project named '%s'", name)
			}
			if len(pNode.Children) > 0 {
				return fmt.Errorf("failed to move project: '%s' has %d child project(s)", name, len(pNode.Children))
			}

			if !t.WorkspaceExists(workspace) {
				return fmt.Errorf("failed to move project: workspace '%s' does not exist", workspace)
			}
			if t.Workspace() == workspace {
				return fmt.Errorf("failed to move project: project is already in workspace '%s'", workspace)
			}

			prevWorkspace := t.Config.Workspace
			t.Config.Workspace = workspace

			if t.ProjectExists(name) {
				return fmt.Errorf("failed to move project: a project '%s' already exists in workspace '%s'", name, workspace)
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
				return fmt.Errorf("failed to move project: %s", err)
			}

			t.Config.Workspace = workspace

			if !*dryRun {
				err = t.SaveProject(project, false)
				if err != nil {
					return fmt.Errorf("failed to move project: %s", err)
				}

				for _, rec := range records {
					err = t.SaveRecord(&rec, false)
					if err != nil {
						return fmt.Errorf("failed to move project: %s", err)
					}
				}
			}

			t.Config.Workspace = prevWorkspace

			_, err = t.DeleteProject(&project, true, *dryRun)
			if err != nil {
				return fmt.Errorf("failed to move project: %s", err)
			}

			if *dryRun {
				out.Success("Moved project '%s' to workspace '%s' (%d records) - dry-run", name, workspace, len(records))
			} else {
				out.Success("Moved project '%s' to workspace '%s' (%d records)", name, workspace, len(records))
			}
			return nil
		},
	}

	return moveProject
}
