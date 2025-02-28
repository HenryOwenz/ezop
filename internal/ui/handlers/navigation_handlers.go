package handlers

import (
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleEnter processes the enter key press based on the current view
func HandleEnter(m *core.Model) (tea.Model, tea.Cmd) {
	// Special handling for manual input in AWS config view
	if m.CurrentView == constants.ViewAWSConfig && m.ManualInput {
		newModel := m.Clone()

		// Get the entered value
		value := strings.TrimSpace(m.TextInput.Value())
		if value == "" {
			// If empty, just exit manual input mode
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}

		// Set the appropriate value based on context
		if m.AwsProfile == "" {
			// Setting profile
			newModel.AwsProfile = value
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
		} else {
			// Setting region and moving to next view
			newModel.AwsRegion = value
			newModel.ManualInput = false
			newModel.ResetTextInput()

			// Create the provider with the selected profile and region
			_, err := providers.CreateProvider(newModel.Registry, "AWS", newModel.AwsProfile, newModel.AwsRegion)
			if err != nil {
				return WrapModel(newModel), func() tea.Msg {
					return core.ErrMsg{Err: err}
				}
			}

			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(newModel)
		}

		return WrapModel(newModel), nil
	}

	// Regular view handling
	switch m.CurrentView {
	case constants.ViewProviders:
		return HandleProviderSelection(m)
	case constants.ViewAWSConfig:
		return HandleAWSConfigSelection(m)
	case constants.ViewSelectService:
		return HandleServiceSelection(m)
	case constants.ViewSelectCategory:
		return HandleCategorySelection(m)
	case constants.ViewSelectOperation:
		return HandleOperationSelection(m)
	case constants.ViewApprovals:
		return HandleApprovalSelection(m)
	case constants.ViewConfirmation:
		return HandleConfirmationSelection(m)
	case constants.ViewSummary:
		if !m.ManualInput {
			if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
				if selected := m.Table.SelectedRow(); len(selected) > 0 {
					newModel := m.Clone()
					switch selected[0] {
					case "Latest Commit":
						newModel.CurrentView = constants.ViewExecutingAction
						newModel.Summary = "" // Empty string means use latest commit
						view.UpdateTableForView(newModel)
						return WrapModel(newModel), nil
					case "Manual Input":
						newModel.ManualInput = true
						newModel.TextInput.Focus()
						newModel.TextInput.Placeholder = constants.MsgEnterCommitID
						return WrapModel(newModel), nil
					}
				}
			}
		}
		return HandleSummaryConfirmation(m)
	case constants.ViewExecutingAction:
		return HandleExecutionSelection(m)
	case constants.ViewPipelineStatus:
		if selected := m.Table.SelectedRow(); len(selected) > 0 {
			newModel := m.Clone()
			for _, pipeline := range m.Pipelines {
				if pipeline.Name == selected[0] {
					if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
						newModel.CurrentView = constants.ViewExecutingAction
						newModel.SelectedPipeline = &pipeline
						view.UpdateTableForView(newModel)
						return WrapModel(newModel), nil
					}
					newModel.CurrentView = constants.ViewPipelineStages
					newModel.SelectedPipeline = &pipeline
					view.UpdateTableForView(newModel)
					return WrapModel(newModel), nil
				}
			}
		}
	case constants.ViewPipelineStages:
		// Just view only, no action
	}
	return WrapModel(m), nil
}

// HandleProviderSelection handles the selection of a cloud provider
func HandleProviderSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		providerName := selected[0]

		// Initialize providers if not already done
		if len(m.Registry.GetProviderNames()) == 0 {
			providers.InitializeProviders(m.Registry)
		}

		// Check if the provider exists in the registry
		_, exists := m.Registry.GetProvider(providerName)
		if exists {
			newModel := m.Clone()
			newModel.CurrentView = constants.ViewAWSConfig
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleAWSConfigSelection handles the selection of AWS profile or region
func HandleAWSConfigSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()

		// Handle "Manual Entry" option
		if selected[0] == "Manual Entry" {
			newModel.ManualInput = true
			newModel.TextInput.Focus()

			// Set appropriate placeholder based on context
			if m.AwsProfile == "" {
				newModel.TextInput.Placeholder = constants.MsgEnterProfile
			} else {
				newModel.TextInput.Placeholder = constants.MsgEnterRegion
			}

			return WrapModel(newModel), nil
		}

		// Handle regular selection
		if m.AwsProfile == "" {
			newModel.AwsProfile = selected[0]
			view.UpdateTableForView(newModel)
		} else {
			newModel.AwsRegion = selected[0]

			// Create the provider with the selected profile and region
			_, err := providers.CreateProvider(newModel.Registry, "AWS", newModel.AwsProfile, newModel.AwsRegion)
			if err != nil {
				return WrapModel(newModel), func() tea.Msg {
					return core.ErrMsg{Err: err}
				}
			}

			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(newModel)
		}
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}
