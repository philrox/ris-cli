package ui

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-isatty"
)

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }

// NewPagerWriter returns a writer that pipes output through a pager (e.g. less)
// when stdout is a TTY, and a cleanup function that must be called when done.
// If stdout is not a TTY or the pager cannot be started, it falls back to stdout.
func NewPagerWriter(noPager bool) (io.WriteCloser, func()) {
	if noPager || !isatty.IsTerminal(os.Stdout.Fd()) {
		return nopCloser{os.Stdout}, func() {}
	}

	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "less -FIRX"
	}

	parts := strings.Fields(pagerCmd)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	w, err := cmd.StdinPipe()
	if err != nil {
		return nopCloser{os.Stdout}, func() {}
	}

	if err := cmd.Start(); err != nil {
		return nopCloser{os.Stdout}, func() {}
	}

	cleanup := func() {
		w.Close()
		cmd.Wait()
	}

	return w, cleanup
}
