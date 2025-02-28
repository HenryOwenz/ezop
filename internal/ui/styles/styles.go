package styles

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
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

	// Use color constants for consistent styling
	subtle := lipgloss.Color(constants.ColorSubtle)
	highlight := lipgloss.Color(constants.ColorPrimary)
	special := lipgloss.Color(constants.ColorSuccess)
	contextColor := lipgloss.Color(constants.ColorSubtle)
	darkGray := lipgloss.Color(constants.ColorBgAlt)
	titleColor := lipgloss.Color(constants.ColorTitle)
	headerColor := lipgloss.Color(constants.ColorHeader)

	s.App = lipgloss.NewStyle().
		Padding(1, 2).
		MaxWidth(100).
		Height(17) // Increased to accommodate more context

	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(titleColor).
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
		Foreground(lipgloss.Color(constants.ColorError)).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(constants.ColorError))

	// Table styles with fixed height
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(special).
		BorderBottom(true).
		Bold(true).
		Padding(0, 1).
		Foreground(headerColor).
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
