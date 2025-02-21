package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the UI styling configurations
type Styles struct {
	App     lipgloss.Style
	Title   lipgloss.Style
	Help    lipgloss.Style
	Context lipgloss.Style
	Error   lipgloss.Style
	Table   table.Styles
}

// DefaultStyles returns the default styling configuration
func DefaultStyles() Styles {
	s := Styles{}

	// Define a modern color scheme
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	contextColor := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}

	s.App = lipgloss.NewStyle().
		Padding(1, 2).
		MaxWidth(100)

	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBottom(true).
		BorderForeground(subtle).
		Padding(0, 1)

	s.Help = lipgloss.NewStyle().
		Foreground(subtle).
		MarginTop(1)

	s.Context = lipgloss.NewStyle().
		Foreground(contextColor).
		Italic(true).
		PaddingLeft(4). // Increased padding for better visual hierarchy
		MaxWidth(100).  // Allow full width for context
		Height(6).      // Allow up to 6 lines for context
		Align(lipgloss.Left)

	s.Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	// Table styles
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		BorderBottom(true).
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Center)

	ts.Selected = ts.Selected.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(highlight).
		Bold(true).
		Padding(0, 1)

	// Ensure all cells have consistent padding
	ts.Cell = ts.Cell.
		Padding(0, 1)

	s.Table = ts

	return s
}
