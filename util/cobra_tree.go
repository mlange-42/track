package util

import (
	"github.com/mlange-42/track/tree"
	"github.com/spf13/cobra"
)

// FormatCmdTree creates a tree-like representation of a command and its sub-commands
func FormatCmdTree(command *cobra.Command) string {
	cmdTree := NewCmdTree(command)

	formatter := NewTreeFormatter(
		func(t *CmdNode, indent int) string {
			return t.Value.Use
		},
		2,
	)
	return formatter.FormatTree(cmdTree)
}

// CmdTree is a tree of cobra commands
type CmdTree = tree.MapTree[*cobra.Command]

// CmdNode is a tree of cobra commands
type CmdNode = tree.MapNode[*cobra.Command]

// NewCmdTree creates a new project tree
func NewCmdTree(command *cobra.Command) *CmdTree {
	t := tree.NewTree(
		command,
		func(c *cobra.Command) string { return c.Name() },
	)

	buildTree(t, t.Root)
	return t
}

func buildTree(t *CmdTree, node *tree.MapNode[*cobra.Command]) {
	for _, cmd := range node.Value.Commands() {
		child := t.Add(node, cmd)
		buildTree(t, child)
	}
}
