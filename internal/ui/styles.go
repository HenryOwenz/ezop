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

	s.App = lipgloss.NewStyle().
		Padding(1, 2)

	s.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true).
		MarginBottom(1)

	s.Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	s.Context = lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		MarginTop(1)

	s.Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	// Table styles
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Center)

	ts.Selected = ts.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false).
		Padding(0, 1)

	// Ensure all cells have consistent padding
	ts.Cell = ts.Cell.
		Padding(0, 1)

	s.Table = ts

	return s
}
