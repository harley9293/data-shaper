package main

import (
	"github.com/harley9293/data-shaper/internal/pbz"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "shaper",
		Short: "for the conversion and generation between various types of file data",
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(pbz.Cmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
