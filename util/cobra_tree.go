package util

import (
	"fmt"

	"github.com/mlange-42/track/tree"
	"github.com/spf13/cobra"
)

// FormatCmdTree creates a tree-like representation of a command and its sub-commands
func FormatCmdTree(command *cobra.Command) (string, error) {
	cmdTree, err := NewCmdTree(command)
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

func nodePath(command *cobra.Command) string {
	if command.HasParent() {
		return fmt.Sprintf("%s/%s", nodePath(command.Parent()), command.Name())
	}
	return command.Name()
}

// NewCmdTree creates a new project tree
func NewCmdTree(command *cobra.Command) (*CmdTree, error) {

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
