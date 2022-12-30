package cli

import (
	"unicode/utf8"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func createCommand(t *core.Track) *cobra.Command {
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a new resource",
		Long:    "Create a new resource",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	create.AddCommand(createWorkspaceCommand(t))
	create.AddCommand(createProjectCommand(t))
	create.Long += "\n\n" + formatCmdTree(create)
	return create
}

func createProjectCommand(t *core.Track) *cobra.Command {
	var parent string
	var color uint8
	var fgColor uint8
	var symbol string

	createProject := &cobra.Command{
		Use:     "project PROJECT",
		Short:   "Create a new project",
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if !cmd.Flags().Changed("symbol") {
				symbol = string([]rune(name)[0])
			}
			if utf8.RuneCountInString(symbol) != 1 {
				out.Err("failed to create project: --symbol must be a single character")
				return
			}

			project := core.NewProject(name, parent, symbol, fgColor, color)

			if err := t.CheckParents(project); err != nil {
				out.Err("failed to create project: %s", err)
				return
			}

			if err := t.SaveProject(project, false); err != nil {
				out.Err("failed to create project: %s", err.Error())
				return
			}

			out.Success("Created project '%s'", name)
		},
	}

	createProject.Flags().StringVarP(&parent, "parent", "p", "", "Parent project of this project")
	createProject.Flags().Uint8VarP(&color, "color", "c", 0, "Background color for the project, as color index 0..256.\nSee: $ track list colors")
	createProject.Flags().Uint8VarP(&fgColor, "fg-color", "f", 15, "Foreground color for the project, as color index 0..256.\nSee: $ track list colors")
	createProject.Flags().StringVarP(&symbol, "symbol", "s", "", "Symbol for the project. Defaults to the first letter of the name")

	return createProject
}

func createWorkspaceCommand(t *core.Track) *cobra.Command {
	createWorkspace := &cobra.Command{
		Use:     "workspace WORKSPACE",
		Short:   "Create a new workspace",
		Aliases: []string{"w"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			err := t.CreateWorkspace(name)
			if err != nil {
				out.Err("failed to create workspace: %s", err.Error())
				return
			}

			out.Success("Created workspace '%s'", name)
		},
	}

	return createWorkspace
}
