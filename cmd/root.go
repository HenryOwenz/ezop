package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ezop",
	Short: "Interactive CLI tool for managing cloud operations",
	Long: `ezop is an interactive terminal UI for managing cloud operations across multiple providers.
It provides a beautiful, user-friendly interface to:

- Select from multiple cloud providers (AWS, Azure, GCP)
- Choose from available services for each provider
- Execute operations with guided prompts
- Approve or reject actions with clear context
- Confirm actions with safety checks
- Get immediate feedback with color-coded status updates`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel("", ""))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add completion command with custom documentation
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for ezop.
		
To load completions:

Bash:
  # Linux:
  $ ezop completion bash > /etc/bash_completion.d/ezop
  
  # macOS (requires bash-completion):
  $ ezop completion bash > $(brew --prefix)/etc/bash_completion.d/ezop

Zsh:
  # If shell completion is not already enabled in your environment:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  
  # Load the ezop completion code:
  $ ezop completion zsh > "${fpath[1]}/_ezop"

Fish:
  $ ezop completion fish > ~/.config/fish/completions/ezop.fish

PowerShell:
  PS> ezop completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> ezop completion powershell > ezop.ps1
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
