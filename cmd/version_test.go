package cmd

import (
	"bytes"
	"fmt"
	"testing"
)

// TestVersionCmd tests the version command functionality.
// Important testing practices for Cobra commands:
//   - Always use the parent (root) command for setting arguments, output, etc.
//     When you call Execute() on the root command, Cobra processes the arguments
//     and delegates to the appropriate subcommand. Setting arguments directly on
//     a child command bypasses this delegation mechanism.
//   - Use NewRootCmd() to generate a new command instance for each test,
//     ensuring clean state and avoiding modified/dirty command instances.
//   - Capture output by setting cmd.SetOut() and cmd.SetErr() on the root command.
func TestVersionCmd(t *testing.T) {
	t.Parallel()

	t.Run("with no args", func(t *testing.T) {
		t.Parallel()

		// Create a new root command instance for this test
		rootCmd := NewRootCmd()

		// Create a buffer to capture output
		out := new(bytes.Buffer)

		// Set output and error streams on the root command
		rootCmd.SetOut(out)
		rootCmd.SetErr(out)

		// Set the arguments on the root command (not the child command)
		rootCmd.SetArgs([]string{"version"})

		// Execute the command
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("got error %q", err)
		}

		// Verify the output
		outStr := out.String()
		expectedText := fmt.Sprintf("%s\n", version)
		if outStr != expectedText {
			t.Errorf("expected: %q; got: %q", expectedText, outStr)
		}
	})
}
