package handlers

import (
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// NavigateBack handles navigation to the previous view
func NavigateBack(m *core.Model) *core.Model {
	newModel := m.Clone()

	switch m.CurrentView {
	case constants.ViewAWSConfig:
		if m.AwsProfile != "" {
			// If we're in region selection, just clear region and stay in AWS config
			newModel.AwsRegion = ""
			newModel.AwsProfile = ""
			// Don't change the view - we'll stay in AWS config to show profiles
		} else {
			// If we're in profile selection, go back to providers
			newModel.CurrentView = constants.ViewProviders
		}
		newModel.ManualInput = false
		newModel.ResetTextInput()
	case constants.ViewSelectService:
		newModel.CurrentView = constants.ViewAWSConfig
		newModel.SelectedService = nil
	case constants.ViewSelectCategory:
		newModel.CurrentView = constants.ViewSelectService
		newModel.SelectedCategory = nil
	case constants.ViewSelectOperation:
		newModel.CurrentView = constants.ViewSelectCategory
		newModel.SelectedOperation = nil
	case constants.ViewApprovals:
		newModel.CurrentView = constants.ViewSelectOperation
		newModel.ResetApprovalState()
	case constants.ViewConfirmation:
		newModel.CurrentView = constants.ViewApprovals
		newModel.SelectedApproval = nil
	case constants.ViewSummary:
		newModel.CurrentView = constants.ViewConfirmation
		newModel.Summary = ""
		newModel.ResetTextInput()
	case constants.ViewExecutingAction:
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			// For pipeline start flow, go back to pipeline selection
			newModel.CurrentView = constants.ViewPipelineStatus
			newModel.SelectedPipeline = nil
		} else {
			// For approval flow, go back to summary
			newModel.CurrentView = constants.ViewSummary
			// When going back to summary, restore the previous comment and focus
			newModel.TextInput.SetValue(m.Summary)
			newModel.TextInput.Focus()
			if newModel.ApproveAction {
				newModel.TextInput.Placeholder = constants.MsgEnterApprovalComment
			} else {
				newModel.TextInput.Placeholder = constants.MsgEnterRejectionComment
			}
		}
	case constants.ViewPipelineStages:
		newModel.CurrentView = constants.ViewPipelineStatus
		newModel.SelectedPipeline = nil
	case constants.ViewPipelineStatus:
		newModel.CurrentView = constants.ViewSelectOperation
		newModel.Pipelines = nil
		newModel.Provider = nil
	}
	return newModel
}

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
