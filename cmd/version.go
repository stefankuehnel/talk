package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is the version of the application.
// This variable can be set at build time using ldflags:
//
//	go build -ldflags "-X stefanco.de/talk/cmd.version=v1.2.3"
//
// If not set during build, it defaults to "dev".
var version = "dev"

// NewVersionCmd creates and returns the version command.
// This function returns a new instance of the version command for each execution,
// which is necessary for testing to avoid modified/dirty command instances
// in subsequent tests.
//
// The version information is written to cmd.OutOrStdout() to enable output
// capture in tests.
func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Write output to cmd.OutOrStdout() so it can be captured in tests
			_, err := fmt.Fprintln(cmd.OutOrStdout(), version)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return versionCmd
}
