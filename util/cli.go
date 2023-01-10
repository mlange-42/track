package util

import (
	"fmt"
	"os"

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
type CmdTree = MapTree[CmdWrapper]

// CmdNode is a tree of cobra commands
type CmdNode = MapNode[CmdWrapper]

// CmdWrapper wraps *cobra.Command to implement the Named interface
type CmdWrapper struct {
	*cobra.Command
}

// GetName implements the Named interface required for the MapTree
func (cmd CmdWrapper) GetName() string {
	return nodePath(cmd.Command)
}

// NewCmdTree creates a new project tree
func newCmdTree(command *cobra.Command) (*CmdTree, error) {

	t := NewTree(
		CmdWrapper{command},
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

func buildTree(t *CmdTree, node *MapNode[CmdWrapper]) error {
	for _, cmd := range node.Value.Commands() {
		child, err := t.Add(node, CmdWrapper{cmd})
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
