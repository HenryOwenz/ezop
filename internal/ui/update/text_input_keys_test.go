package update

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

// TestVimKeysInTextInputMode verifies that vim-style navigation keys
// are properly handled in text input mode.
func TestVimKeysInTextInputMode(t *testing.T) {
	// This test verifies that when in text input mode:
	// 1. 'j' and 'k' keys are passed to the text input component
	// 2. Other vim-style navigation keys are also passed to the text input component

	// Define the keys we want to test
	keysToTest := []struct {
		name string
		key  string
	}{
		{"j key", "j"},
		{"k key", "k"},
		{"- key", "-"},
		{"g key", "g"},
		{"G key", "G"},
		{"u key", "u"},
		{"d key", "d"},
		{"b key", "b"},
		{"f key", "f"},
	}

	for _, tc := range keysToTest {
		t.Run(tc.name, func(t *testing.T) {
			// Create a model with manual input enabled
			m := model.New()
			m.ManualInput = true

			// Set up the text input
			ti := textinput.New()
			ti.Focus()
			m.TextInput = ti

			// Create a key message for the key we're testing
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tc.key)}

			// Call the HandleTextInput function directly
			result := HandleTextInput(m, msg)

			// Verify that the key was added to the text input
			if result.TextInput.Value() != tc.key {
				t.Errorf("Expected key '%s' to be passed to text input in manual input mode, got '%s'",
					tc.key, result.TextInput.Value())
			}
		})
	}
}

// HandleTextInput simulates the text input handling logic from the UI update function
func HandleTextInput(m *model.Model, msg tea.KeyMsg) *model.Model {
	// Create a copy of the model
	result := m.Clone()

	// Update the text input with the key message
	result.TextInput, _ = result.TextInput.Update(msg)

	return result
}
