package update

import (
	"sort"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// NavigateBack handles navigation to the previous view
func NavigateBack(m *model.Model) *model.Model {
	newModel := m.Clone()

	switch m.CurrentView {
	case constants.ViewAuthMethodSelect:
		// Go back to provider selection
		newModel.CurrentView = constants.ViewProviders
		newModel.ProviderState.ProviderName = ""
	case constants.ViewAuthConfig:
		// Go back to auth method selection or provider selection
		if len(m.ProviderState.AuthState.AvailableMethods) > 1 {
			newModel.CurrentView = constants.ViewAuthMethodSelect
		} else {
			newModel.CurrentView = constants.ViewProviders
			newModel.ProviderState.ProviderName = ""
		}
	case constants.ViewProviderConfig:
		// Go back to auth config, auth method selection, or provider selection
		if len(m.ProviderState.AuthState.AvailableMethods) > 0 {
			if len(m.ProviderState.AuthState.AuthConfig) > 0 {
				newModel.CurrentView = constants.ViewAuthConfig
			} else if len(m.ProviderState.AuthState.AvailableMethods) > 1 {
				newModel.CurrentView = constants.ViewAuthMethodSelect
			} else {
				newModel.CurrentView = constants.ViewProviders
				newModel.ProviderState.ProviderName = ""
			}
		} else {
			newModel.CurrentView = constants.ViewProviders
			newModel.ProviderState.ProviderName = ""
		}
	case constants.ViewAWSConfig:
		if m.GetAwsProfile() != "" {
			// If we're in region selection, just clear region and stay in AWS config
			newModel.SetAwsRegion("")
			newModel.SetAwsProfile("")
			// Don't change the view - we'll stay in AWS config to show profiles
		} else {
			// If we're in profile selection, go back to providers
			newModel.CurrentView = constants.ViewProviders
		}
		newModel.ManualInput = false
		newModel.ResetTextInput()
	case constants.ViewSelectService:
		// For backward compatibility, go back to AWS config
		// In the future, this will go back to provider config
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
		// For pipeline start flow, go back to pipeline status view
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			newModel.CurrentView = constants.ViewPipelineStatus
		} else {
			// For approval flow, go back to confirmation view
			newModel.CurrentView = constants.ViewConfirmation
		}
		newModel.Summary = ""
		newModel.ResetTextInput()
	case constants.ViewExecutingAction:
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			// For pipeline start flow, go back to pipeline status view
			newModel.CurrentView = constants.ViewPipelineStatus

			// Make sure we're showing the pipeline selection table, not the approval table
			// This ensures we stay in the pipeline start flow, not the approval flow
			if newModel.SelectedPipeline != nil {
				// Reset any approval-related state that might be present
				newModel.SelectedApproval = nil
				newModel.ApproveAction = false
				newModel.ApprovalComment = ""

				// Make sure we're in the correct view with the right table
				view.UpdateTableForView(newModel)
			}
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
	case constants.ViewFunctionStatus:
		newModel.CurrentView = constants.ViewSelectOperation
		newModel.Functions = nil
		newModel.Provider = nil
	case constants.ViewFunctionDetails:
		newModel.CurrentView = constants.ViewFunctionStatus
		newModel.SetSelectedFunction(nil)
	}
	return newModel
}

// HandleTableSelect handles table row selection based on the current view
func HandleTableSelect(m *model.Model) (tea.Model, tea.Cmd) {
	switch m.CurrentView {
	case constants.ViewProviders:
		return HandleProviderSelection(m)
	case constants.ViewAWSConfig:
		return HandleAWSConfigSelection(m)
	case constants.ViewAuthMethodSelect:
		return HandleAuthMethodSelection(m)
	case constants.ViewAuthConfig:
		return HandleAuthConfigSelection(m)
	case constants.ViewProviderConfig:
		return HandleProviderConfigSelection(m)
	case constants.ViewSelectService:
		return SelectService(m)
	case constants.ViewSelectCategory:
		return SelectCategory(m)
	case constants.ViewSelectOperation:
		return SelectOperation(m)
	case constants.ViewApprovals:
		return SelectApproval(m)
	case constants.ViewConfirmation:
		return HandleConfirmationSelection(m)
	case constants.ViewSummary:
		return HandleSummaryConfirmation(m)
	case constants.ViewExecutingAction:
		return HandleExecutionSelection(m)
	case constants.ViewPipelineStatus:
		return HandlePipelineSelection(m)
	case constants.ViewFunctionStatus:
		return HandleFunctionSelection(m)
	default:
		return WrapModel(m), nil
	}
}

// HandleProviderSelection handles the selection of a provider
func HandleProviderSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		providerName := selected[0]

		// Get the provider from the registry
		provider, err := m.Registry.Get(providerName)
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		newModel := m.Clone()
		newModel.ProviderState.ProviderName = providerName

		// Special handling for AWS provider for backward compatibility
		if providerName == "AWS" {
			// Get profiles from the registry
			profiles, err := provider.GetProfiles()
			if err != nil {
				return WrapModel(m), func() tea.Msg {
					return model.ErrMsg{Err: err}
				}
			}

			// Sort the profiles alphabetically
			sort.Strings(profiles)

			// Set profiles and transition to AWS config view
			newModel.Profiles = profiles
			newModel.CurrentView = constants.ViewAWSConfig
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}

		// For other providers, follow the new flow
		// Check if the provider is already authenticated
		if provider.IsAuthenticated() {
			// If authenticated, go to service selection
			newModel.CurrentView = constants.ViewSelectService
		} else {
			// Get authentication methods
			authMethods := provider.GetAuthenticationMethods()

			newModel.ProviderState.AuthState.AvailableMethods = authMethods

			if len(authMethods) == 1 {
				// If only one auth method, select it automatically
				newModel.ProviderState.AuthState.Method = authMethods[0]

				// Get auth config keys for the selected method
				configKeys := provider.GetAuthConfigKeys(authMethods[0])

				if len(configKeys) > 0 {
					// Initialize auth config map if needed
					if newModel.ProviderState.AuthState.AuthConfig == nil {
						newModel.ProviderState.AuthState.AuthConfig = make(map[string]string)
					}
					newModel.CurrentView = constants.ViewAuthConfig
				} else {
					newModel.CurrentView = constants.ViewProviderConfig
				}
			} else if len(authMethods) > 1 {
				// If multiple auth methods, go to auth method selection
				newModel.CurrentView = constants.ViewAuthMethodSelect
			} else {
				// If no auth methods, go to provider config
				newModel.CurrentView = constants.ViewProviderConfig
			}
		}

		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleAWSConfigSelection handles the selection of AWS configuration options
func HandleAWSConfigSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := m.Clone()

		if m.GetAwsProfile() == "" {
			// If no profile is selected, this is profile selection
			profile := selected[0]

			// Check if "Manual Entry" was selected
			if profile == "Manual Entry" {
				newModel.ManualInput = true
				newModel.TextInput.Focus()
				newModel.TextInput.Placeholder = constants.MsgEnterProfile
				return WrapModel(newModel), nil
			}

			newModel.SetAwsProfile(profile)
			view.UpdateTableForView(newModel)
		} else {
			// If profile is already selected, this is region selection
			region := selected[0]

			// Check if "Manual Entry" was selected
			if region == "Manual Entry" {
				newModel.ManualInput = true
				newModel.TextInput.Focus()
				newModel.TextInput.Placeholder = constants.MsgEnterRegion
				return WrapModel(newModel), nil
			}

			newModel.SetAwsRegion(region)

			// Configure the provider with the selected profile and region
			provider, err := m.Registry.Get("AWS")
			if err != nil {
				return WrapModel(m), func() tea.Msg {
					return model.ErrMsg{Err: err}
				}
			}

			// Use LoadConfig instead of Configure to properly initialize the services
			err = provider.LoadConfig(newModel.GetAwsProfile(), region)
			if err != nil {
				return WrapModel(m), func() tea.Msg {
					return model.ErrMsg{Err: err}
				}
			}

			// Move to service selection
			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(newModel)
		}

		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleAuthMethodSelection handles the selection of an authentication method
func HandleAuthMethodSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		methodName := selected[0]

		// Find the selected method in available methods
		var selectedMethod string
		for _, method := range m.ProviderState.AuthState.AvailableMethods {
			if method == methodName {
				selectedMethod = method
				break
			}
		}

		if selectedMethod != "" {
			newModel := m.Clone()
			newModel.ProviderState.AuthState.Method = selectedMethod

			// Get the provider from the registry
			provider, err := m.Registry.Get(m.ProviderState.ProviderName)
			if err != nil {
				return WrapModel(m), func() tea.Msg {
					return model.ErrMsg{Err: err}
				}
			}

			// Get auth config keys for the selected method
			configKeys := provider.GetAuthConfigKeys(selectedMethod)

			if len(configKeys) > 0 {
				// Initialize auth config map if needed
				if newModel.ProviderState.AuthState.AuthConfig == nil {
					newModel.ProviderState.AuthState.AuthConfig = make(map[string]string)
				}
				newModel.CurrentView = constants.ViewAuthConfig
			} else {
				newModel.CurrentView = constants.ViewProviderConfig
			}

			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleAuthConfigSelection handles the selection of authentication configuration options
func HandleAuthConfigSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		configKey := selected[0]

		newModel := m.Clone()
		newModel.ProviderState.AuthState.CurrentAuthConfigKey = configKey
		newModel.ManualInput = true
		newModel.TextInput.Focus()
		newModel.TextInput.Placeholder = "Enter value for " + configKey

		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleProviderConfigSelection handles the selection of provider configuration options
func HandleProviderConfigSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		configKey := selected[0]

		newModel := m.Clone()
		newModel.ProviderState.CurrentConfigKey = configKey
		newModel.ManualInput = true
		newModel.TextInput.Focus()
		newModel.TextInput.Placeholder = "Enter value for " + configKey

		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleEnter handles the Enter key press based on the current view
func HandleEnter(m *model.Model) (tea.Model, tea.Cmd) {
	// If manual input is enabled, handle text input submission
	if m.ManualInput {
		return HandleTextInputSubmission(m)
	}

	// Otherwise, handle table selection
	return HandleTableSelect(m)
}

// HandleTextInputSubmission handles the submission of text input
func HandleTextInputSubmission(m *model.Model) (tea.Model, tea.Cmd) {
	newModel := m.Clone()
	value := m.TextInput.Value()

	switch m.CurrentView {
	case constants.ViewAWSConfig:
		// Handle AWS config input
		if m.GetAwsProfile() == "" {
			// This is profile input
			if value != "" {
				newModel.SetAwsProfile(value)
				newModel.ManualInput = false
				newModel.ResetTextInput()
				view.UpdateTableForView(newModel)
			}
		} else {
			// This is region input
			if value != "" {
				newModel.SetAwsRegion(value)
				newModel.ManualInput = false
				newModel.ResetTextInput()

				// Configure the provider with the selected profile and region
				provider, err := m.Registry.Get("AWS")
				if err != nil {
					return WrapModel(m), func() tea.Msg {
						return model.ErrMsg{Err: err}
					}
				}

				err = provider.Configure(map[string]string{
					"profile": newModel.GetAwsProfile(),
					"region":  value,
				})
				if err != nil {
					return WrapModel(m), func() tea.Msg {
						return model.ErrMsg{Err: err}
					}
				}

				// Move to service selection
				newModel.CurrentView = constants.ViewSelectService
				view.UpdateTableForView(newModel)
			}
		}
	case constants.ViewAuthConfig:
		// Handle auth config input
		if m.ProviderState.AuthState.CurrentAuthConfigKey != "" {
			if newModel.ProviderState.AuthState.AuthConfig == nil {
				newModel.ProviderState.AuthState.AuthConfig = make(map[string]string)
			}
			newModel.ProviderState.AuthState.AuthConfig[m.ProviderState.AuthState.CurrentAuthConfigKey] = value
			newModel.ProviderState.AuthState.CurrentAuthConfigKey = ""
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
		}
	case constants.ViewProviderConfig:
		// Handle provider config input
		if m.ProviderState.CurrentConfigKey != "" {
			if newModel.ProviderState.Config == nil {
				newModel.ProviderState.Config = make(map[string]string)
			}
			newModel.ProviderState.Config[m.ProviderState.CurrentConfigKey] = value
			newModel.ProviderState.CurrentConfigKey = ""
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
		}
	case constants.ViewSummary:
		// Handle summary input (comments for approvals or commit IDs for pipeline starts)
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			newModel.CommitID = value
			newModel.ManualInput = false
			newModel.ResetTextInput()
			newModel.CurrentView = constants.ViewExecutingAction
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		} else {
			newModel.ApprovalComment = value
			newModel.Summary = value
			newModel.ManualInput = false
			newModel.ResetTextInput()
			newModel.CurrentView = constants.ViewExecutingAction
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}

	return WrapModel(newModel), nil
}
