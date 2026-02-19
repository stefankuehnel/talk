package main

import (
	"log"

	"stefanco.de/talk/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	err := rootCmd.Execute()
	if err != nil {
		logger := log.Default()
		logger.SetPrefix("talk: ")
		logger.SetOutput(rootCmd.ErrOrStderr())
		logger.SetFlags(0)

		logger.Fatalln(err)
	}
}
