package ui

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

// NewSpinner creates a configured spinner that writes to stderr.
func NewSpinner(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " " + msg
	return s
}
