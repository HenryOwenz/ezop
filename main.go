package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/HenryOwenz/ezop/v2/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Clear screen based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()

	p := tea.NewProgram(ui.New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
