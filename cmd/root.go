package cmd

import (
	"fmt"
	"os"

	"github.com/HenryOwenz/ezop/internal/ui/handlers"
	"github.com/HenryOwenz/ezop/internal/ui/model"
	"github.com/HenryOwenz/ezop/internal/ui/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type program struct {
	model model.Model
}

func (p program) Init() tea.Cmd {
	return nil
}

func (p program) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	p.model, cmd = handlers.Update(p.model, msg)
	return p, cmd
}

func (p program) View() string {
	return views.View(p.model)
}

var rootCmd = &cobra.Command{
	Use:   "ezop",
	Short: "A user-friendly interactive CLI tool for managing cloud operations",
	Run: func(cmd *cobra.Command, args []string) {
		p := program{
			model: model.NewModel(),
		}

		if _, err := tea.NewProgram(p).Run(); err != nil {
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

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v", err)
		os.Exit(1)
	}
}
