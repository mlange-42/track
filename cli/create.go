package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/spf13/cobra"
)

func createCommand(t *core.Track) *cobra.Command {
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a new resource",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	create.AddCommand(createProjectCommand(t))

	return create
}

func createProjectCommand(t *core.Track) *cobra.Command {
	var parent string

	createProject := &cobra.Command{
		Use:     "project <NAME>",
		Short:   "Create a new project",
		Aliases: []string{"p"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if parent != "" && !t.ProjectExists(parent) {
				out.Err("failed to create project: parent project '%s' does not exist", parent)
				return
			}

			project := core.Project{
				Name:   name,
				Parent: parent,
			}

			if err := t.SaveProject(project, false); err != nil {
				out.Err("failed to create project: %s", err.Error())
				return
			}

			out.Success("Created project '%s'", name)
		},
	}

	createProject.Flags().StringVarP(&parent, "parent", "p", "", "Parent project of this project")

	return createProject
}
