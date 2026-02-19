// Command docs generates man pages for ris.
//
// Usage:
//
//	go run ./cmd/docs --dir ./man/man1
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/philrox/ris-cli/cmd"
)

func main() {
	dir := "./man/man1"
	if len(os.Args) > 2 && os.Args[1] == "--dir" {
		dir = os.Args[2]
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	header := &doc.GenManHeader{
		Title:   "RIS",
		Section: "1",
		Source:  "ris-cli",
		Manual:  "RIS CLI Manual",
	}

	root := cmd.RootCmd()
	if err := doc.GenManTree(root, header, dir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating man pages: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Man pages generated in %s\n", dir)
}
