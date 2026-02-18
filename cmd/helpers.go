package cmd

import (
	"fmt"

	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
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
