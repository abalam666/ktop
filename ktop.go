package main

import (
	"os"

	"github.com/ynqa/ktop/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
