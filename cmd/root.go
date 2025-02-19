package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ciselect",
	Short: "Interactive CLI tool for managing AWS CodePipeline approvals",
	Long: `ciselect is an interactive terminal UI for managing AWS CodePipeline manual approvals.
It provides a beautiful, user-friendly interface to:

- Select from available AWS profiles or type your own
- Choose from common regions or type a custom one
- View and select pending approvals in a clear, formatted list
- Approve or reject actions with guided prompts
- Confirm actions with clear context
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
