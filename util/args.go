package util

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
