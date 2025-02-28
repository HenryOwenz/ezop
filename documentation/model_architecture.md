# Model Architecture

## Overview

This document describes the model architecture implemented in CloudGate to support multiple cloud providers. The architecture is designed to be flexible, extensible, and maintainable, allowing for easy addition of new cloud providers and services.

## Model Structure

The model architecture consists of several key components:

1. **Core Model**: The central state container for the application
2. **Provider State**: Manages provider-specific state and configuration
3. **Authentication State**: Handles provider authentication
4. **Input State**: Manages user input across different views

### Core Model

The `Model` struct in `internal/ui/model/model.go` is the central state container for the application. It contains:

- UI components (table, text input, spinner)
- View state (current view, manual input mode)
- Provider registry
- Provider state
- Input state
- Legacy fields for backward compatibility

```go
package model

import (
    "sort"
    
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

    // Legacy fields for backward compatibility
    // These will be gradually migrated to the new structure
    // ...
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

### Model Initialization

The model is initialized with default values and provider-specific functionality:

```go
func (m *Model) Init() tea.Cmd {
    m.Regions = constants.DefaultAWSRegions

    // Initialize the AWS provider to get profiles
    if providers.CreateAWSProvider != nil {
        awsProvider := providers.CreateAWSProvider()
        if awsProvider != nil {
            profiles, err := awsProvider.GetProfiles()
            if err == nil && len(profiles) > 0 {
                // Sort the profiles alphabetically
                sort.Strings(profiles)
                m.Profiles = profiles
            } else {
                // Fallback to default profile if there's an error
                m.Profiles = []string{"default"}
            }
        }
    }

    return m.Spinner.Tick
}
```

### Backward Compatibility

To maintain backward compatibility with existing code, the model includes helper methods that bridge the new structure with legacy fields:

```go
// GetAwsProfile returns the AWS profile from the provider config
func (m *Model) GetAwsProfile() string {
    // First check the new structure
    profile := m.GetProviderConfig("profile")
    if profile != "" {
        return profile
    }
    // Fall back to legacy field
    return m.AwsProfile
}

// SetAwsProfile sets the AWS profile in the provider config
func (m *Model) SetAwsProfile(profile string) {
    m.SetProviderConfig("profile", profile)
    // Also set in legacy field for backward compatibility
    m.AwsProfile = profile
}
```

## Provider Architecture

The provider architecture is designed to support multiple cloud providers through a common interface:

1. **Provider Interface**: Defines methods that all cloud providers must implement
2. **Cloud Provider**: Implements cloud-specific functionality
3. **Provider Adapter**: Adapts cloud providers to the provider interface

### Provider Interface

The `Provider` interface in `internal/providers/interfaces.go` defines methods that all cloud providers must implement:

```go
type Provider interface {
    // Basic provider information
    Name() string
    Description() string
    Services() []Service

    // Authentication and configuration
    GetAuthenticationMethods() []string
    GetAuthConfigKeys(method string) []string
    Authenticate(method string, authConfig map[string]string) error
    IsAuthenticated() bool
    GetConfigKeys() []string
    GetConfigOptions(key string) ([]string, error)
    Configure(config map[string]string) error

    // Profile and region management
    GetProfiles() ([]string, error)
    LoadConfig(profile, region string) error
    
    // Cloud-specific operations
    GetApprovals(ctx context.Context) ([]ApprovalAction, error)
    ApproveAction(ctx context.Context, action ApprovalAction, approved bool, comment string) error
    GetStatus(ctx context.Context) ([]PipelineStatus, error)
    StartPipeline(ctx context.Context, pipelineName string, commitID string) error
}
```

### Provider Registration

Providers are registered with the application using a factory pattern to avoid import cycles:

```go
// In internal/providers/factory.go
func InitializeProviders(registry *ProviderRegistry) {
    // Register AWS provider
    awsProvider := CreateAWSProvider()
    if awsProvider != nil {
        registry.Register(awsProvider)
    } else {
        panic("Failed to create AWS provider")
    }
}

// In internal/providers/aws/register.go
func init() {
    // Set the CreateAWSProvider function in the providers package
    providers.CreateAWSProvider = func() providers.Provider {
        return New()
    }
}
```

### Cloud Provider

The cloud provider implementations in `internal/cloud/` contain the actual cloud-specific functionality:

```go
type Provider struct {
    profile  string
    region   string
    services []cloud.Service
}
```

### Provider Adapter

The provider adapters in `internal/providers/` adapt cloud providers to the provider interface:

```go
type CloudProviderAdapter struct {
    provider cloud.Provider
    profile  string
    region   string
}
```

## View Flow

The view flow is managed by the update package, which handles user input and updates the model accordingly:

1. **Navigation**: Handles navigation between views
2. **Selection**: Handles selection of providers, services, categories, and operations
3. **Authentication**: Handles provider authentication
4. **Configuration**: Handles provider configuration

## Direct Provider Interaction

The application now interacts directly with providers for cloud-specific operations:

```go
// Get approvals directly from the provider
ctx := context.Background()
approvals, err := provider.GetApprovals(ctx)
if err != nil {
    return model.ErrMsg{Err: err}
}

// Execute the approval action
ctx := context.Background()
err = provider.ApproveAction(ctx, providerApproval, m.ApproveAction, m.ApprovalComment)
if err != nil {
    return model.ApprovalResultMsg{Err: err}
}

// Get pipeline status directly from the provider
ctx := context.Background()
pipelines, err := provider.GetStatus(ctx)
if err != nil {
    return model.ErrMsg{Err: err}
}

// Execute the pipeline
ctx := context.Background()
err = provider.StartPipeline(ctx, m.SelectedPipeline.Name, m.CommitID)
if err != nil {
    return model.PipelineExecutionMsg{Err: err}
}
```

## Conclusion

The model architecture provides a flexible and extensible foundation for CloudGate, allowing for easy addition of new cloud providers and services. The backward compatibility layer ensures that existing code continues to work while new features are added.

The direct provider interaction pattern simplifies the code and makes it more maintainable by removing the need for complex reflection-based type conversions and intermediate layers. This approach also makes it easier to add new cloud providers in the future.
