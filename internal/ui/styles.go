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

	// Define a green-based professional scheme with orange highlight
	primary := lipgloss.AdaptiveColor{Light: "#276749", Dark: "#48BB78"}      // Forest green
	secondary := lipgloss.AdaptiveColor{Light: "#1C4532", Dark: "#68D391"}    // Deep/Light green
	subtle := lipgloss.AdaptiveColor{Light: "#718096", Dark: "#CBD5E0"}       // Keep gray for subtle elements
	highlight := lipgloss.AdaptiveColor{Light: "#DD6B20", Dark: "#ED8936"}    // Professional orange
	special := lipgloss.AdaptiveColor{Light: "#2F855A", Dark: "#68D391"}      // Medium green
	contextColor := lipgloss.AdaptiveColor{Light: "#4A5568", Dark: "#A0AEC0"} // Keep gray for context
	darkGray := lipgloss.AdaptiveColor{Light: "#1A202C", Dark: "#2D3748"}     // Dark gray for selected text

	s.App = lipgloss.NewStyle().
		Padding(1, 2).
		MaxWidth(100).
		Height(17) // Increased to accommodate more context

	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(primary).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(subtle).
		Padding(0, 1).
		Height(1)

	s.Help = lipgloss.NewStyle().
		Foreground(subtle).
		MarginTop(1).
		Height(1)

	s.Context = lipgloss.NewStyle().
		Foreground(contextColor).
		Padding(0, 1).
		Height(5) // Increased from 3 to 5 lines for context

	s.Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E53E3E")).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#E53E3E"))

	// Table styles with fixed height
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(special).
		BorderBottom(true).
		Bold(true).
		Padding(0, 1).
		Foreground(secondary).
		Align(lipgloss.Center)

	ts.Selected = ts.Selected.
		Foreground(darkGray).
		Background(highlight).
		Bold(true).
		Padding(0, 1)

	ts.Cell = ts.Cell.
		BorderForeground(subtle).
		Padding(0, 1)

	s.Table = ts

	return s
}
