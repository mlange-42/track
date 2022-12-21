package cli

import (
	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/spf13/cobra"
)

func createCommand(t *core.Track) *cobra.Command {
	create := &cobra.Command{
		Use:   "create",
		Short: "Create a new resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	create.AddCommand(createProjectCommand(t))

	return create
}

func createProjectCommand(t *core.Track) *cobra.Command {
	createProject := &cobra.Command{
		Use:   "project <NAME>",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			project := core.Project{
				Name: name,
			}

			if err := t.SaveProject(project); err != nil {
				out.Err("failed to create project: %s", err.Error())
				return
			}

			out.Success("Created project %s", name)
		},
	}

	return createProject
}
