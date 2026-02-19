package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Shell-Autovervollständigung generieren",
	Long: `Shell-Autovervollständigung für risgo generieren.

Bash:
  source <(risgo completion bash)

  # Dauerhaft installieren:
  risgo completion bash > /etc/bash_completion.d/risgo

Zsh:
  source <(risgo completion zsh)

  # Dauerhaft installieren:
  risgo completion zsh > "${fpath[1]}/_risgo"

Fish:
  risgo completion fish | source

  # Dauerhaft installieren:
  risgo completion fish > ~/.config/fish/completions/risgo.fish

PowerShell:
  risgo completion powershell | Out-String | Invoke-Expression`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
