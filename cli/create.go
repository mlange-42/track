package cli

import (
	"fmt"
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
	var requiredTags []string
	var color uint8
	var fgColor uint8
	var symbol string

	createProject := &cobra.Command{
		Use:     "project PROJECT",
		Short:   "Create a new project",
		Aliases: []string{"p"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if !cmd.Flags().Changed("symbol") {
				symbol = string([]rune(name)[0])
			}
			if utf8.RuneCountInString(symbol) != 1 {
				return fmt.Errorf("failed to create project: --symbol must be a single character")
			}

			requiredTags = util.Unique(requiredTags)
			project := core.NewProject(name, parent, symbol, requiredTags, fgColor, color)

			if err := t.CheckParents(project); err != nil {
				return fmt.Errorf("failed to create project: %s", err)
			}

			if err := t.SaveProject(project, false); err != nil {
				return fmt.Errorf("failed to create project: %s", err.Error())
			}

			out.Success("Created project '%s'", name)
			return nil
		},
	}

	createProject.Flags().StringVarP(&parent, "parent", "p", "", "Parent project of this project")
	createProject.Flags().StringSliceVarP(&requiredTags, "tags", "t", []string{}, "Tags that are required for records in this project")
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
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			err := t.CreateWorkspace(name)
			if err != nil {
				return fmt.Errorf("failed to create workspace: %s", err.Error())
			}

			out.Success("Created workspace '%s'", name)
			return nil
		},
	}

	return createWorkspace
}
