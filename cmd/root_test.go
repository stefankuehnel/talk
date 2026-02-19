package cmd

import "testing"

func TestRootCmdSilencesUsageAndErrors(t *testing.T) {
	t.Parallel()

	rootCmd := NewRootCmd()

	if !rootCmd.SilenceUsage {
		t.Fatal("expected SilenceUsage to be true")
	}

	if !rootCmd.SilenceErrors {
		t.Fatal("expected SilenceErrors to be true")
	}
}
