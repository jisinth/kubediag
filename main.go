package main

import (
	"os"

	"github.com/jisinth/kubediag/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
