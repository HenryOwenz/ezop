package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ciselect",
	Short: "A CLI tool for managing AWS CodePipeline approvals",
	Long: `ciselect is a command line tool that helps you manage AWS CodePipeline
manual approvals efficiently. It provides a simple interface to list and
approve/reject pending approvals in your pipelines.

Required AWS Configuration:
  - AWS Profile: Use --profile or -p to specify which AWS profile to use
  - AWS Region: Use --region or -r to specify which AWS region to use

Examples:
  # List all pending approvals
  ciselect list --profile prod-account --region us-west-2

  # Approve a specific action
  ciselect approve pipeline-name stage-name action-name \
    --profile prod-account \
    --region us-west-2 \
    --summary "Approved by John Doe"

  # Reject a specific action
  ciselect reject pipeline-name stage-name action-name \
    --profile prod-account \
    --region us-west-2 \
    --summary "Rejected due to test failures"

For more information and examples, visit:
https://github.com/HenryOwenz/ciselect`,
}

func init() {
	// Add completion command with custom documentation
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for ciselect.
		
To load completions:

Bash:
  # Linux:
  $ ciselect completion bash > /etc/bash_completion.d/ciselect
  
  # macOS (requires bash-completion):
  $ ciselect completion bash > $(brew --prefix)/etc/bash_completion.d/ciselect

Zsh:
  # If shell completion is not already enabled in your environment:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  
  # Load the ciselect completion code:
  $ ciselect completion zsh > "${fpath[1]}/_ciselect"

Fish:
  $ ciselect completion fish > ~/.config/fish/completions/ciselect.fish

PowerShell:
  PS> ciselect completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> ciselect completion powershell > ciselect.ps1
  # and source this file from your PowerShell profile.`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	rootCmd.AddCommand(completionCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
