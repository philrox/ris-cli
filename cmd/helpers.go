package cmd

import (
	"fmt"
	"os"

	"github.com/briandowns/spinner"
	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
	"github.com/philrox/ris-cli/internal/format"
	"github.com/philrox/ris-cli/internal/parser"
	"github.com/philrox/ris-cli/internal/ui"
	"github.com/spf13/cobra"
)

// ValidationError represents a user input validation error.
// main.go uses this to set exit code 2.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string { return e.msg }

// errValidation creates a validation error with fmt.Sprintf formatting.
func errValidation(format string, args ...any) error {
	return &ValidationError{msg: fmt.Sprintf(format, args...)}
}

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

// executeSearch runs the common search pipeline: spinner → API call → parse → output.
func executeSearch(cmd *cobra.Command, endpoint, spinnerMsg string, params *api.Params) error {
	setPageParams(cmd, params)

	client := newClient(cmd)
	s := startSpinner(cmd, spinnerMsg)
	body, err := client.Search(endpoint, params)
	stopSpinner(s)
	if err != nil {
		return fmt.Errorf("API-Anfrage fehlgeschlagen: %w", err)
	}

	result, err := parser.ParseSearchResponse(body)
	if err != nil {
		return fmt.Errorf("Antwort konnte nicht verarbeitet werden: %w", err)
	}

	if useJSON(cmd) {
		return format.JSON(os.Stdout, result)
	}
	return format.Text(os.Stdout, result)
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
		params.Set("DokumenteProSeite", constants.PageSizes[20])
	}
}
