package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/HenryOwenz/cloudgate/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Clear the screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to clear screen: %v\n", err)
	}

	// Create and run the program
	p := tea.NewProgram(ui.New())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
