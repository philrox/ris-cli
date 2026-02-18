package main

import (
	"os"

	"github.com/philrox/ris-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
