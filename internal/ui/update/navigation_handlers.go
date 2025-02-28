package update

import (
	"fmt"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/providers"
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
func HandleEnter(m *model.Model) (tea.Model, tea.Cmd) {
	// Special handling for manual input in different views
	if m.ManualInput {
		newModel := m.Clone()
		value := strings.TrimSpace(m.TextInput.Value())

		if value == "" {
			// If empty, just exit manual input mode
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}

		// Handle manual input based on the current view
		switch m.CurrentView {
		case constants.ViewAuthConfig:
			// Set auth config value
			newModel.SetAuthConfig(m.ProviderState.AuthState.CurrentAuthConfigKey, value)
			newModel.ManualInput = false
			newModel.ResetTextInput()

			// Get the provider
			provider, err := m.Registry.Get(m.ProviderState.ProviderName)
			if err != nil {
				return WrapModel(newModel), func() tea.Msg {
					return model.ErrMsg{Err: err}
				}
			}

			// Get the next auth config key
			authConfigKeys := provider.GetAuthConfigKeys(m.ProviderState.AuthState.Method)
			currentKeyIndex := -1
			for i, key := range authConfigKeys {
				if key == m.ProviderState.AuthState.CurrentAuthConfigKey {
					currentKeyIndex = i
					break
				}
			}

			// If there are more auth config keys, show the next one
			if currentKeyIndex < len(authConfigKeys)-1 {
				nextKey := authConfigKeys[currentKeyIndex+1]
				newModel.ProviderState.AuthState.CurrentAuthConfigKey = nextKey

				// Get options for the next key
				options, err := provider.GetConfigOptions(nextKey)
				if err == nil {
					newModel.ProviderState.ConfigOptions[nextKey] = options
				}

				view.UpdateTableForView(newModel)
				return WrapModel(newModel), nil
			}

			// If all auth config keys are set, authenticate
			err = provider.Authenticate(newModel.ProviderState.AuthState.Method, newModel.ProviderState.AuthState.AuthConfig)
			if err != nil {
				newModel.ProviderState.AuthState.AuthError = err.Error()
				newModel.CurrentView = constants.ViewAuthError
			} else {
				newModel.ProviderState.AuthState.IsAuthenticated = true

				// Move to provider configuration
				configKeys := provider.GetConfigKeys()
				if len(configKeys) > 0 {
					firstKey := configKeys[0]
					options, err := provider.GetConfigOptions(firstKey)
					if err == nil {
						newModel.ProviderState.ConfigOptions[firstKey] = options
					}

					newModel.ProviderState.CurrentConfigKey = firstKey
					newModel.CurrentView = constants.ViewProviderConfig
				} else {
					newModel.CurrentView = constants.ViewSelectService
				}
			}

			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil

		case constants.ViewProviderConfig:
			// Set provider config value
			newModel.SetProviderConfig(m.ProviderState.CurrentConfigKey, value)
			newModel.ManualInput = false
			newModel.ResetTextInput()

			// Get the provider
			provider, err := m.Registry.Get(m.ProviderState.ProviderName)
			if err != nil {
				return WrapModel(newModel), func() tea.Msg {
					return model.ErrMsg{Err: err}
				}
			}

			// Get the next config key
			configKeys := provider.GetConfigKeys()
			currentKeyIndex := -1
			for i, key := range configKeys {
				if key == m.ProviderState.CurrentConfigKey {
					currentKeyIndex = i
					break
				}
			}

			// If there are more config keys, show the next one
			if currentKeyIndex < len(configKeys)-1 {
				nextKey := configKeys[currentKeyIndex+1]
				newModel.ProviderState.CurrentConfigKey = nextKey

				// Get options for the next key
				options, err := provider.GetConfigOptions(nextKey)
				if err == nil {
					newModel.ProviderState.ConfigOptions[nextKey] = options
				}

				view.UpdateTableForView(newModel)
				return WrapModel(newModel), nil
			}

			// If all config keys are set, configure the provider
			err = provider.Configure(newModel.ProviderState.Config)
			if err != nil {
				return WrapModel(newModel), func() tea.Msg {
					return model.ErrMsg{Err: err}
				}
			}

			// Move to service selection
			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil

		case constants.ViewAWSConfig:
			// For backward compatibility
			// Get the entered value
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
						return model.ErrMsg{Err: err}
					}
				}

				newModel.CurrentView = constants.ViewSelectService
				view.UpdateTableForView(newModel)
			}
			return WrapModel(newModel), nil

		case constants.ViewSummary:
			// For backward compatibility
			if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
				newModel.CommitID = value
				newModel.ManualCommitID = true
			} else if m.SelectedApproval != nil {
				newModel.ApprovalComment = value
			}

			newModel.ManualInput = false
			newModel.CurrentView = constants.ViewExecutingAction
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}

	// Regular view handling
	switch m.CurrentView {
	case constants.ViewProviders:
		return HandleProviderSelection(m)
	case constants.ViewAuthMethodSelect:
		return HandleAuthMethodSelection(m)
	case constants.ViewAuthConfig:
		return HandleAuthConfigSelection(m)
	case constants.ViewProviderConfig:
		return HandleProviderConfigSelection(m)
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
func HandleProviderSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		providerName := selected[0]

		// Initialize providers if not already done
		if len(m.Registry.GetProviderNames()) == 0 {
			providers.InitializeProviders(m.Registry)
		}

		// Get the provider
		provider, err := m.Registry.Get(providerName)
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		// Set the provider name
		newModel := m.Clone()
		newModel.ProviderState.ProviderName = providerName

		// Get authentication methods
		authMethods := provider.GetAuthenticationMethods()
		newModel.ProviderState.AuthState.AvailableMethods = authMethods

		// If no authentication methods, skip to configuration
		if len(authMethods) == 0 {
			// Move to provider configuration
			configKeys := provider.GetConfigKeys()
			if len(configKeys) > 0 {
				// Get options for first config key
				firstKey := configKeys[0]
				options, err := provider.GetConfigOptions(firstKey)
				if err == nil {
					newModel.ProviderState.ConfigOptions[firstKey] = options
				}

				newModel.ProviderState.CurrentConfigKey = firstKey
				newModel.CurrentView = constants.ViewProviderConfig
			} else {
				newModel.CurrentView = constants.ViewSelectService
			}
		} else {
			// If there's only one authentication method, use it
			if len(authMethods) == 1 {
				newModel.ProviderState.AuthState.Method = authMethods[0]

				// Get auth config keys
				authConfigKeys := provider.GetAuthConfigKeys(authMethods[0])

				// If no auth config keys, authenticate directly
				if len(authConfigKeys) == 0 {
					err := provider.Authenticate(authMethods[0], map[string]string{})
					if err != nil {
						newModel.ProviderState.AuthState.AuthError = err.Error()
						newModel.CurrentView = constants.ViewAuthError
					} else {
						newModel.ProviderState.AuthState.IsAuthenticated = true

						// Move to provider configuration
						configKeys := provider.GetConfigKeys()
						if len(configKeys) > 0 {
							firstKey := configKeys[0]
							options, err := provider.GetConfigOptions(firstKey)
							if err == nil {
								newModel.ProviderState.ConfigOptions[firstKey] = options
							}

							newModel.ProviderState.CurrentConfigKey = firstKey
							newModel.CurrentView = constants.ViewProviderConfig
						} else {
							newModel.CurrentView = constants.ViewSelectService
						}
					}
				} else {
					// Show auth config view
					newModel.ProviderState.AuthState.CurrentAuthConfigKey = authConfigKeys[0]

					// Get options for the auth config key
					options, err := provider.GetConfigOptions(authConfigKeys[0])
					if err == nil {
						newModel.ProviderState.ConfigOptions[authConfigKeys[0]] = options
					}

					newModel.CurrentView = constants.ViewAuthConfig
				}
			} else {
				// Show auth method selection view
				newModel.CurrentView = constants.ViewAuthMethodSelect
			}
		}

		// For backward compatibility with AWS, also set the AWS config view
		if providerName == "AWS" {
			newModel.CurrentView = constants.ViewAWSConfig
		}

		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleAWSConfigSelection handles the selection of AWS profile or region
func HandleAWSConfigSelection(m *model.Model) (tea.Model, tea.Cmd) {
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
					return model.ErrMsg{Err: err}
				}
			}

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
		authMethod := selected[0]

		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		newModel := m.Clone()
		newModel.ProviderState.AuthState.Method = authMethod

		// Get auth config keys
		authConfigKeys := provider.GetAuthConfigKeys(authMethod)

		// If no auth config keys, authenticate directly
		if len(authConfigKeys) == 0 {
			err := provider.Authenticate(authMethod, map[string]string{})
			if err != nil {
				newModel.ProviderState.AuthState.AuthError = err.Error()
				newModel.CurrentView = constants.ViewAuthError
			} else {
				newModel.ProviderState.AuthState.IsAuthenticated = true

				// Move to provider configuration
				configKeys := provider.GetConfigKeys()
				if len(configKeys) > 0 {
					firstKey := configKeys[0]
					options, err := provider.GetConfigOptions(firstKey)
					if err == nil {
						newModel.ProviderState.ConfigOptions[firstKey] = options
					}

					newModel.ProviderState.CurrentConfigKey = firstKey
					newModel.CurrentView = constants.ViewProviderConfig
				} else {
					newModel.CurrentView = constants.ViewSelectService
				}
			}
		} else {
			// Show auth config view
			newModel.ProviderState.AuthState.CurrentAuthConfigKey = authConfigKeys[0]

			// Get options for the auth config key
			options, err := provider.GetConfigOptions(authConfigKeys[0])
			if err == nil {
				newModel.ProviderState.ConfigOptions[authConfigKeys[0]] = options
			}

			newModel.CurrentView = constants.ViewAuthConfig
		}

		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleAuthConfigSelection handles the selection of an authentication configuration option
func HandleAuthConfigSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		configValue := selected[0]

		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		newModel := m.Clone()

		// Handle "Manual Entry" option
		if configValue == "Manual Entry" {
			newModel.ManualInput = true
			newModel.TextInput.Focus()
			newModel.TextInput.Placeholder = fmt.Sprintf("Enter %s", newModel.ProviderState.AuthState.CurrentAuthConfigKey)
			return WrapModel(newModel), nil
		}

		// Set the auth config value
		newModel.SetAuthConfig(newModel.ProviderState.AuthState.CurrentAuthConfigKey, configValue)

		// Get the next auth config key
		authConfigKeys := provider.GetAuthConfigKeys(newModel.ProviderState.AuthState.Method)
		currentKeyIndex := -1
		for i, key := range authConfigKeys {
			if key == newModel.ProviderState.AuthState.CurrentAuthConfigKey {
				currentKeyIndex = i
				break
			}
		}

		// If there are more auth config keys, show the next one
		if currentKeyIndex < len(authConfigKeys)-1 {
			nextKey := authConfigKeys[currentKeyIndex+1]
			newModel.ProviderState.AuthState.CurrentAuthConfigKey = nextKey

			// Get options for the next key
			options, err := provider.GetConfigOptions(nextKey)
			if err == nil {
				newModel.ProviderState.ConfigOptions[nextKey] = options
			}

			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}

		// If all auth config keys are set, authenticate
		err = provider.Authenticate(newModel.ProviderState.AuthState.Method, newModel.ProviderState.AuthState.AuthConfig)
		if err != nil {
			newModel.ProviderState.AuthState.AuthError = err.Error()
			newModel.CurrentView = constants.ViewAuthError
		} else {
			newModel.ProviderState.AuthState.IsAuthenticated = true

			// Move to provider configuration
			configKeys := provider.GetConfigKeys()
			if len(configKeys) > 0 {
				firstKey := configKeys[0]
				options, err := provider.GetConfigOptions(firstKey)
				if err == nil {
					newModel.ProviderState.ConfigOptions[firstKey] = options
				}

				newModel.ProviderState.CurrentConfigKey = firstKey
				newModel.CurrentView = constants.ViewProviderConfig
			} else {
				newModel.CurrentView = constants.ViewSelectService
			}
		}

		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleProviderConfigSelection handles the selection of a provider configuration option
func HandleProviderConfigSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		configValue := selected[0]

		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		newModel := m.Clone()

		// Handle "Manual Entry" option
		if configValue == "Manual Entry" {
			newModel.ManualInput = true
			newModel.TextInput.Focus()
			newModel.TextInput.Placeholder = fmt.Sprintf("Enter %s", newModel.ProviderState.CurrentConfigKey)
			return WrapModel(newModel), nil
		}

		// Set the config value
		newModel.SetProviderConfig(newModel.ProviderState.CurrentConfigKey, configValue)

		// Get the next config key
		configKeys := provider.GetConfigKeys()
		currentKeyIndex := -1
		for i, key := range configKeys {
			if key == newModel.ProviderState.CurrentConfigKey {
				currentKeyIndex = i
				break
			}
		}

		// If there are more config keys, show the next one
		if currentKeyIndex < len(configKeys)-1 {
			nextKey := configKeys[currentKeyIndex+1]
			newModel.ProviderState.CurrentConfigKey = nextKey

			// Get options for the next key
			options, err := provider.GetConfigOptions(nextKey)
			if err == nil {
				newModel.ProviderState.ConfigOptions[nextKey] = options
			}

			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}

		// If all config keys are set, configure the provider
		err = provider.Configure(newModel.ProviderState.Config)
		if err != nil {
			return WrapModel(newModel), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		// Move to service selection
		newModel.CurrentView = constants.ViewSelectService
		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}
