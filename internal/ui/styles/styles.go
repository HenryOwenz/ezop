package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the UI styling configurations
type Styles struct {
	// Layout
	Frame       lipgloss.Style // Main window frame
	ContentArea lipgloss.Style // Content area inside frame
	Section     lipgloss.Style // Individual sections within content

	// Text and UI elements (keeping our existing colors)
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
	// Get terminal size
	w, _ := lipgloss.Size("")

	// Set frame width to 90% of terminal width
	frameWidth := int(float64(w) * 0.9)

	// Calculate content area size
	contentWidth := frameWidth - 4 // Account for frame borders and padding

	return Styles{
		// Layout styles
		Frame: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(1).
			MarginLeft(2).
			MarginRight(2).
			Width(frameWidth),

		ContentArea: lipgloss.NewStyle().
			Width(contentWidth).
			MarginBottom(1),

		Section: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(0, 1).
			MarginTop(1),

		// Text styles without center alignment
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			MarginBottom(1).
			Padding(1, 0),

		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),

		Unselected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),

		Instruction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Italic(true).
			MarginTop(1),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),

		Disabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")),

		Loading: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Italic(true).
			MarginBottom(1),
	}
}

// GetWindowSize returns the terminal window size
func GetWindowSize() (width int, height int) {
	w, h := lipgloss.Size("")
	return w, h
}

// JoinHorizontal joins strings horizontally with the given width
func JoinHorizontal(width int, strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

// JoinVertical joins strings vertically with the given height
func JoinVertical(height int, strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, strs...)
}

// Place positions a string at the given coordinates
func Place(width, height int, position lipgloss.Position, str string) string {
	return lipgloss.Place(width, height, position, position, str)
}
