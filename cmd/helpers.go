package cmd

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
	"github.com/philrox/ris-cli/internal/ui"
	"github.com/spf13/cobra"
)

// newClient creates an API client from the root command's global flags.
func newClient(cmd *cobra.Command) *api.Client {
	root := cmd.Root()
	timeout, _ := root.PersistentFlags().GetDuration("timeout")
	verbose, _ := root.PersistentFlags().GetBool("verbose")

	return api.NewClient(api.ClientOptions{
		Timeout: timeout,
		Verbose: verbose,
	})
}

// useJSON returns true if --json flag is set on the root command.
func useJSON(cmd *cobra.Command) bool {
	root := cmd.Root()
	j, _ := root.PersistentFlags().GetBool("json")
	return j
}

// isVerbose returns true if --verbose flag is set.
func isVerbose(cmd *cobra.Command) bool {
	root := cmd.Root()
	v, _ := root.PersistentFlags().GetBool("verbose")
	return v
}

// usePlain returns true if plain text output should be used.
// This is the case when --plain is set or stdout is not a TTY.
func usePlain(cmd *cobra.Command) bool {
	root := cmd.Root()
	p, _ := root.PersistentFlags().GetBool("plain")
	return p || !isTTY
}

// IsTTY reports whether stdout is connected to a terminal.
func IsTTY() bool {
	return isTTY
}

// startSpinner starts a progress spinner on stderr if conditions allow it.
// Returns nil if spinner should not be shown (JSON mode, quiet, non-TTY).
func startSpinner(cmd *cobra.Command, msg string) *spinner.Spinner {
	if !isTTY || useJSON(cmd) || quiet {
		return nil
	}
	s := ui.NewSpinner(msg)
	s.Start()
	return s
}

// stopSpinner stops a running spinner. Safe to call with nil.
func stopSpinner(s *spinner.Spinner) {
	if s != nil {
		s.Stop()
	}
}

// setPageParams sets pagination parameters from the root command's global flags.
func setPageParams(cmd *cobra.Command, params *api.Params) {
	root := cmd.Root()
	page, _ := root.PersistentFlags().GetInt("page")
	limit, _ := root.PersistentFlags().GetInt("limit")

	if page > 0 {
		params.Set("Seitennummer", fmt.Sprintf("%d", page))
	}

	pageSize, ok := constants.PageSizes[limit]
	if ok {
		params.Set("DokumenteProSeite", pageSize)
	} else {
		params.Set("DokumenteProSeite", "Twenty")
	}
}
