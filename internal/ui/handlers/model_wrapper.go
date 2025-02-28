package handlers

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	tea "github.com/charmbracelet/bubbletea"
)

// ModelWrapper wraps a core.Model to implement the tea.Model interface
type ModelWrapper struct {
	Model *core.Model
}

// Update implements the tea.Model interface
func (m ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// This is just a placeholder - the actual update logic will be in the UI package
	return m, nil
}

// View implements the tea.Model interface
func (m ModelWrapper) View() string {
	// This is just a placeholder - the actual view logic will be in the UI package
	return ""
}

// Init implements the tea.Model interface
func (m ModelWrapper) Init() tea.Cmd {
	// This is just a placeholder - the actual init logic will be in the UI package
	return nil
}

// WrapModel wraps a core.Model in a ModelWrapper
func WrapModel(m *core.Model) ModelWrapper {
	return ModelWrapper{Model: m}
}
