package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/philrox/ris-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		var ve *cmd.ValidationError
		if errors.As(err, &ve) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
