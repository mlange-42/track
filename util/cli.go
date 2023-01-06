package util

import (
	"fmt"
	"os"

	"github.com/mlange-42/track/tree"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// TerminalSize returns the size of the terminal
func TerminalSize() (width int, height int, err error) {
	return terminal.GetSize(int(os.Stdout.Fd()))
}

// WrappedArgs are PositionalArgs that print usage on error
func WrappedArgs(fn cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := fn(cmd, args)
		if err != nil {
			return fmt.Errorf("%s\nUsage: %s", err, cmd.UseLine())
		}
		return nil
	}
}

// FormatCmdTree creates a tree-like representation of a command and its sub-commands
func FormatCmdTree(command *cobra.Command) (string, error) {
	cmdTree, err := newCmdTree(command)
	if err != nil {
		return "", err
	}

	formatter := NewTreeFormatter(
		func(t *CmdNode, indent int) string {
			return t.Value.Use
		},
		2,
	)
	return formatter.FormatTree(cmdTree), nil
}

// CmdTree is a tree of cobra commands
type CmdTree = tree.MapTree[*cobra.Command]

// CmdNode is a tree of cobra commands
type CmdNode = tree.MapNode[*cobra.Command]

// NewCmdTree creates a new project tree
func newCmdTree(command *cobra.Command) (*CmdTree, error) {

	t := tree.NewTree(
		command,
		func(c *cobra.Command) string { return nodePath(c) },
	)

	err := buildTree(t, t.Root)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func nodePath(command *cobra.Command) string {
	if command.HasParent() {
		return fmt.Sprintf("%s/%s", nodePath(command.Parent()), command.Name())
	}
	return command.Name()
}

func buildTree(t *CmdTree, node *tree.MapNode[*cobra.Command]) error {
	for _, cmd := range node.Value.Commands() {
		child, err := t.Add(node, cmd)
		if err != nil {
			return err
		}
		err = buildTree(t, child)
		if err != nil {
			return err
		}
	}
	return nil
}
