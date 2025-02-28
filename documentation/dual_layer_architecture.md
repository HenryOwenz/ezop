# Dual-Layer Architecture Pattern

## Overview

CloudGate implements a dual-layer architecture pattern to provide a clean separation between cloud provider implementations and the application's business logic. This document describes this architecture, its benefits, and how to extend it when adding new cloud providers or services.

## Architecture Layers

The architecture consists of two primary layers:

1. **Cloud Layer** (`internal/cloud/`) - Contains the core implementations for each cloud provider
2. **Providers Layer** (`internal/providers/`) - Contains adapters and interfaces used by the application

### Cloud Layer

The cloud layer is responsible for:

- Implementing cloud provider-specific logic
- Communicating with cloud provider APIs
- Defining the core interfaces for providers, services, categories, and operations

Key components:
- `cloud.Provider` - Interface for cloud providers
- `cloud.Service` - Interface for cloud services
- `cloud.Category` - Interface for service categories
- `cloud.Operation` - Interface for service operations

### Providers Layer

The providers layer is responsible for:

- Adapting cloud implementations to the application's needs
- Providing a consistent interface for the UI and other components
- Adding additional functionality or validation
- Implementing direct provider operations for cloud-specific functionality

Key components:
- `providers.Provider` - Interface for adapted cloud providers
- `providers.Service` - Interface for adapted cloud services
- `providers.Category` - Interface for adapted service categories
- `providers.Operation` - Interface for adapted service operations
- `providers.ApprovalAction` - Type for pipeline approval actions
- `providers.PipelineStatus` - Type for pipeline status information

## Enhanced Provider Interface

The provider interface has been enhanced to include direct cloud-specific operations:

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

## Provider Registration

To avoid import cycles, providers are registered using a factory pattern:

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

## Adapter Pattern

The connection between the two layers is implemented using the Adapter pattern:

```go
// CloudProviderAdapter adapts a cloud.Provider to a providers.Provider
type CloudProviderAdapter struct {
    provider cloud.Provider
}

// CloudServiceAdapter adapts a cloud.Service to a providers.Service
type CloudServiceAdapter struct {
    service cloud.Service
}

// CloudCategoryAdapter adapts a cloud.Category to a providers.Category
type CloudCategoryAdapter struct {
    category cloud.Category
}

// CloudOperationAdapter adapts a cloud.Operation to a providers.Operation
type CloudOperationAdapter struct {
    operation cloud.Operation
}
```

These adapters ensure that implementations from the cloud layer can be used through the providers layer interfaces.

## Direct Provider Implementation

For cloud-specific operations, the provider implementation directly handles the functionality:

```go
// GetApprovals returns all pending approvals for the provider
func (p *Provider) GetApprovals(ctx context.Context) ([]providers.ApprovalAction, error) {
    if !p.IsAuthenticated() {
        return nil, fmt.Errorf("provider not authenticated")
    }

    approvals, err := p.GetPendingApprovals(ctx)
    if err != nil {
        return nil, err
    }

    // Convert internal ApprovalAction to providers.ApprovalAction
    providerApprovals := make([]providers.ApprovalAction, len(approvals))
    for i, approval := range approvals {
        providerApprovals[i] = providers.ApprovalAction{
            PipelineName: approval.PipelineName,
            StageName:    approval.StageName,
            ActionName:   approval.ActionName,
            Token:        approval.Token,
        }
    }

    return providerApprovals, nil
}
```

## Flow of Execution

1. The UI interacts with the providers layer through the `providers.Provider` interface
2. For basic operations, the providers layer adapts calls to the cloud layer through the adapters
3. For cloud-specific operations, the providers layer directly implements the functionality
4. Results flow back through the providers layer to the UI

## Adding New Components

### Adding a New Cloud Provider

1. Create a new package in `internal/cloud/{provider_name}/`
2. Implement the `cloud.Provider` interface
3. Create a new package in `internal/providers/{provider_name}/`
4. Implement the `providers.Provider` interface, adapting the cloud provider
5. Create a registration function in `internal/providers/{provider_name}/register.go`
6. Update `internal/providers/factory.go` to include the new provider

### Adding a New Service

1. Create a new package in `internal/cloud/{provider_name}/{service_name}/`
2. Implement the `cloud.Service` interface
3. Register the service with the provider

### Adding a New Category

1. Create a new type in the service package
2. Implement the `cloud.Category` interface
3. Register the category with the service

### Adding a New Operation

1. Create a new type in the category package
2. Implement the `cloud.Operation` interface
3. Register the operation with the category

## Benefits

1. **Separation of Concerns**: Clear separation between cloud provider implementations and application logic
2. **Testability**: Each layer can be tested independently
3. **Extensibility**: New cloud providers can be added without changing the application logic
4. **Consistency**: Consistent interfaces across different cloud providers
5. **Maintainability**: Changes to one layer don't necessarily affect the other
6. **Direct Implementation**: Cloud-specific operations are implemented directly in the provider, simplifying the code

## Example: AWS Provider Implementation

```go
// Cloud Layer (internal/cloud/aws/)
type Provider struct {
    profile  string
    region   string
    services []cloud.Service
}

// Providers Layer (internal/providers/aws/)
type Provider struct {
    cloudProvider *aws.Provider
    profile       string
    region        string
    authenticated bool
}

// Direct implementation of cloud-specific operations
func (p *Provider) GetApprovals(ctx context.Context) ([]providers.ApprovalAction, error) {
    // Implementation details...
}

func (p *Provider) ApproveAction(ctx context.Context, action providers.ApprovalAction, approved bool, comment string) error {
    // Implementation details...
}

func (p *Provider) GetStatus(ctx context.Context) ([]providers.PipelineStatus, error) {
    // Implementation details...
}

func (p *Provider) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
    // Implementation details...
}
```

## Best Practices

1. **Keep Layers Separate**: Avoid direct dependencies between the UI and the cloud layer
2. **Consistent Naming**: Use consistent naming across both layers
3. **Interface Alignment**: Keep interfaces in both layers aligned
4. **Minimal Adapters**: Adapters should be thin and add minimal functionality
5. **Complete Implementation**: Ensure all methods are properly implemented in both layers
6. **Direct Implementation**: Implement cloud-specific operations directly in the provider
7. **Avoid Import Cycles**: Use factory patterns to avoid import cycles

## Conclusion

The dual-layer architecture provides a robust foundation for CloudGate, allowing for clean separation of concerns and easy extension to support multiple cloud providers. By following this pattern, we ensure that the codebase remains maintainable and scalable as we add more functionality. The direct implementation of cloud-specific operations in the provider layer simplifies the code and makes it more maintainable. 