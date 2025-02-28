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

Key components:
- `providers.Provider` - Interface for adapted cloud providers
- `providers.Service` - Interface for adapted cloud services
- `providers.Category` - Interface for adapted service categories
- `providers.Operation` - Interface for adapted service operations

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

## Flow of Execution

1. The UI interacts with the providers layer through the `providers.Provider` interface
2. The providers layer adapts calls to the cloud layer through the adapters
3. The cloud layer executes the actual cloud provider operations
4. Results flow back through the adapters to the UI

## Adding New Components

### Adding a New Cloud Provider

1. Create a new package in `internal/cloud/{provider_name}/`
2. Implement the `cloud.Provider` interface
3. Register the provider with the adapter in `internal/providers/registry.go`

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

## Example: AWS CodePipeline Implementation

```go
// Cloud Layer (internal/cloud/aws/codepipeline/)
type Service struct {
    profile    string
    region     string
    categories []cloud.Category
}

// Providers Layer (internal/providers/aws/codepipeline/)
type Service struct {
    profile    string
    region     string
    categories []providers.Category
}
```

## Best Practices

1. **Keep Layers Separate**: Avoid direct dependencies between the UI and the cloud layer
2. **Consistent Naming**: Use consistent naming across both layers
3. **Interface Alignment**: Keep interfaces in both layers aligned
4. **Minimal Adapters**: Adapters should be thin and add minimal functionality
5. **Complete Implementation**: Ensure all methods are properly implemented in both layers

## Conclusion

The dual-layer architecture provides a robust foundation for CloudGate, allowing for clean separation of concerns and easy extension to support multiple cloud providers. By following this pattern, we ensure that the codebase remains maintainable and scalable as we add more functionality. 