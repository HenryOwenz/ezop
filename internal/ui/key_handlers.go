package ui

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m.handleKeyPressWithError(msg)
	}

	if m.isLoading {
		return m.handleKeyPressWhileLoading(msg)
	}

	if m.manualInput || m.currentView == constants.ViewSummary {
		return m.handleKeyPressInTextInput(msg)
	}

	return m.handleKeyPressInNormalMode(msg)
}

func (m *Model) handleKeyPressWithError(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "ctrl+c":
		return m, tea.Quit
	case "-":
		newModel := *m
		newModel.err = nil
		return newModel.navigateBack(), nil
	}
	return m, nil
}

func (m *Model) handleKeyPressWhileLoading(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, m.spinner.Tick
}

func (m *Model) handleKeyPressInTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		newModel := *m
		if m.currentView == constants.ViewSummary && m.selectedApproval != nil {
			// For approval summary, go back to confirmation
			newModel.currentView = constants.ViewConfirmation
			newModel.resetTextInput()
		} else {
			newModel.manualInput = false
		}
		newModel.updateTableForView()
		return &newModel, nil
	case "enter":
		if m.textInput.Value() != "" {
			newModel := *m
			if m.currentView == constants.ViewSummary {
				newModel.summary = m.textInput.Value()
				if m.selectedApproval != nil {
					// For approval summary, move to execution
					newModel.currentView = constants.ViewExecutingAction
					newModel.textInput.Blur()
				} else {
					// For pipeline start summary
					newModel.textInput.Blur()
					newModel.manualInput = false
					newModel.currentView = constants.ViewExecutingAction
				}
				newModel.updateTableForView()
				return &newModel, nil
			} else if m.awsProfile == "" {
				newModel.awsProfile = m.textInput.Value()
			} else {
				newModel.awsRegion = m.textInput.Value()
				newModel.currentView = constants.ViewSelectService
			}
			newModel.resetTextInput()
			newModel.manualInput = false
			newModel.updateTableForView()
			return &newModel, nil
		}
		return m, nil
	default:
		var tiCmd tea.Cmd
		m.textInput, tiCmd = m.textInput.Update(msg)
		return m, tiCmd
	}
}

func (m *Model) handleKeyPressInNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "-", "esc":
		if m.currentView > constants.ViewProviders {
			return m.navigateBack(), nil
		}
	case "tab":
		if m.currentView == constants.ViewAWSConfig || m.currentView == constants.ViewSummary {
			newModel := *m
			newModel.manualInput = !m.manualInput
			if newModel.manualInput {
				newModel.textInput.Focus()
				newModel.textInput.SetValue("")
			} else {
				newModel.textInput.Blur()
			}
			return &newModel, nil
		}
	case "enter":
		return m.handleEnter()
	}

	// Handle table navigation for non-input views
	if !m.manualInput && m.currentView != constants.ViewSummary {
		var tableCmd tea.Cmd
		newModel := *m
		newModel.table, tableCmd = m.table.Update(msg)
		return &newModel, tableCmd
	}

	return m, nil
} 