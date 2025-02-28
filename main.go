package main

import (
	"fmt"
	"os"

	// Import the AWS provider package to ensure its init() function is called
	_ "github.com/HenryOwenz/cloudgate/internal/providers/aws"

	"github.com/HenryOwenz/cloudgate/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Clear the screen using ANSI escape codes (works cross-platform)
	fmt.Print("\033[H\033[2J")

	// Create and run the program
	p := tea.NewProgram(ui.New())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
