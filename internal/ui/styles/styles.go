package styles

import "github.com/charmbracelet/lipgloss"

// Styles holds all the UI styling configurations
type Styles struct {
	Title       lipgloss.Style
	Selected    lipgloss.Style
	Unselected  lipgloss.Style
	Instruction lipgloss.Style
	Error       lipgloss.Style
	Disabled    lipgloss.Style
	Loading     lipgloss.Style
}

// DefaultStyles returns the default styling configuration
func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(1, 0),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		Unselected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Instruction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Italic(true),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		Disabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")),
		Loading: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Italic(true),
	}
}
