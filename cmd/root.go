package cmd

import "github.com/spf13/cobra"

// NewRootCmd creates and returns the root command for the CLI application.
// This function should be used to create a new instance of the root command
// for each execution, ensuring a clean state for testing and avoiding
// modified/dirty command instances in subsequent tests.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "talk",
		Short:         "Send messages to Nextcloud Talk chats.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add child commands to the root command
	rootCmd.AddCommand(
		NewSendCmd(),
		NewVersionCmd(),
	)

	return rootCmd
}
