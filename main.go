package main

import (
	"os"

	"github.com/harley9293/data-shaper/internal/pbz"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "shaper",
		Short: "用于各种文件格式之间的转换和生成",
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(pbz.Cmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
