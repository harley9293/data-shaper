package main

import (
	"github.com/harley9293/data-shaper/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
