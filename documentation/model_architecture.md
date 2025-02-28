# Model Architecture

## Overview

This document describes the model architecture for CloudGate, focusing on how we structure the application state to support multiple cloud providers, services, categories, and operations. The design aims to be flexible, maintainable, and extensible as we add support for more cloud providers beyond AWS.

## Design Goals

1. **Provider Agnostic**: The model should not be tied to any specific cloud provider.
2. **Hierarchical Structure**: The model should mirror the provider-service-category-operation hierarchy.
3. **Extensible**: Adding new providers should not require changing the model structure.
4. **Consistent Interface**: The UI should interact with a consistent interface regardless of the provider.
5. **Separation of Concerns**: Provider-specific logic should be isolated from the UI logic.

## Model Structure

```go
package model

import (
    "github.com/HenryOwenz/cloudgate/internal/providers"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/HenryOwenz/cloudgate/internal/ui/constants"
    "github.com/HenryOwenz/cloudgate/internal/ui/styles"
    tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
    // UI Components
    Table     table.Model
    TextInput textinput.Model
    Spinner   spinner.Model
    Styles    styles.Styles

    // Window dimensions
    Width  int
    Height int

    // View state
    CurrentView constants.View
    ManualInput bool
    Err         error
    Error       string // Error message
    Success     string // Success message

    // Loading state
    IsLoading  bool
    LoadingMsg string

    // Provider Registry
    Registry *providers.ProviderRegistry

    // Provider state
    ProviderState ProviderState

    // Input state
    InputState InputState
}

// ProviderState represents the state of the selected provider, service, category, and operation
type ProviderState struct {
    // Selected provider
    ProviderName string
    
    // Provider configuration
    Config map[string]string // Generic configuration (e.g., "profile", "region")
    
    // Available configuration options
    ConfigOptions map[string][]string // e.g., "profile" -> ["default", "dev", "prod"]
    
    // Current configuration key being set
    CurrentConfigKey string
    
    // Authentication state
    AuthState AuthenticationState
    
    // Selected service, category, and operation
    SelectedService   *ServiceInfo
    SelectedCategory  *CategoryInfo
    SelectedOperation *OperationInfo
    
    // Provider-specific state (stored as generic interface{})
    ProviderSpecificState map[string]interface{}
}

// AuthenticationState represents the authentication state for different providers
type AuthenticationState struct {
    // Current authentication method
    Method string
    
    // Authentication configuration
    AuthConfig map[string]string
    
    // Available authentication methods
    AvailableMethods []string
    
    // Current authentication config key being set
    CurrentAuthConfigKey string
    
    // Authentication status
    IsAuthenticated bool
    
    // Error message if authentication failed
    AuthError string
}

// ServiceInfo represents information about a service
type ServiceInfo struct {
    ID          string
    Name        string
    Description string
    Available   bool
}

// CategoryInfo represents information about a category
type CategoryInfo struct {
    ID          string
    Name        string
    Description string
    Available   bool
}

// OperationInfo represents information about an operation
type OperationInfo struct {
    ID          string
    Name        string
    Description string
}

// InputState represents the state of user input
type InputState struct {
    // Generic input fields
    TextValues map[string]string // e.g., "comment" -> "This is a comment"
    BoolValues map[string]bool   // e.g., "approve" -> true
    
    // Operation-specific state
    OperationState map[string]interface{} // Operation-specific state
}
```

## Provider Interface

To support this model structure, we extend the provider interface:

```go
// Provider interface defines methods that all cloud providers must implement
type Provider interface {
    // Name returns the provider's name
    Name() string

    // Description returns the provider's description
    Description() string

    // Services returns all available services for this provider
    Services() []Service

    // GetAuthenticationMethods returns the available authentication methods
    GetAuthenticationMethods() []string
    
    // GetAuthConfigKeys returns the configuration keys required for an authentication method
    GetAuthConfigKeys(method string) []string
    
    // Authenticate authenticates with the provider using the given method and configuration
    Authenticate(method string, authConfig map[string]string) error
    
    // IsAuthenticated returns whether the provider is authenticated
    IsAuthenticated() bool
    
    // GetConfigKeys returns the configuration keys required by this provider
    GetConfigKeys() []string
    
    // GetConfigOptions returns the available options for a configuration key
    GetConfigOptions(key string) ([]string, error)
    
    // Configure configures the provider with the given configuration
    Configure(config map[string]string) error
}
```

## Authentication Strategies

### AWS Authentication

For AWS, we use a simplified authentication approach:

```go
// GetAuthenticationMethods returns the available authentication methods
func (p *AWSProvider) GetAuthenticationMethods() []string {
    // AWS doesn't need explicit authentication methods
    return []string{}
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (p *AWSProvider) GetAuthConfigKeys(method string) []string {
    // AWS doesn't need auth config keys
    return []string{}
}

// Authenticate authenticates with the provider using the given method and configuration
func (p *AWSProvider) Authenticate(method string, authConfig map[string]string) error {
    // AWS doesn't need explicit authentication
    return nil
}

// IsAuthenticated returns whether the provider is authenticated
func (p *AWSProvider) IsAuthenticated() bool {
    // AWS is always "authenticated" if we have a profile
    return p.profile != ""
}

// GetConfigKeys returns the configuration keys required by this provider
func (p *AWSProvider) GetConfigKeys() []string {
    return []string{"profile", "region"}
}

// Configure configures the provider with the given configuration
func (p *AWSProvider) Configure(config map[string]string) error {
    profile := config["profile"]
    if profile == "" {
        return fmt.Errorf("profile is required")
    }
    
    region := config["region"]
    if region == "" {
        return fmt.Errorf("region is required")
    }
    
    p.profile = profile
    p.region = region
    
    // Initialize AWS client
    ctx := context.Background()
    cfg, err := aws.LoadDefaultConfig(ctx,
        aws.WithSharedConfigProfile(profile),
        aws.WithRegion(region),
    )
    if err != nil {
        return fmt.Errorf("failed to load AWS config: %w", err)
    }
    
    p.client = codepipeline.NewFromConfig(cfg)
    return nil
}
```

### Azure Authentication (Future)

For Azure, we'll support multiple authentication methods:

```go
// GetAuthenticationMethods returns the available authentication methods
func (p *AzureProvider) GetAuthenticationMethods() []string {
    return []string{"cli", "config-dir"}
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (p *AzureProvider) GetAuthConfigKeys(method string) []string {
    switch method {
    case "cli":
        return []string{} // No config needed, uses default CLI auth
    case "config-dir":
        return []string{"config-dir"}
    default:
        return []string{}
    }
}

// Authenticate authenticates with the provider using the given method and configuration
func (p *AzureProvider) Authenticate(method string, authConfig map[string]string) error {
    switch method {
    case "cli":
        // Use default Azure CLI auth
        cmd := exec.Command("az", "login")
        return cmd.Run()
    case "config-dir":
        configDir := authConfig["config-dir"]
        if configDir == "" {
            return fmt.Errorf("config directory is required")
        }
        // Set AZURE_CONFIG_DIR environment variable
        os.Setenv("AZURE_CONFIG_DIR", configDir)
        return nil
    default:
        return fmt.Errorf("unsupported authentication method: %s", method)
    }
}
```

### GCP Authentication (Future)

For GCP, we'll support service account authentication:

```go
// GetAuthenticationMethods returns the available authentication methods
func (p *GCPProvider) GetAuthenticationMethods() []string {
    return []string{"service-account", "adc"}
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (p *GCPProvider) GetAuthConfigKeys(method string) []string {
    switch method {
    case "service-account":
        return []string{"service-account-path"}
    case "adc":
        return []string{} // No config needed, uses ADC
    default:
        return []string{}
    }
}

// Authenticate authenticates with the provider using the given method and configuration
func (p *GCPProvider) Authenticate(method string, authConfig map[string]string) error {
    switch method {
    case "service-account":
        serviceAccountPath := authConfig["service-account-path"]
        if serviceAccountPath == "" {
            return fmt.Errorf("service account path is required")
        }
        // Set GOOGLE_APPLICATION_CREDENTIALS environment variable
        os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", serviceAccountPath)
        return nil
    case "adc":
        // Use Application Default Credentials
        return nil
    default:
        return fmt.Errorf("unsupported authentication method: %s", method)
    }
}

// In internal/ui/update/navigation_handlers.go

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
        
        view.UpdateTableForView(newModel)
        return WrapModel(newModel), nil
    }
    return WrapModel(m), nil
}

// In internal/ui/update/navigation_handlers.go

// HandleConfigSelection handles the selection of a configuration option
func HandleConfigSelection(m *model.Model) (tea.Model, tea.Cmd) {
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

// In internal/ui/update/navigation_handlers.go

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

// In internal/ui/update/navigation_handlers.go

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

// In internal/ui/view/table.go

// getRowsForView returns the appropriate rows for the current view
func getRowsForView(m *model.Model) []table.Row {
    switch m.CurrentView {
    // ... existing cases for other views
    
    case constants.ViewAuthMethodSelect:
        rows := make([]table.Row, len(m.ProviderState.AuthState.AvailableMethods))
        for i, method := range m.ProviderState.AuthState.AvailableMethods {
            description := getAuthMethodDescription(m.ProviderState.ProviderName, method)
            rows[i] = table.Row{method, description}
        }
        return rows
    
    case constants.ViewAuthConfig:
        key := m.ProviderState.AuthState.CurrentAuthConfigKey
        options, ok := m.ProviderState.ConfigOptions[key]
        if !ok {
            return []table.Row{}
        }
        
        rows := make([]table.Row, len(options)+1)
        rows[0] = table.Row{"Manual Entry"}
        for i, option := range options {
            rows[i+1] = table.Row{option}
        }
        return rows
    
    case constants.ViewProviderConfig:
        key := m.ProviderState.CurrentConfigKey
        options, ok := m.ProviderState.ConfigOptions[key]
        if !ok {
            return []table.Row{}
        }
        
        rows := make([]table.Row, len(options)+1)
        rows[0] = table.Row{"Manual Entry"}
        for i, option := range options {
            rows[i+1] = table.Row{option}
        }
        return rows
    
    // ... other cases
    }
    
    return []table.Row{}
}

// getAuthMethodDescription returns a description for an authentication method
func getAuthMethodDescription(providerName, method string) string {
    descriptions := map[string]map[string]string{
        "AWS": {
            "profile": "Use AWS profile from ~/.aws/credentials",
        },
        "Azure": {
            "cli": "Use Azure CLI authentication",
            "config-dir": "Use Azure configuration directory",
        },
        "GCP": {
            "service-account": "Use GCP service account key file",
            "adc": "Use Application Default Credentials",
        },
    }
    
    if providerDescriptions, ok := descriptions[providerName]; ok {
        if description, ok := providerDescriptions[method]; ok {
            return description
        }
    }
    
    return ""
}

// In internal/ui/constants/constants.go

// View represents the current view state
type View int

const (
    // View states
    ViewProviders View = iota
    ViewAWSConfig
    ViewSelectService
    ViewSelectCategory
    ViewSelectOperation
    ViewApprovals
    ViewConfirmation
    ViewSummary
    ViewExecutingAction
    ViewPipelineStatus
    ViewPipelineStages
    ViewError
    ViewSuccess
    ViewHelp
    
    // New view states for the updated model
    ViewAuthMethodSelect
    ViewAuthConfig
    ViewProviderConfig
    ViewAuthError
)

// Authentication method constants
const (
    // AWS authentication methods
    AWSProfileAuth = "profile"
    
    // Azure authentication methods
    AzureCliAuth = "cli"
    AzureConfigDirAuth = "config-dir"
    
    // GCP authentication methods
    GCPServiceAccountAuth = "service-account"
    GCPApplicationDefaultAuth = "adc"
)

// Configuration key constants
const (
    // AWS configuration keys
    AWSProfileKey = "profile"
    AWSRegionKey = "region"
    
    // Azure configuration keys
    AzureSubscriptionKey = "subscription"
    AzureLocationKey = "location"
    AzureTenantKey = "tenant"
    AzureConfigDirKey = "config-dir"
    
    // GCP configuration keys
    GCPProjectKey = "project"
    GCPZoneKey = "zone"
    GCPRegionKey = "region"
    GCPServiceAccountKey = "service-account-path"
)

// In internal/ui/model/model.go

// Backward compatibility methods

// GetAwsProfile returns the AWS profile from the provider config
func (m *Model) GetAwsProfile() string {
    return m.GetProviderConfig(AWSProfileKey)
}

// SetAwsProfile sets the AWS profile in the provider config
func (m *Model) SetAwsProfile(profile string) {
    m.SetProviderConfig(AWSProfileKey, profile)
}

// GetAwsRegion returns the AWS region from the provider config
func (m *Model) GetAwsRegion() string {
    return m.GetProviderConfig(AWSRegionKey)
}

// SetAwsRegion sets the AWS region in the provider config
func (m *Model) SetAwsRegion(region string) {
    m.SetProviderConfig(AWSRegionKey, region)
}

// GetApprovalComment returns the approval comment from the input state
func (m *Model) GetApprovalComment() string {
    return m.GetInputText("approval-comment")
}

// SetApprovalComment sets the approval comment in the input state
func (m *Model) SetApprovalComment(comment string) {
    m.SetInputText("approval-comment", comment)
}

// GetApproveAction returns the approve action from the input state
func (m *Model) GetApproveAction() bool {
    return m.GetInputBool("approve-action")
}

// SetApproveAction sets the approve action in the input state
func (m *Model) SetApproveAction(approve bool) {
    m.SetInputBool("approve-action", approve)
}

// GetCommitID returns the commit ID from the input state
func (m *Model) GetCommitID() string {
    return m.GetInputText("commit-id")
}

// SetCommitID sets the commit ID in the input state
func (m *Model) SetCommitID(commitID string) {
    m.SetInputText("commit-id", commitID)
}

// GetManualCommitID returns whether to use manual commit ID
func (m *Model) GetManualCommitID() bool {
    return m.GetInputBool("manual-commit-id")
}

// SetManualCommitID sets whether to use manual commit ID
func (m *Model) SetManualCommitID(manual bool) {
    m.SetInputBool("manual-commit-id", manual)
}

// GetSelectedApproval returns the selected approval from the provider-specific state
func (m *Model) GetSelectedApproval() interface{} {
    return m.ProviderState.ProviderSpecificState["selected-approval"]
}

// SetSelectedApproval sets the selected approval in the provider-specific state
func (m *Model) SetSelectedApproval(approval interface{}) {
    m.ProviderState.ProviderSpecificState["selected-approval"] = approval
}

// GetSelectedPipeline returns the selected pipeline from the provider-specific state
func (m *Model) GetSelectedPipeline() interface{} {
    return m.ProviderState.ProviderSpecificState["selected-pipeline"]
}

// SetSelectedPipeline sets the selected pipeline in the provider-specific state
func (m *Model) SetSelectedPipeline(pipeline interface{}) {
    m.ProviderState.ProviderSpecificState["selected-pipeline"] = pipeline
}

// GetApprovals returns the approvals from the provider-specific state
func (m *Model) GetApprovals() []interface{} {
    approvals, ok := m.ProviderState.ProviderSpecificState["approvals"].([]interface{})
    if !ok {
        return nil
    }
    return approvals
}

// SetApprovals sets the approvals in the provider-specific state
func (m *Model) SetApprovals(approvals []interface{}) {
    m.ProviderState.ProviderSpecificState["approvals"] = approvals
}

// GetPipelines returns the pipelines from the provider-specific state
func (m *Model) GetPipelines() []interface{} {
    pipelines, ok := m.ProviderState.ProviderSpecificState["pipelines"].([]interface{})
    if !ok {
        return nil
    }
    return pipelines
}

// SetPipelines sets the pipelines in the provider-specific state
func (m *Model) SetPipelines(pipelines []interface{}) {
    m.ProviderState.ProviderSpecificState["pipelines"] = pipelines
}

```

Clarification on Backward Compatibility
The backward compatibility adapter methods are not about leaving AWS out of our new model design. On the contrary, AWS will fully participate in the new model structure just like Azure and GCP will.
Here's what the adapter methods are actually for:
Purpose of the Adapter Methods
1. Gradual Migration: These methods allow us to gradually migrate the codebase without breaking existing functionality. We can implement the new model structure and start using it for new code while existing code continues to work.
Code Transition: As we refactor the codebase, we can update one file at a time to use the new model structure, rather than having to update everything at once.
Maintaining Functionality: During the transition period, these methods ensure that both old and new code can interact with the model correctly.
AWS in the New Model
AWS will be fully integrated into the new model structure:
Same Provider Interface: AWS will implement the same extended provider interface as Azure and GCP.
Same Model Structure: AWS configuration and state will be stored in the same model structure as other providers.
Same Authentication Flow: While AWS has a simpler authentication flow (it doesn't need explicit authentication), it will still go through the same flow in the code, just with fewer steps.
Same Configuration Flow: AWS will use the same configuration flow as other providers, with profile and region as its configuration keys.
Example of AWS in the New Model
Here's how AWS will fit into the new model structure:

```go
// AWS provider state in the new model
m.ProviderState.ProviderName = "AWS"
m.ProviderState.Config = map[string]string{
    "profile": "default",
    "region": "us-west-2",
}

// And when we need to access AWS-specific data:
// Store AWS-specific data
m.ProviderState.ProviderSpecificState["approvals"] = awsApprovals
m.ProviderState.ProviderSpecificState["pipelines"] = awsPipelines

// Retrieve AWS-specific data
approvals := m.ProviderState.ProviderSpecificState["approvals"].([]aws.ApprovalAction)
pipelines := m.ProviderState.ProviderSpecificState["pipelines"].([]aws.PipelineStatus)

```
Transition Strategy
The transition strategy is:
Implement New Model: Create the new model structure with all the necessary fields and methods.
Add Adapter Methods: Add the backward compatibility methods to bridge between old and new code.
Update AWS Implementation: Update the AWS provider to implement the extended provider interface.
Migrate Code Gradually: Update one file at a time to use the new model structure.
5. Remove Adapter Methods: Once all code has been migrated, remove the adapter methods.
This approach ensures that AWS is fully integrated into the new model structure from the beginning, while allowing us to gradually migrate the codebase without breaking existing functionality.



## UI Flow

The UI flow for provider selection and configuration:

1. User selects a provider
2. If the provider has authentication methods:
   a. User selects an authentication method
   b. User provides authentication configuration
   c. Provider authenticates
3. User provides provider configuration
4. Provider is configured
5. User selects a service, category, and operation

## Benefits

1. **Flexibility**: Each provider can define whatever configuration parameters make sense for its authentication model.
2. **Extensibility**: Adding a new provider doesn't require changing the model structure, just adding new provider-specific code.
3. **Consistency**: The UI code can work with a consistent interface regardless of the provider.
4. **Separation of Concerns**: Provider-specific logic is isolated from the UI logic.

## Implementation Strategy

1. Create the new model structure in `internal/ui/model/model.go`
2. Update the provider interface in `internal/providers/interfaces.go`
3. Implement the provider interface for AWS in `internal/cloud/aws/provider.go`
4. Update the handlers in `internal/ui/update/` to use the new model structure
5. Update the views in `internal/ui/view/` to use the new model structure
6. Add support for Azure and GCP in the future

## Conclusion

This model architecture provides a flexible, maintainable, and extensible foundation for cloudgate. It allows us to support multiple cloud providers with different authentication and configuration requirements while maintaining a consistent user interface. The design is provider-agnostic and can be extended to support new providers without changing the core model structure. 
