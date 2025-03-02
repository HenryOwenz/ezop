# Model Architecture

## Overview

This document describes the model architecture implemented in cloudgate to support multiple cloud providers. The architecture is designed to be flexible, extensible, and maintainable, allowing for easy addition of new cloud providers and services.

## Model Structure

The model architecture consists of several key components:

1. **Core Model**: The central state container for the application
2. **Provider State**: Manages provider-specific state and configuration
3. **Authentication State**: Handles provider authentication
4. **Input State**: Manages user input across different views
5. **Operation-Specific State**: Manages state for specific operations like Lambda Function Status

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
    GetFunctionStatus(ctx context.Context) ([]FunctionStatus, error)
    
    // UI Operation interfaces
    GetCodePipelineManualApprovalOperation() (CodePipelineManualApprovalOperation, error)
    GetPipelineStatusOperation() (PipelineStatusOperation, error)
    GetStartPipelineOperation() (StartPipelineOperation, error)
    GetFunctionStatusOperation() (FunctionStatusOperation, error)
}
```

### Service, Category, and Operation Interfaces

The provider architecture includes interfaces for services, categories, and operations:

```go
type Service interface {
    // Name returns the service's name
    Name() string
    
    // Description returns the service's description
    Description() string
    
    // Categories returns all available categories for this service
    Categories() []Category
}

type Category interface {
    // Name returns the category's name
    Name() string
    
    // Description returns the category's description
    Description() string
    
    // Operations returns all available operations for this category
    Operations() []Operation
    
    // IsUIVisible returns whether this category should be visible in the UI
    IsUIVisible() bool
}

type Operation interface {
    // Name returns the operation's name
    Name() string
    
    // Description returns the operation's description
    Description() string
    
    // Execute executes the operation with the given parameters
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
    
    // IsUIVisible returns whether this operation should be visible in the UI
    IsUIVisible() bool
}
```

### UI Operation Interfaces

The provider architecture includes specialized interfaces for UI operations:

```go
type UIOperation interface {
    // Name returns the operation's name
    Name() string
    
    // Description returns the operation's description
    Description() string
    
    // IsUIVisible returns whether this operation should be visible in the UI
    IsUIVisible() bool
}

type FunctionStatusOperation interface {
    UIOperation
    
    // GetFunctionStatus returns the status of all Lambda functions
    GetFunctionStatus(ctx context.Context) ([]FunctionStatus, error)
}
```

### Data Types

The provider architecture includes data types for various cloud resources:

```go
// FunctionStatus represents the status of a Lambda function
type FunctionStatus struct {
    Name         string
    Runtime      string
    Memory       int32
    Timeout      int32
    LastUpdate   string
    Role         string
    Handler      string
    Description  string
    FunctionArn  string
    CodeSize     int64
    Version      string
    PackageType  string
    Architecture string
    LogGroup     string
}
```

## Provider Registry

The provider registry manages the available providers:

```go
type ProviderRegistry struct {
    providers map[string]Provider
    mu        sync.RWMutex
}

func NewProviderRegistry() *ProviderRegistry {
    return &ProviderRegistry{
        providers: make(map[string]Provider),
    }
}

func (r *ProviderRegistry) Register(provider Provider) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.providers[provider.Name()] = provider
}

func (r *ProviderRegistry) GetProvider(name string) (Provider, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    provider, ok := r.providers[name]
    return provider, ok
}
```

## AWS Provider Implementation

The AWS provider implements the Provider interface:

```go
type Provider struct {
    cloudProvider *aws.Provider
    profile       string
    region        string
    authenticated bool
}

func (p *Provider) GetFunctionStatus(ctx context.Context) ([]providers.FunctionStatus, error) {
    if !p.IsAuthenticated() {
        return nil, fmt.Errorf("provider not authenticated")
    }

    // Load AWS config
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithSharedConfigProfile(p.profile),
        config.WithRegion(p.region),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to load AWS config: %w", err)
    }

    client := lambda.NewFromConfig(cfg)

    // List all functions
    var functions []providers.FunctionStatus
    var marker *string

    for {
        output, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{
            Marker: marker,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to list functions: %w", err)
        }

        // Convert Lambda functions to FunctionStatus
        for _, function := range output.Functions {
            memory := int32(0)
            if function.MemorySize != nil {
                memory = *function.MemorySize
            }

            timeout := int32(0)
            if function.Timeout != nil {
                timeout = *function.Timeout
            }

            // Get architecture (default to x86_64 if not specified)
            architecture := "x86_64"
            if len(function.Architectures) > 0 {
                architecture = string(function.Architectures[0])
            }

            // Get log group if available
            logGroup := ""
            if function.LoggingConfig != nil && function.LoggingConfig.LogGroup != nil {
                logGroup = *function.LoggingConfig.LogGroup
            }

            functions = append(functions, providers.FunctionStatus{
                Name:         aws.ToString(function.FunctionName),
                Runtime:      string(function.Runtime),
                Memory:       memory,
                Timeout:      timeout,
                LastUpdate:   aws.ToString(function.LastModified),
                Role:         aws.ToString(function.Role),
                Handler:      aws.ToString(function.Handler),
                Description:  aws.ToString(function.Description),
                FunctionArn:  aws.ToString(function.FunctionArn),
                CodeSize:     function.CodeSize,
                Version:      aws.ToString(function.Version),
                PackageType:  string(function.PackageType),
                Architecture: architecture,
                LogGroup:     logGroup,
            })
        }

        if output.NextMarker == nil {
            break
        }
        marker = output.NextMarker
    }

    return functions, nil
}

func (p *Provider) GetFunctionStatusOperation() (providers.FunctionStatusOperation, error) {
    if !p.IsAuthenticated() {
        return nil, fmt.Errorf("provider not authenticated")
    }
    return &functionStatusOperation{provider: p}, nil
}
```

## Conclusion

The model architecture provides a flexible and maintainable foundation for cloudgate. By clearly separating concerns and using interfaces, we've created a system that can easily support multiple cloud providers and services. The recent addition of Lambda Function Status demonstrates how the architecture can be extended to support new services.
